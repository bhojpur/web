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
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/bhojpur/web/pkg/core/utils"
)

// Statistics struct
type Statistics struct {
	RequestURL        string
	RequestController string
	RequestNum        int64
	MinTime           time.Duration
	MaxTime           time.Duration
	TotalTime         time.Duration
}

// URLMap contains several statistics struct to log different data
type URLMap struct {
	lock        sync.RWMutex
	LengthLimit int // limit the urlmap's length if it's equal to 0 there's no limit
	urlmap      map[string]map[string]*Statistics
}

// AddStatistics add statistics task.
// it needs request method, request url, request controller and statistics time duration
func (m *URLMap) AddStatistics(requestMethod, requestURL, requestController string, requesttime time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if method, ok := m.urlmap[requestURL]; ok {
		if s, ok := method[requestMethod]; ok {
			s.RequestNum++
			if s.MaxTime < requesttime {
				s.MaxTime = requesttime
			}
			if s.MinTime > requesttime {
				s.MinTime = requesttime
			}
			s.TotalTime += requesttime
		} else {
			nb := &Statistics{
				RequestURL:        requestURL,
				RequestController: requestController,
				RequestNum:        1,
				MinTime:           requesttime,
				MaxTime:           requesttime,
				TotalTime:         requesttime,
			}
			m.urlmap[requestURL][requestMethod] = nb
		}
	} else {
		if m.LengthLimit > 0 && m.LengthLimit <= len(m.urlmap) {
			return
		}
		methodmap := make(map[string]*Statistics)
		nb := &Statistics{
			RequestURL:        requestURL,
			RequestController: requestController,
			RequestNum:        1,
			MinTime:           requesttime,
			MaxTime:           requesttime,
			TotalTime:         requesttime,
		}
		methodmap[requestMethod] = nb
		m.urlmap[requestURL] = methodmap
	}
}

// GetMap put url statistics result in io.Writer
func (m *URLMap) GetMap() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	fields := []string{"requestUrl", "method", "times", "used", "max used", "min used", "avg used"}

	var resultLists [][]string
	content := make(map[string]interface{})
	content["Fields"] = fields

	for k, v := range m.urlmap {
		for kk, vv := range v {
			result := []string{
				fmt.Sprintf("% -50s", template.HTMLEscapeString(k)),
				fmt.Sprintf("% -10s", kk),
				fmt.Sprintf("% -16d", vv.RequestNum),
				fmt.Sprintf("%d", vv.TotalTime),
				fmt.Sprintf("% -16s", utils.ToShortTimeFormat(vv.TotalTime)),
				fmt.Sprintf("%d", vv.MaxTime),
				fmt.Sprintf("% -16s", utils.ToShortTimeFormat(vv.MaxTime)),
				fmt.Sprintf("%d", vv.MinTime),
				fmt.Sprintf("% -16s", utils.ToShortTimeFormat(vv.MinTime)),
				fmt.Sprintf("%d", time.Duration(int64(vv.TotalTime)/vv.RequestNum)),
				fmt.Sprintf("% -16s", utils.ToShortTimeFormat(time.Duration(int64(vv.TotalTime)/vv.RequestNum))),
			}
			resultLists = append(resultLists, result)
		}
	}
	content["Data"] = resultLists
	return content
}

// GetMapData return all mapdata
func (m *URLMap) GetMapData() []map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var resultLists []map[string]interface{}

	for k, v := range m.urlmap {
		for kk, vv := range v {
			result := map[string]interface{}{
				"request_url": k,
				"method":      kk,
				"times":       vv.RequestNum,
				"total_time":  utils.ToShortTimeFormat(vv.TotalTime),
				"max_time":    utils.ToShortTimeFormat(vv.MaxTime),
				"min_time":    utils.ToShortTimeFormat(vv.MinTime),
				"avg_time":    utils.ToShortTimeFormat(time.Duration(int64(vv.TotalTime) / vv.RequestNum)),
			}
			resultLists = append(resultLists, result)
		}
	}
	return resultLists
}

// StatisticsMap hosld global statistics data map
var StatisticsMap *URLMap

func init() {
	StatisticsMap = &URLMap{
		urlmap: make(map[string]map[string]*Statistics),
	}
}
