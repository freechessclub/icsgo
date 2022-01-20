// Copyright © 2019 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	gameMoveRE  *regexp.Regexp
	gameStartRE *regexp.Regexp
	gameEndRE   *regexp.Regexp
	chTellRE    *regexp.Regexp
	pTellRE     *regexp.Regexp
	kibitzRE    *regexp.Regexp
	toldMsgRE   *regexp.Regexp
)

// type of game end messages
const (
	Unknown = iota
	Resign
	Disconnect
	Checkmate
	TimeForfeit
	Draw
	Adjourn
	Abort
)

func init() {
	// game move
	// <12> rnbqkb-r pppppppp -----n-- -------- ----P--- -------- PPPPKPPP RNBQ-BNR B -1 0 0 1 1 0 7 Newton Einstein 1 2 12 39 39 119 122 2 K/e1-e2 (0:06) Ke2 0
	gameMoveRE = regexp.MustCompile(`<12>\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([rnbqkpRNBQKP\-]{8})\s([BW\-])\s(?:\-?[0-7])\s(?:[01])\s(?:[01])\s(?:[01])\s(?:[01])\s(?:[0-9]+)\s([0-9]+)\s([a-zA-Z]+)\s([a-zA-Z]+)\s(\-?[0-3])\s([0-9]+)\s([0-9]+)\s(?:[0-9]+)\s(?:[0-9]+)\s(\-?[0-9]+)\s(\-?[0-9]+)\s([0-9]+)\s(?:\S+)\s\((?:[0-9]+)\:(?:[0-9]+)\)\s(\S+)\s(?:[01])\s(?:[0-9]+)\s(?:[0-9]+)\s*`)

	// {Game 117 (GuestMDPS vs. guestl) Creating unrated blitz match.}
	gameStartRE = regexp.MustCompile(`(?s)^\s*\{Game\s([0-9]+)\s\(([a-zA-Z]+)\svs\.\s([a-zA-Z]+)\)\sCreating.*\}.*`)

	gameEndRE = regexp.MustCompile(`(?s)^[^\(\):]*(?:Game\s[0-9]+:.*)?\{Game\s([0-9]+)\s\(([a-zA-Z]+)\svs\.\s([a-zA-Z]+)\)\s([a-zA-Z]+)\s([a-zA-Z0-9\s]+)\}\s(?:[012/]+-[012/]+)?.*`)

	// channel tell
	chTellRE = regexp.MustCompile(`(?s)^([a-zA-Z]+)(?:\([A-Z\*]+\))*\(([0-9]+)\):\s+(.*)`)

	// private tell
	pTellRE = regexp.MustCompile(`(?s)^([a-zA-Z]+)(?:[\(\[][A-Z0-9\*\-]+[\)\]])* (?:tells you|says):\s+(.*)`)

	// kibitz/whispers
	kibitzRE = regexp.MustCompile(`(?s)^([a-zA-Z]+)(?:\([A-Z0-9\*\-]+\))*\[([0-9]+)\] (?:kibitzes|whispers):\s+(.*)`)

	// told status
	toldMsgRE = regexp.MustCompile(`\((?:told|kibitzed) .+\)`)
}

func style12ToFEN(b []byte) string {
	str := string(b[:])
	var fen string
	count := 0
	for i := 0; i < 8; i++ {
		if str[i] == '-' {
			count++
			if i == 7 {
				fen += strconv.Itoa(count)
			}
		} else {
			if count > 0 {
				fen += strconv.Itoa(count)
				count = 0
			}
			fen += string(str[i])
		}
	}
	return fen
}

func unsafeAtoi(b []byte) uint32 {
	i, _ := strconv.Atoi(string(b))
	return uint32(i)
}

