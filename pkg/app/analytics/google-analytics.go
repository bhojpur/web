package analytics

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

	"github.com/bhojpur/web/pkg/app"
)

// GoogleAnalyticsHeader returns the header to use in the app.Handler.RawHeader
// field to initialize Google Analytics.
func GoogleAnalyticsHeader(propertyID string) string {
	return fmt.Sprintf(`<!-- Global site tag (gtag.js) - Google Analytics -->
	<script defer src="https://www.googletagmanager.com/gtag/js?id=%s"></script>
	<script>
	  window.dataLayer = window.dataLayer || [];
	  function gtag(){dataLayer.push(arguments);}
	  gtag('js', new Date());
	
	  gtag('config', '%s', {'send_page_view': false});
	</script>`, propertyID, propertyID)
}

func NewGoogleAnalytics() Backend {
	return googleAnalytics{}
}

type googleAnalytics struct {
	propertyID string
}

func (a googleAnalytics) Identify(userID string, traits map[string]interface{}) {
	a.gtag("set", map[string]interface{}{
		"user_id": userID,
	})
}

func (a googleAnalytics) Track(event string, properties map[string]interface{}) {
	a.gtag("event", event, properties)
}

func (a googleAnalytics) Page(name string, properties map[string]interface{}) {
	a.gtag("event", "page_view", map[string]interface{}{
		"page_title":    properties["title"],
		"page_location": properties["url"],
		"page_path":     properties["path"],
	})
}

func (a googleAnalytics) gtag(args ...interface{}) {
	gtag := app.Window().Get("gtag")
	if !gtag.Truthy() {
		return
	}
	app.Window().Call("gtag", args...)
}
