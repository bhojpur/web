package cmd

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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bhojpur/web/pkg/synthesis"
	"github.com/spf13/cobra"
)

var generateCmdOpts = synthesis.NewConfig()

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "To generate source code of web applications in specific language",
	Run: func(cmd *cobra.Command, args []string) {
		// Create input configurations.
		generateCmdOpts.Input = make([]synthesis.InputConfig, len(args))
		for i := range generateCmdOpts.Input {
			generateCmdOpts.Input[i] = parseInput(args[i])
		}
		// generate Go source code from template files
		err := synthesis.Translate(generateCmdOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "webctl generate: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.Debug, "debug", generateCmdOpts.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.Dev, "dev", generateCmdOpts.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated source code's package.")
	generateCmd.PersistentFlags().StringVar(&generateCmdOpts.Tags, "tags", generateCmdOpts.Tags, "Optional set of build tags to include.")
	generateCmd.PersistentFlags().StringVar(&generateCmdOpts.Prefix, "prefix", generateCmdOpts.Prefix, "Optional path prefix to strip off asset names.")
	generateCmd.PersistentFlags().StringVar(&generateCmdOpts.Package, "pkg", generateCmdOpts.Package, "Package name to use in the generated code.")
	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoMemCopy, "nomemcopy", generateCmdOpts.NoMemCopy, "Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoCompress, "nocompress", generateCmdOpts.NoCompress, "Assets will *not* be GZIP compressed when this flag is specified.")
	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoMetadata, "nometadata", generateCmdOpts.NoMetadata, "Assets will not preserve size, mode, and modtime info.")
	generateCmd.PersistentFlags().BoolVar(&generateCmdOpts.HttpFileSystem, "fs", generateCmdOpts.HttpFileSystem, "Whether generate instance http.FileSystem interface code.")
	generateCmd.PersistentFlags().UintVar(&generateCmdOpts.Mode, "mode", generateCmdOpts.Mode, "Optional file mode override for all files.")
	generateCmd.PersistentFlags().Int64Var(&generateCmdOpts.ModTime, "modtime", generateCmdOpts.ModTime, "Optional modification unix timestamp override for all files.")
	generateCmd.PersistentFlags().StringVar(&generateCmdOpts.Output, "o", generateCmdOpts.Output, "Optional name of the output file to be generated.")

	ignore := make([]string, 0)
	flag.Var((*AppendSliceValue)(&ignore), "ignore", "Regex pattern to ignore")
	flag.Parse()

	patterns := make([]*regexp.Regexp, 0)
	for _, pattern := range ignore {
		patterns = append(patterns, regexp.MustCompile(pattern))
	}
	generateCmdOpts.Ignore = patterns
}

// parseRecursive determines whether the given path has a recursive indicator and
// returns a new path with recursive indicator chopped off if it does. For example
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(path string) synthesis.InputConfig {
	if strings.HasSuffix(path, "/...") {
		return synthesis.InputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	} else {
		return synthesis.InputConfig{
			Path:      filepath.Clean(path),
			Recursive: false,
		}
	}

}

type AppendSliceValue []string

func (s *AppendSliceValue) String() string {
	return strings.Join(*s, ",")
}

func (s *AppendSliceValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	*s = append(*s, value)
	return nil
}
