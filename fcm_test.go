// Package fcm contains Firebase Cloud Messaging client.
// Only http requests are implemented, if xmpp is required
// open an issue.
/*
{
	"multicast_id": 216,
	"success": 3,
	"failure": 3,
	"canonical_ids": 1,
	"results": [
		{ "message_id": "1:0408" },
		{ "error": "Unavailable" },
		{ "error": "InvalidRegistration" },
		{ "message_id": "1:1516" },
		{ "message_id": "1:2342", "registration_id": "32" },
		{ "error": "NotRegistered"}
	]
}
*/
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
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client, err := New("someServerKey12345678901234567890123456", time.Second)
	if err != nil {
		t.Logf("err : [%s]", err)
		t.Fail()
		return
	}

	msg := NewMessage("TOKENFAIL")
	_, err = client.Send(msg)
	if err == nil {
		t.Logf("err : [%s]", err)
		t.Fail()
	}
}

type T struct {
	Purpose          string
	Input            *Message
	ExpectedError    error
	ExpectedResponse *Response
}

// TestTable tests most common errors.
// TODO; add cases.
func TestTable(t *testing.T) {
	client, err := New("someServerKey12345678901234567890123456", time.Second)
	if err != nil {
		t.Logf("err : [%s]", err)
		t.Fail()
		return
	}
	table := []T{
		T{
			Purpose: "Demonstrate some case",
			Input: &Message{
				To: "someDeviceToken",
			},
			ExpectedError:    ErrRequestFail,
			ExpectedResponse: nil,
		},
	}
	for i := range table {
		x := table[i]
		actual, err := client.Send(x.Input)
		if err != x.ExpectedError {
			t.Logf("error : expected [%v] actual [%v]", x.ExpectedError, err)
			t.Fail()
			continue
		}
		if actual != x.ExpectedResponse {
			t.Logf("response : expected [%v] actual [%v]", x.ExpectedResponse, actual)
			t.Fail()
			continue
		}
	}
}
