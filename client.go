// Copyright © 2019 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate protoc --go_out=. types.proto

package icsgo

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// Config represents the configuration parameters supported by the icsgo client
type Config struct {
	UserPrompt       string
	PasswordPrompt   string
	ICSPrompt        string
	DisableKeepAlive bool
	DisableTimeseal  bool
	TimesealHello    string
	ConnTimeout      int
	ConnRetries      int
	Debug            bool
}

// DefaultConfig represents the default configuration of icsgo client
var DefaultConfig = &Config{
	UserPrompt:       "login:",
	PasswordPrompt:   "password:",
	ICSPrompt:        "fics%",
	DisableKeepAlive: false,
	DisableTimeseal:  false,
	TimesealHello:    "TIMESEAL2|freeseal|icsgo|",
	ConnTimeout:      2,
	ConnRetries:      5,
	Debug:            false,
}

// Client represents a new ICS client
type Client struct {
	config   *Config
	conn     *Conn
	username string
}

func getConfig(cfg *Config) *Config {
	if cfg == nil {
		cfg = DefaultConfig
	}

	// merge partial config with default config parameters
	if cfg.UserPrompt == "" {
		cfg.UserPrompt = DefaultConfig.UserPrompt
	}

	if cfg.PasswordPrompt == "" {
		cfg.PasswordPrompt = DefaultConfig.PasswordPrompt
	}

	if cfg.ICSPrompt == "" {
		cfg.ICSPrompt = DefaultConfig.ICSPrompt
	}

	if cfg.TimesealHello == "" {
		cfg.TimesealHello = DefaultConfig.TimesealHello
	}

	if cfg.ConnTimeout == 0 {
		cfg.ConnTimeout = DefaultConfig.ConnTimeout
	}

	if cfg.ConnRetries == 0 {
		cfg.ConnRetries = DefaultConfig.ConnRetries
	}

	return cfg
}

// NewClient creates a new ICS client
func NewClient(cfg *Config, addr, username, password string) (*Client, error) {
	cfg = getConfig(cfg)
	retries := cfg.ConnRetries
	timeout := time.Duration(cfg.ConnTimeout) * time.Second
	conn, err := Dial(addr, retries, timeout, !cfg.DisableTimeseal, cfg.Debug)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new connection")
	}

	if !cfg.DisableTimeseal {
		conn.Write(cfg.TimesealHello)
	}

	username, err = login(conn, username, password, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate to server")
	}

	if !cfg.DisableKeepAlive {
		go keepAlive(conn)
	}

	return &Client{
		config:   cfg,
		conn:     conn,
		username: username,
	}, nil
}

// Send sends a message to the ICS server
func (client *Client) Send(msg string) error {
	return client.conn.Write(msg)
}

// Recv receives messages from the ICS server
func (client *Client) Recv() ([]interface{}, error) {
	out, err := client.conn.ReadUntil(client.config.ICSPrompt)
	if err != nil {
		return nil, err
	}

	return decodeMessages(out), nil
}

// Destroy destroys a client instance
func (client *Client) Destroy() {
	client.Send("exit")
	client.conn.Close()
}

func keepAlive(conn *Conn) {
	for {
		time.Sleep(58 * time.Minute)
		conn.Write("ping")
	}
}

//
func login(conn *Conn, username, password string, cfg *Config) (string, error) {
	if conn == nil {
		return "", fmt.Errorf("client not connected")
	}

	// wait for the login prompt
	_, err := conn.ReadUntil(cfg.UserPrompt)
	if err != nil {
		return "", fmt.Errorf("creating new login session for %s: %v", username, err)
	}

	conn.Write(username)

	var prompt string
	// guests have no passwords
	if username != "guest" && len(password) > 0 {
		prompt = cfg.PasswordPrompt
	} else {
		prompt = ":"
		password = ""
	}

	// wait for the password prompt
	_, err = conn.ReadUntil(prompt)
	if err != nil {
		return "", fmt.Errorf("creating new login session for %s: %v", username, err)
	}

	conn.Write(password)

	out, err := conn.ReadUntil("****\n")
	if err != nil {
		return "", fmt.Errorf("failed authentication for %s: %v", username, err)
	}

	re := regexp.MustCompile("\\*\\*\\*\\* Starting FICS session as ([a-zA-Z]+)(?:\\(U\\))?")
	user := re.FindSubmatch(out)
	if user != nil && len(user) > 1 {
		username = string(user[1][:])
		log.Printf("Logged in as %s", username)
		return username, nil
	}

	return "", fmt.Errorf("invalid password for %s", username)
}
