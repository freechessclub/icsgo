// Copyright © 2018 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

// Config represents the configuration parameters supported by icsgo
type Config struct {
	ServerAddr     *string
	UserPrompt     *string
	User           *string
	PasswordPrompt *string
	Password       *string
	ICSPrompt      *string
	KeepAlive      *bool
	Timeseal       *bool
	ConnTimeout    *int
	ConnRetries    *int
}

// String value for setting configuration parameter
func String(s string) *string {
	return &s
}

// Bool value for setting configuration parameter
func Bool(b bool) *bool {
	return &b
}

// Int value for setting configuration parameter
func Int(i int) *int {
	return &i
}

// NewConfig returns a new default configuration
func NewConfig(cfg *Config) *Config {
	defaultCfg := Config{
		ServerAddr:     String("freechess.org:5000"),
		UserPrompt:     String("login:"),
		User:           String("guest"),
		PasswordPrompt: String("password:"),
		ICSPrompt:      String("fics%"),
		KeepAlive:      Bool(false),
		Timeseal:       Bool(true),
		ConnTimeout:    Int(2),
		ConnRetries:    Int(5),
	}

	if cfg == nil {
		return &defaultCfg
	}

	if cfg.ServerAddr == nil {
		cfg.ServerAddr = defaultCfg.ServerAddr
	}

	if cfg.UserPrompt == nil {
		cfg.UserPrompt = defaultCfg.UserPrompt
	}

	if cfg.User == nil {
		cfg.User = defaultCfg.User
	}

	if cfg.PasswordPrompt == nil {
		cfg.PasswordPrompt = defaultCfg.PasswordPrompt
	}

	if cfg.ICSPrompt == nil {
		cfg.ICSPrompt = defaultCfg.ICSPrompt
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = defaultCfg.KeepAlive
	}

	if cfg.Timeseal == nil {
		cfg.Timeseal = defaultCfg.Timeseal
	}

	if cfg.ConnTimeout == nil {
		cfg.ConnTimeout = defaultCfg.ConnTimeout
	}

	if cfg.ConnRetries == nil {
		cfg.ConnRetries = defaultCfg.ConnRetries
	}

	return cfg
}

// SetServerAddr sets the ICS server address (default: "freechess.org:5000")
func (c *Config) SetServerAddr(addr string) *Config {
	c.ServerAddr = String(addr)
	return c
}

// SetUserPrompt sets the ICS server login prompt (default: "login:")
func (c *Config) SetUserPrompt(prompt string) *Config {
	c.UserPrompt = String(prompt)
	return c
}

// SetUser sets the ICS server username (default: "guest")
func (c *Config) SetUser(user string) *Config {
	c.User = String(user)
	return c
}

// SetPasswordPrompt sets the ICS server password prompt (default: "password")
func (c *Config) SetPasswordPrompt(prompt string) *Config {
	c.PasswordPrompt = String(prompt)
	return c
}

// SetPassword sets the ICS server password
func (c *Config) SetPassword(password string) *Config {
	c.Password = String(password)
	return c
}

// SetICSPrompt sets the ICS server prompt, (default: "fics%")
func (c *Config) SetICSPrompt(prompt string) *Config {
	c.ICSPrompt = String(prompt)
	return c
}

// SetTimeseal enables timeseal v2 on the ICS connection (default: true)
func (c *Config) SetTimeseal(value bool) *Config {
	c.Timeseal = Bool(value)
	return c
}

// SetKeepAlive enables ICS server connection from idling out (default: false)
func (c *Config) SetKeepAlive(value bool) *Config {
	c.KeepAlive = Bool(value)
	return c
}

// SetConnTimeout sets the timeout (in seconds) when connecting to the
// ICS server (default: 2 seconds)
func (c *Config) SetConnTimeout(timeout int) *Config {
	c.ConnTimeout = Int(timeout)
	return c
}

// SetConnRetries sets the number of times the client should retry
// when connecting to the ICS server (default: 5)
func (c *Config) SetConnRetries(retries int) *Config {
	c.ConnRetries = Int(retries)
	return c
}