func getGameResult(p1, p2, who, action string) (string, string, uint32) {
	action = strings.TrimSpace(action)
	switch action {
	case "resigns":
		if p1 == who {
			return p2, p1, Resign
		} else if p2 == who {
			return p1, p2, Resign
		}
	case "forfeits by disconnection":
		if p1 == who {
			return p2, p1, Disconnect
		} else if p2 == who {
			return p1, p2, Disconnect
		}
	case "checkmated":
		if p1 == who {
			return p2, p1, Checkmate
		} else if p2 == who {
			return p1, p2, Checkmate
		}
	case "forfeits on time":
		if p1 == who {
			return p2, p1, TimeForfeit
		} else if p2 == who {
			return p1, p2, TimeForfeit
		}
	case "aborted on move 1":
	case "aborted by mutual agreement":
		return p1, p2, Abort
	case "drawn by mutual agreement":
	case "drawn because both players ran out of time":
	case "drawn by repetition":
	case "drawn by the 50 move rule":
	case "drawn due to length":
	case "was drawn":
	case "player has mating material":
	case "drawn by adjudication":
	case "drawn by stalemate":
		return p1, p2, Draw
	case "adjourned by mutual agreement":
		return p1, p2, Adjourn
	}
	return p1, p2, Unknown
}

func decodeMessages(msg []byte) []interface{} {
	if len(msg) == 0 {
		return nil
	}

	msg = toldMsgRE.ReplaceAll(msg, []byte{})
	if msg == nil || bytes.Equal(msg, []byte("\n")) {
		return nil
	}

	matches := gameMoveRE.FindSubmatch(msg)
	if matches != nil && len(matches) >= 18 {
		m := bytes.Split(msg, []byte("\n"))
		if len(m) > 1 {
			var msgs []interface{}
			for i := 0; i < len(m); i++ {
				if len(m[i]) > 0 {
					msgs = append(msgs, decodeMessages(m[i])...)
				}
			}
			return msgs
		}

		fen := ""
		for i := 1; i < 8; i++ {
			fen += style12ToFEN(matches[i][:])
			fen += "/"
		}
		fen += style12ToFEN(matches[8][:])
		r, _ := strconv.Atoi(string(matches[13][:]))
		role := int32(r)

		return []interface{}{
			&GameMove{
				Fen:       fen,
				Turn:      string(matches[9][:]),
				GameId:    unsafeAtoi(matches[10][:]),
				WhiteName: string(matches[11][:]),
				BlackName: string(matches[12][:]),
				Role:      role,
				Time:      unsafeAtoi(matches[14][:]),
				Inc:       unsafeAtoi(matches[15][:]),
				WhiteTime: unsafeAtoi(matches[16][:]),
				BlackTime: unsafeAtoi(matches[17][:]),
				MoveNo:    unsafeAtoi(matches[18][:]),
				Move:      string(matches[19][:]),
			},
		}
	}

	matches = gameStartRE.FindSubmatch(msg)
	if matches != nil && len(matches) > 2 {
		return []interface{}{
			&GameStart{
				GameId:    unsafeAtoi(matches[1][:]),
				PlayerOne: string(matches[2][:]),
				PlayerTwo: string(matches[3][:]),
			},
		}
	}

	matches = gameEndRE.FindSubmatch(msg)
	if matches != nil && len(matches) > 4 {
		p1 := string(matches[2][:])
		p2 := string(matches[3][:])
		who := string(matches[4][:])
		action := string(matches[5][:])

		winner, loser, reason := getGameResult(p1, p2, who, action)
		return []interface{}{
			&GameEnd{
				GameId:  unsafeAtoi(matches[1][:]),
				Winner:  winner,
				Loser:   loser,
				Reason:  reason,
				Message: string(msg),
			},
		}
	}

	matches = chTellRE.FindSubmatch(msg)
	if matches != nil && len(matches) > 3 {
		return []interface{}{
			&ChannelTell{
				Channel: string(matches[2][:]),
				User:    string(matches[1][:]),
				Message: string(bytes.Replace(matches[3][:], []byte("\n"), []byte{}, -1)),
			},
		}
	}

	matches = pTellRE.FindSubmatch(msg)
	if matches != nil && len(matches) > 2 {
		return []interface{}{
			&PrivateTell{
				User:    string(matches[1][:]),
				Message: string(bytes.Replace(matches[2][:], []byte("\n"), []byte{}, -1)),
			},
		}
	}

	matches = kibitzRE.FindSubmatch(msg)
	if matches != nil && len(matches) > 3 {
		return []interface{}{
			&ChannelTell{
				Channel: "Game " + string(matches[2][:]),
				User:    string(matches[1][:]),
				Message: string(bytes.Replace(matches[3][:], []byte("\n"), []byte{}, -1)),
			},
		}
	}

	return []interface{}{
		&Message{
			Message: string(msg),
		},
	}
}
