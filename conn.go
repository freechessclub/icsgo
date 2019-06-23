// Copyright © 2019 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/ziutek/telnet"
)

const (
	timesealHello = "TIMESEAL2|freeseal|icsgo|"
)

// Conn represents a connection to the ICS server
type Conn struct {
	// whether timeseal is enabled on the connection
	timeseal bool
	// whether debug/verbose logging is enabled
	debug bool
	// the underlying telnet connection
	conn *telnet.Conn
}

// Dial creates a new connection
func Dial(addr string, retries int, timeout time.Duration, timeseal, debug bool) (*Conn, error) {
	connected := false

	var conn *telnet.Conn
	var err error

	for attempts := 1; attempts <= retries && connected != true; attempts++ {
		log.Printf("connecting to ICS server %s (attempt %d of %d)...", addr, attempts, retries)
		conn, err = telnet.DialTimeout("tcp", addr, timeout)
		if err != nil {
			timeout = time.Duration(float64(timeout) * 1.5)
			continue
		}
		connected = true
	}

	if err != nil || connected == false {
		return nil, fmt.Errorf("connecting to server %s: %v", addr, err)
	}

	log.Printf("connected to ICS server %s! (timeseal: %t)", addr, timeseal)

	c := &Conn{
		timeseal: timeseal,
		debug:    debug,
		conn:     conn,
	}

	if timeseal {
		c.Write(timesealHello)
	}

	return c, nil
}

// ReadUntilTimeout reads messages from the connection until the given prompt is encountered
// or until the given timeout duration has surpassed
func (c *Conn) ReadUntilTimeout(prompt string, timeout time.Duration) ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(timeout))

	bs, err := c.conn.ReadUntil(prompt)
	if err != nil {
		return nil, err
	}

	if c.debug {
		log.Printf("< %s", string(bs))
	}

	if c.timeseal {
		for {
			i := bytes.Index(bs, []byte{'[', 'G', ']', 0x00})
			if i == -1 {
				break
			}
			c.Write(string([]byte{0x02, 0x39}))
			bs = append(bs[:i], bs[i+4:]...)
		}
	}

	bs = bytes.Replace(bs, []byte("\u0007"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\x00"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\\   "), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\r"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte(prompt), []byte{}, -1)
	bs = bytes.TrimSpace(bs)

	return bs, nil
}

// ReadUntil reads messages from the connection until the given prompt is encountered
func (c *Conn) ReadUntil(prompt string) ([]byte, error) {
	return c.ReadUntilTimeout(prompt, 3600*time.Second)
}

// Write writes the given message on the open connection
func (c *Conn) Write(msg string) error {
	c.conn.SetWriteDeadline(time.Now().Add(20 * time.Second))
	if c.debug {
		log.Printf("> %s", msg)
	}

	bs := []byte(msg)
	if c.timeseal {
		bs = encode(bs, len(msg))
	}

	_, err := c.conn.Conn.Write(bs)
	return err
}

// Close closes the connection to the ICS server
func (c *Conn) Close() {
	c.conn.Close()
}
