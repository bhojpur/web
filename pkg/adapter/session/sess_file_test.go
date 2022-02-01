package session

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

const sid = "Session_id"
const sidNew = "Session_id_new"
const sessionPath = "./_session_runtime"

var (
	mutex sync.Mutex
)

func TestFileProvider_SessionExist(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	if fp.SessionExist(sid) {
		t.Error()
	}

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}
}

func TestFileProvider_SessionExist2(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	if fp.SessionExist(sid) {
		t.Error()
	}

	if fp.SessionExist("") {
		t.Error()
	}

	if fp.SessionExist("1") {
		t.Error()
	}
}

func TestFileProvider_SessionRead(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	s, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	_ = s.Set("sessionValue", 18975)
	v := s.Get("sessionValue")

	if v.(int) != 18975 {
		t.Error()
	}
}

func TestFileProvider_SessionRead1(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead("")
	if err == nil {
		t.Error(err)
	}

	_, err = fp.SessionRead("1")
	if err == nil {
		t.Error(err)
	}
}

func TestFileProvider_SessionAll(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 546

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	if fp.SessionAll() != sessionCount {
		t.Error()
	}
}

func TestFileProvider_SessionRegenerate(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}

	_, err = fp.SessionRegenerate(sid, sidNew)
	if err != nil {
		t.Error(err)
	}

	if fp.SessionExist(sid) {
		t.Error()
	}

	if !fp.SessionExist(sidNew) {
		t.Error()
	}
}

func TestFileProvider_SessionDestroy(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	_, err := fp.SessionRead(sid)
	if err != nil {
		t.Error(err)
	}

	if !fp.SessionExist(sid) {
		t.Error()
	}

	err = fp.SessionDestroy(sid)
	if err != nil {
		t.Error(err)
	}

	if fp.SessionExist(sid) {
		t.Error()
	}
}

func TestFileProvider_SessionGC(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(1, sessionPath)

	sessionCount := 412

	for i := 1; i <= sessionCount; i++ {
		_, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
	}

	time.Sleep(2 * time.Second)

	fp.SessionGC()
	if fp.SessionAll() != 0 {
		t.Error()
	}
}

func TestFileSessionStore_Set(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		err := s.Set(i, i)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestFileSessionStore_Get(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(i, i)

		v := s.Get(i)
		if v.(int) != i {
			t.Error()
		}
	}
}

func TestFileSessionStore_Delete(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	s, _ := fp.SessionRead(sid)
	s.Set("1", 1)

	if s.Get("1") == nil {
		t.Error()
	}

	s.Delete("1")

	if s.Get("1") != nil {
		t.Error()
	}
}

func TestFileSessionStore_Flush(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 100
	s, _ := fp.SessionRead(sid)
	for i := 1; i <= sessionCount; i++ {
		_ = s.Set(i, i)
	}

	_ = s.Flush()

	for i := 1; i <= sessionCount; i++ {
		if s.Get(i) != nil {
			t.Error()
		}
	}
}

func TestFileSessionStore_SessionID(t *testing.T) {
	mutex.Lock()
	defer mutex.Unlock()
	os.RemoveAll(sessionPath)
	defer os.RemoveAll(sessionPath)
	fp := &FileProvider{}

	_ = fp.SessionInit(180, sessionPath)

	sessionCount := 85

	for i := 1; i <= sessionCount; i++ {
		s, err := fp.SessionRead(fmt.Sprintf("%s_%d", sid, i))
		if err != nil {
			t.Error(err)
		}
		if s.SessionID() != fmt.Sprintf("%s_%d", sid, i) {
			t.Error(err)
		}
	}
}
