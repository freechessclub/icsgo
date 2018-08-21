// Copyright © 2018 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

import (
	"testing"
)

var defaultCfg = Config{
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

func TestDefaultConfig(t *testing.T) {
	cfg := NewConfig(&Config{})
	if cfg.ServerAddr != defaultCfg.ServerAddr {
		t.Errorf("error getting default ServerAddr configuration")
	}
}

func TestPartialConfig(t *testing.T) {
	cfg := NewConfig(&Config{}).SetServerAddr("chessclub.com:5000")
	if cfg.ServerAddr != "chessclub.com:5000" {
		t.Errorf("error setting ServerAddr configuration")
	}
}
