// Package fcm contains Firebase Cloud Messaging client.
// Only http requests are implemented, if xmpp is required
// open an issue.
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
package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// fcmEndpoint defines FCM endpoint for internal usage.
	// It's declared a var to allow tests run multiple cases.
	fcmEndpoint = Endpoint

	// ErrInvalidServerKey is returned when server key is
	// empty or doesn't have serverKeyLen length.
	ErrInvalidServerKey = errors.New("fcm: invalid server key")

	// ErrInvalidToken is returned when device token is
	// empty.
	ErrInvalidToken = errors.New("fcm: empty token")

	// ErrRequestFail is returned when server response status
	// is not 200.
	ErrRequestFail = errors.New("fcm: request fail")
)

const (
	// Endpoint defines Firebase Cloud Message endpoint.
	Endpoint = "https://fcm.googleapis.com/fcm/send"

	// serverKeyLen expetected server key length.
	serverKeyLen = 39
)

// FCM Firebase Cloud Messaging client.
type FCM struct {
	ServerKey string

	client *http.Client
}

// New returns a new FCM client.
func New(serverKey string, timeout time.Duration) (*FCM, error) {
	if len(serverKey) != serverKeyLen {
		return nil, ErrInvalidServerKey
	}
	fcm := &FCM{
		ServerKey: serverKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
	return fcm, nil
}

// Send sends an http request to FCM endpoint.
// server key.
//
// see: https://firebase.google.com/docs/cloud-messaging/server
func (f *FCM) Send(message *Message) (*Response, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("fcm: marshal: %s", err)
	}
	req, err := http.NewRequest("POST", fcmEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("fcm: encoding: %s", err)
	}
	req.Header.Add("Authorization", "key="+f.ServerKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fcm: do: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrRequestFail
	}

	// read multiple times because we want to see body
	// response if can't be unmarshalled.
	buf := bytes.NewBuffer([]byte{})
	body := io.TeeReader(resp.Body, buf)

	var res *Response
	err = json.NewDecoder(body).Decode(&res)
	if err != nil {
		// Response can't be unmarshalled, this will mean
		// that response is a string from FCM server.
		return nil, fmt.Errorf("fcm: %s", buf.String())
	}
	return res, nil
}

// Message contains all available parameters for FCM request.
//
// see: https://firebase.google.com/docs/cloud-messaging/http-server-ref#table1
type Message struct {
	To                    string   `json:"to,omitempty"`
	RegistrationIDs       []string `json:"registration_ids,omitempty"`
	Condition             string   `json:"condition,omitempty"`
	CollapseKey           string   `json:"collapse_key"`
	Priority              int      `json:"priority"`
	ContentAvailable      bool     `json:"content_available,omitempty"`
	DelayWhileIdle        bool     `json:"delay_while_idle,omitempty"`
	TimeToLive            int      `json:"time_to_live,omitempty"`
	RestrictedPackageName string   `json:"restricted_package_name,omitempty"`
	DryRun                bool     `json:"dry_run,omitempty"`
	Data                  KeyValue `json:"data,omitempty"`
	Notification          KeyValue `json:"notification,omitempty"`
}

// NewMessage returns a new Message with 'to' property set.
// For more options use '&Message{}' directly.
func NewMessage(to string) *Message {
	m := &Message{
		To:           to,
		Data:         KeyValue{},
		Notification: KeyValue{},
	}
	return m
}

// RegIDs shorthand to set registration_ids field.
func (m *Message) RegIDs(tokens []string) {
	m.RegistrationIDs = tokens
}

// KeyValue type allows add fields to FCM message.
type KeyValue map[string]interface{}

// Add adds a key-value inside KeyValue struct.
func (d KeyValue) Add(key string, value interface{}) {
	d[key] = value
}

// Set sets a key-value inside KeyValue struct.
func (d KeyValue) Set(key string, value interface{}) {
	d[key] = value
}

// Get returns valueGet key inside KeyValue struct.
func (d KeyValue) Get(key string) interface{} {
	return d[key]
}

// Response FCM response.
//
// see: https://firebase.google.com/docs/cloud-messaging/http-server-ref#table5
type Response struct {
	MulticastID  int       `json:"multicast_id"`
	Success      int       `json:"success"`
	Failure      int       `json:"failure"`
	CanonicalIDs int       `json:"canonical_ids"`
	Results      []*Result `json:"results"`
}

// Result results inside FCM response. see Response struct
//
// see: https://firebase.google.com/docs/cloud-messaging/http-server-ref#table5
type Result struct {
	MessageID      string `json:"message_id"`
	RegistrationID string `json:"registration_id"`
	Error          string `json:"error"`
}
