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
	"time"

	web "github.com/bhojpur/web/pkg/engine"
)

// Statistics struct
type Statistics web.Statistics

// URLMap contains several statistics struct to log different data
type URLMap web.URLMap

// AddStatistics add statistics task.
// it needs request method, request url, request controller and statistics time duration
func (m *URLMap) AddStatistics(requestMethod, requestURL, requestController string, requesttime time.Duration) {
	(*web.URLMap)(m).AddStatistics(requestMethod, requestURL, requestController, requesttime)
}

// GetMap put url statistics result in io.Writer
func (m *URLMap) GetMap() map[string]interface{} {
	return (*web.URLMap)(m).GetMap()
}

// GetMapData return all mapdata
func (m *URLMap) GetMapData() []map[string]interface{} {
	return (*web.URLMap)(m).GetMapData()
}

// StatisticsMap hosld global statistics data map
var StatisticsMap *URLMap

func init() {
	StatisticsMap = (*URLMap)(web.StatisticsMap)
}
