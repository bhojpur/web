package common

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
	webapp "github.com/bhojpur/web/pkg/app"
	webui "github.com/bhojpur/web/pkg/app/ui"
)

const (
	githubURL        = "https://github.com/bhojpur/web"
	githubSponsorURL = "https://github.com/shashi-rai"
)

type githubSponsor struct {
	webapp.Compo

	Iclass string
}

func newGithubSponsor() *githubSponsor {
	return &githubSponsor{}
}

func (s *githubSponsor) Class(v string) *githubSponsor {
	s.Iclass = webapp.AppendClass(s.Iclass, v)
	return s
}

func (s *githubSponsor) Render() webapp.UI {
	return webui.Stack().
		Class(s.Iclass).
		Center().
		Middle().
		Content(
			webapp.Aside().
				Class("magnify").
				Class("text-center").
				Body(
					webapp.A().
						Class("default").
						Href(githubSponsorURL).
						Body(
							webui.Icon().
								Class("center").
								Class("icon-top").
								Size(72).
								Src(githubSVG),
							webapp.Header().
								Class("h3").
								Class("default").
								Text("Support on GitHub"),
							webapp.P().
								Class("subtext").
								Text("Help framework development by sponsoring it on GitHub"),
						),
				),
		)
}
