package task

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
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	tk := NewTask("taska", "0/30 * * * * *", func(ctx context.Context) error {
		fmt.Println("hello world")
		return nil
	})
	err := tk.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
	m.AddTask("taska", tk)
	m.StartTask()
	time.Sleep(3 * time.Second)
	m.StopTask()
}

func TestModifyTaskListAfterRunning(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	tk := NewTask("taskb", "0/30 * * * * *", func(ctx context.Context) error {
		fmt.Println("hello world")
		return nil
	})
	err := tk.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
	m.AddTask("taskb", tk)
	m.StartTask()
	go func() {
		m.DeleteTask("taskb")
	}()
	go func() {
		m.AddTask("taskb1", tk)
	}()

	time.Sleep(3 * time.Second)
	m.StopTask()
}

func TestSpec(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tk1 := NewTask("tk1", "0 12 * * * *", func(ctx context.Context) error { fmt.Println("tk1"); return nil })
	tk2 := NewTask("tk2", "0,10,20 * * * * *", func(ctx context.Context) error { fmt.Println("tk2"); wg.Done(); return nil })
	tk3 := NewTask("tk3", "0 10 * * * *", func(ctx context.Context) error { fmt.Println("tk3"); wg.Done(); return nil })

	m.AddTask("tk1", tk1)
	m.AddTask("tk2", tk2)
	m.AddTask("tk3", tk3)
	m.StartTask()
	defer m.StopTask()

	select {
	case <-time.After(200 * time.Second):
		t.FailNow()
	case <-wait(wg):
	}
}

func TestTask_Run(t *testing.T) {
	cnt := -1
	task := func(ctx context.Context) error {
		cnt++
		fmt.Printf("Hello, world! %d \n", cnt)
		return errors.New(fmt.Sprintf("Hello, world! %d", cnt))
	}
	tk := NewTask("taska", "0/30 * * * * *", task)
	for i := 0; i < 200; i++ {
		e := tk.Run(nil)
		assert.NotNil(t, e)
	}

	l := tk.Errlist
	assert.Equal(t, 100, len(l))
	assert.Equal(t, "Hello, world! 100", l[0].errinfo)
	assert.Equal(t, "Hello, world! 101", l[1].errinfo)
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
