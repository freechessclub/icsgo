// Copyright © 2018 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package icsgo

// func (c *ICSClient) sendAndReadUntil(cmd string, delims ...string) ([]byte, error) {
// 	err := send(cmd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return readUntil(conn, delims...)
// }

// func (c *ICSClient) keepAlive(timeout time.Duration) {
// 	var lastResponse int64
// 	atomic.StoreInt64(&lastResponse, time.Now().UnixNano())
// 	s.ws.SetPongHandler(func(msg string) error {
// 		atomic.StoreInt64(&lastResponse, time.Now().UnixNano())
// 		return nil
// 	})

// 	for {
// 		// write
// 		time.Sleep(timeout / 2)
// 		if atomic.LoadInt64(&lastResponse) < time.Now().Add(-timeout).UnixNano() {
// 			s.end()
// 			return
// 		}
// 	}
// }
