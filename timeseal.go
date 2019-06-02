// Copyright © 2019 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

import (
	"fmt"
	"time"
)

const (
	tsKey = "Timestamp (FICS) v1.0 - programmed by Henrik Gram."
)

// Encode a byte array using the encoding scheme mandated by timeseal2
func encode(b []byte, l int) []byte {
	s := make([]byte, l+30)
	copy(s[:l], b)
	s[l] = 0x18
	l++
	ts := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	copy(s[l:], ts)
	l += len(ts)
	s[l] = 0x19
	l++
	for ; (l % 12) != 0; l++ {
		s[l] = 0x31
	}

	for n := 0; n < l; n += 12 {
		s[n] ^= s[n+11]
		s[n+11] ^= s[n]
		s[n] ^= s[n+11]
		s[n+2] ^= s[n+9]
		s[n+9] ^= s[n+2]
		s[n+2] ^= s[n+9]
		s[n+4] ^= s[n+7]
		s[n+7] ^= s[n+4]
		s[n+4] ^= s[n+7]
	}

	for n := 0; n < l; n++ {
		var x = int8(((s[n] | 0x80) ^ tsKey[n%50]) - 32)
		s[n] = byte(x)
	}

	s[l] = 0x80
	l++
	s[l] = 0x0a
	l++
	return s[:l]
}
