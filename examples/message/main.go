// Package main contains an example usage of FCM client.
//
// MIT License
//
// Copyright (c) 2016 Angel Del Castillo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/jimmy-go/fcm"
)

var (
	token = flag.String("device-token", "", "FCM token.")
	key   = flag.String("server-key", "", "FCM token.")
)

func main() {
	flag.Parse()

	// 1 second timeout
	timeout := time.Second * 5
	client, err := fcm.New(*key, timeout)
	if err != nil {
		log.Fatal(err)
	}

	msg := fcm.NewMessage(*token)
	msg.TimeToLive = 25
	msg.Data.Add("title", "Title demo")
	msg.Data.Add("message", "Hello world! at: "+time.Now().Format(time.RFC3339))
	msg.Data.Add("some-var", 1)
	res, err := client.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	// show response, don't use marshal this way on production.
	b, _ := json.Marshal(res)
	log.Printf("FCM message response : [%v]", string(b))
}
