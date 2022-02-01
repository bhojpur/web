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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type countTask struct {
	cnt     int
	mockErr error
}

func (c *countTask) GetSpec(ctx context.Context) string {
	return "AAA"
}

func (c *countTask) GetStatus(ctx context.Context) string {
	return "SUCCESS"
}

func (c *countTask) Run(ctx context.Context) error {
	c.cnt++
	return c.mockErr
}

func (c *countTask) SetNext(ctx context.Context, time time.Time) {
}

func (c *countTask) GetNext(ctx context.Context) time.Time {
	return time.Now()
}

func (c *countTask) SetPrev(ctx context.Context, time time.Time) {
}

func (c *countTask) GetPrev(ctx context.Context) time.Time {
	return time.Now()
}

func TestRunTaskCommand_Execute(t *testing.T) {
	task := &countTask{}
	AddTask("count", task)

	cmd := &runTaskCommand{}

	res := cmd.Execute()
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "task name not passed", res.Error.Error())

	res = cmd.Execute(10)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "parameter is invalid", res.Error.Error())

	res = cmd.Execute("CCCC")
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "task with name CCCC not found", res.Error.Error())

	res = cmd.Execute("count")
	assert.NotNil(t, res)
	assert.True(t, res.IsSuccess())

	task.mockErr = errors.New("mock error")
	res = cmd.Execute("count")
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "mock error", res.Error.Error())
}

func TestListTaskCommand_Execute(t *testing.T) {
	task := &countTask{}

	cmd := &listTaskCommand{}

	res := cmd.Execute()

	assert.True(t, res.IsSuccess())

	_, ok := res.Content.([][]string)
	assert.True(t, ok)

	AddTask("count", task)

	res = cmd.Execute()

	assert.True(t, res.IsSuccess())

	rl, ok := res.Content.([][]string)
	assert.True(t, ok)
	assert.Equal(t, 1, len(rl))
}
