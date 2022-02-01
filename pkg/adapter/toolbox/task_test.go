package toolbox

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
	"sync"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	defer ClearTask()

	tk := NewTask("taska", "0/30 * * * * *", func() error { fmt.Println("hello world"); return nil })
	err := tk.Run()
	if err != nil {
		t.Fatal(err)
	}
	AddTask("taska", tk)
	StartTask()
	time.Sleep(6 * time.Second)
	StopTask()
}

func TestSpec(t *testing.T) {
	defer ClearTask()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	tk1 := NewTask("tk1", "0 12 * * * *", func() error { fmt.Println("tk1"); return nil })
	tk2 := NewTask("tk2", "0,10,20 * * * * *", func() error { fmt.Println("tk2"); wg.Done(); return nil })
	tk3 := NewTask("tk3", "0 10 * * * *", func() error { fmt.Println("tk3"); wg.Done(); return nil })

	AddTask("tk1", tk1)
	AddTask("tk2", tk2)
	AddTask("tk3", tk3)
	StartTask()
	defer StopTask()

	select {
	case <-time.After(200 * time.Second):
		t.FailNow()
	case <-wait(wg):
	}
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
