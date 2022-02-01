package engine

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
	"encoding/json"
	"testing"
	"time"
)

func TestStatics(t *testing.T) {
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(2000))
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(120000))
	StatisticsMap.AddStatistics("GET", "/api/user", "&admin.user", time.Duration(13000))
	StatisticsMap.AddStatistics("POST", "/api/admin", "&admin.user", time.Duration(14000))
	StatisticsMap.AddStatistics("POST", "/api/user/bhojpur", "&admin.user", time.Duration(12000))
	StatisticsMap.AddStatistics("POST", "/api/user/yunica", "&admin.user", time.Duration(13000))
	StatisticsMap.AddStatistics("DELETE", "/api/user", "&admin.user", time.Duration(1400))
	t.Log(StatisticsMap.GetMap())

	data := StatisticsMap.GetMapData()
	b, err := json.Marshal(data)
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Log(string(b))
}
