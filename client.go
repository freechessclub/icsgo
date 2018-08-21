// Copyright © 2018 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/ziutek/telnet"
)

// Client represents a connection instance to the ICS server
type Client struct {
	Config   *Config
	Conn     *telnet.Conn
	Username string
}

// Greeting message to use when enabling timeseal on client connection
const (
	hello = "TIMESEAL2|freeseal|Free Chess Club|"
)

//
func NewClient(cfg *Config) (*Client, error) {

	client := &Client{
		Config: cfg,
	}

	conn, err := dial(cfg)
	if err != nil {
		return nil, err
	}

	client.Conn = conn

	if *cfg.Timeseal {
		client.Send(hello)
	}

	username, err := client.Login()
	if err != nil {
		return nil, err
	}

	client.Username = username

	// set seek 0
	// set echo 1
	// set style 12
	// set interface www.freechess.club

	if *cfg.KeepAlive {
		go keepAlive(client)
	}

	go connReader(client)
	return client, nil
}

func dial(cfg *Config) (*telnet.Conn, error) {
	addr := *cfg.ServerAddr
	timeout := float32(*cfg.ConnTimeout)
	retries := *cfg.ConnRetries

	var conn *telnet.Conn
	var connected = false
	var err error

	for attempts := 1; attempts <= retries && connected != true; attempts++ {
		ts := time.Duration(timeout) * time.Second
		log.Printf("Connecting to ICS server %s (attempt %d of %d)...", addr, attempts, retries)
		conn, err = telnet.DialTimeout("tcp", addr, ts)
		if err != nil {
			timeout *= 1.5
			continue
		}
		connected = true
	}

	if err != nil || connected == false {
		return nil, fmt.Errorf("connecting to server %s: %v", addr, err)
	}

	log.Printf("Connected to ICS server %s!", addr)
	return conn, nil
}

func keepAlive(client *Client) {
	for {
		time.Sleep(58 * time.Minute)
		client.Send("ping")
	}
}

func connReader(client *Client) {
	cfg := client.Config
	conn := client.Conn
	for {
		conn.SetReadDeadline(time.Now().Add(3600 * time.Second))
		out, err := client.ReadUntil(*cfg.ICSPrompt)
		if err != nil {
			client.Destroy()
			return
		}

		if len(out) == 0 {
			continue
		}

		msgs, err := decodeMessage(out)
		if err != nil {
			log.Println("Error decoding message")
		}
		if msgs == nil {
			continue
		}

		arr, ok := msgs.([]interface{})
		if ok && len(arr) == 0 {
			continue
		}

		bs, err := json.Marshal(msgs)
		if err != nil {
			log.Println("Error marshaling message")
		}
		if bs == nil {
			continue
		}
	}
}

//
func (client *Client) ReadUntil(delims ...string) ([]byte, error) {
	bs, err := client.Conn.ReadUntil(delims...)
	if err != nil {
		return nil, err
	}

	if *client.Config.Timeseal {
		for {
			i := bytes.Index(bs, []byte{'[', 'G', ']', 0x00})
			if i == -1 {
				break
			}
			client.Send(string([]byte{0x02, 0x39}))
			bs = append(bs[:i], bs[i+4:]...)
		}
	}

	bs = bytes.Replace(bs, []byte("\u0007"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\x00"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\\   "), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("\r"), []byte{}, -1)
	bs = bytes.Replace(bs, []byte("fics%"), []byte{}, -1)
	bs = bytes.TrimSpace(bs)

	return bs, nil
}

//
func (client *Client) Send(cmd string) error {
	client.Conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	bs := []byte(cmd)
	if *client.Config.Timeseal {
		bs = Encode(bs, len(cmd))
	}
	_, err := client.Conn.Write(bs)
	return err
}

//
func (client *Client) Destroy() {
	client.Send("exit")
	client.Conn.Close()
}

//
func (client *Client) Login() (string, error) {
	conn := client.Conn
	cfg := client.Config

	if conn == nil {
		return "", fmt.Errorf("client not connected")
	}

	username := *cfg.User
	password := *cfg.Password
	var prompt string
	// guests have no passwords
	if username != "guest" && len(password) > 0 {
		prompt = *cfg.PasswordPrompt
	} else {
		prompt = "Press return to enter the server as"
		password = ""
	}

	// wait for the login prompt
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	client.ReadUntil(*cfg.UserPrompt)

	_, err := sendAndReadUntil(conn, username, prompt)
	if err != nil {
		return "", fmt.Errorf("creating new login session for %s: %v", username, err)
	}

	// wait for the password prompt
	out, err := sendAndReadUntil(conn, password, "****\n")
	if err != nil {
		return "", fmt.Errorf("failed authentication for %s: %v", username, err)
	}

	re := regexp.MustCompile("\\*\\*\\*\\* Starting FICS session as ([a-zA-Z]+)(?:\\(U\\))? \\*\\*\\*\\*")
	user := re.FindSubmatch(out)
	if user != nil && len(user) > 1 {
		username = string(user[1][:])
		log.Printf("Logged in as %s", username)
		return username, nil
	}

	return "", fmt.Errorf("invalid password for %s", username)
}
