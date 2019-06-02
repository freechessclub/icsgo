// Copyright © 2019 Free Chess Club <hi@freechess.club>
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/freechessclub/icsgo"
)

func main() {
	// connect as guest to the FICS server
	client, err := icsgo.NewClient(&icsgo.Config{
		DisableKeepAlive: true,
		DisableTimeseal:  true,
		Debug:            true,
	}, "freechess.org:5000", "guest", "")
	if err != nil {
		log.Fatalf("error creating new FICS client: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer client.Destroy()
		for {
			select {
			default:
				msgs, err := client.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Fatalf("error receiving server output: %v", err)
					return
				}

				if msgs != nil {
					fmt.Printf("%v\n", msgs)
				}
			case <-done:
				return
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("fics_client% ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			client.Destroy()
			log.Fatalf("error reading console input: %v", err)
		}

		err = client.Send(cmd)
		if err != nil || cmd == "exit\n" {
			close(done)
			break
		}
	}

}
