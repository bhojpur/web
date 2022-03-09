package dockerize

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
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
)

const dockerBuildTemplate = `FROM {{.BaseImage}}

# Godep for vendoring
RUN go get github.com/tools/godep

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR $GOPATH{{.Appdir}}
RUN mkdir -p $APP_DIR

# Set the entrypoint
ENTRYPOINT (cd $APP_DIR && ./{{.Entrypoint}})
ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 godep go build -ldflags '-d -w -s'

EXPOSE {{.Expose}}
`

// Dockerfile holds the information about the Docker container.
type Dockerfile struct {
	BaseImage  string
	Appdir     string
	Entrypoint string
	Expose     string
}

var CmdDockerize = &commands.Command{
	CustomFlags: true,
	UsageLine:   "dockerize",
	Short:       "Generates a Dockerfile for your Bhojpur Web application",
	Long: `Dockerize generates a Dockerfile for your Bhojpur Web application.
  The Dockerfile will compile, get the dependencies with {{"godep"|bold}}, and set the entrypoint.

  {{"Example:"|bold}}
    $ webutl dockerize -expose="3000,80,25"
  `,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    dockerizeApp,
}

var (
	expose    string
	baseImage string
)

func init() {
	fs := flag.NewFlagSet("dockerize", flag.ContinueOnError)
	fs.StringVar(&baseImage, "image", "library/golang", "Set the base image of the Docker container.")
	fs.StringVar(&expose, "expose", "8080", "Port(s) to expose in the Docker container.")
	CmdDockerize.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdDockerize)
}

func dockerizeApp(cmd *commands.Command, args []string) int {
	if err := cmd.Flag.Parse(args); err != nil {
		cliLogger.Log.Fatalf("Error parsing flags: %v", err.Error())
	}

	cliLogger.Log.Info("Generating Dockerfile...")

	gopath := os.Getenv("GOPATH")
	dir, err := filepath.Abs(".")
	if err != nil {
		cliLogger.Log.Error(err.Error())
	}

	appdir := strings.Replace(dir, gopath, "", 1)

	// In case of multiple ports to expose inside the container,
	// replace all the commas with whitespaces.
	// See the verb EXPOSE in the Docker documentation.
	expose = strings.Replace(expose, ",", " ", -1)

	_, entrypoint := path.Split(appdir)
	dockerfile := Dockerfile{
		BaseImage:  baseImage,
		Appdir:     appdir,
		Entrypoint: entrypoint,
		Expose:     expose,
	}

	generateDockerfile(dockerfile)
	return 0
}

func generateDockerfile(df Dockerfile) {
	t := template.Must(template.New("dockerBuildTemplate").Parse(dockerBuildTemplate)).Funcs(utils.BhojpurWebFuncMap())

	f, err := os.Create("Dockerfile")
	if err != nil {
		cliLogger.Log.Fatalf("Error writing Dockerfile: %v", err.Error())
	}
	defer utils.CloseFile(f)

	t.Execute(f, df)

	cliLogger.Log.Success("Dockerfile generated.")
}
