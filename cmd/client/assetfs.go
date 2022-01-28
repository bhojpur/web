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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/bhojpur/web/pkg/synthesis"
	"github.com/spf13/cobra"
)

// generateCmd represents the assetfs command
var assetfsCmd = &cobra.Command{
	Use:   "assetfs",
	Short: "To generate code for embedded assets using template files",
	Run: func(cmd *cobra.Command, args []string) {
		// Create input configurations.
		generateCmdOpts.Input = make([]synthesis.InputConfig, len(args))
		for i := range generateCmdOpts.Input {
			generateCmdOpts.Input[i] = parseInput(args[i])
		}
		err := synthesis.Translate(generateCmdOpts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "webctl assetfs: %v\n", err)
			os.Exit(1)
		}
		// assetFS generated output files
		out, in, err := getBhojpurFile()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot create temporary file", err)
			os.Exit(1)
		}
		debug := isDebug()
		r := bufio.NewReader(in)
		done := false
		for line, isPrefix, err := r.ReadLine(); err == nil; line, isPrefix, err = r.ReadLine() {
			if !isPrefix {
				line = append(line, '\n')
			}
			if _, err := out.Write(line); err != nil {
				fmt.Fprintln(os.Stderr, "Cannot write to ", out.Name(), err)
				return
			}
			if !done && !isPrefix && bytes.HasPrefix(line, []byte("import (")) {
				if debug {
					fmt.Fprintln(out, "\t\"net/http\"")
				} else {
					fmt.Fprintln(out, "\t\"github.com/bhojpur/web/pkg/synthesis\"")
				}
				done = true
			}
		}
		if debug {
			fmt.Fprintln(out, `
func assetFS() http.FileSystem {
	for k := range _bhojpurTree.Children {
		return http.Dir(k)
	}
	panic("unreachable")
}`)
		} else {
			fmt.Fprintln(out, `
func assetFS() *synthesis.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bhojpurTree.Children {
		return &synthesis.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}`)
		}
		// Close files BEFORE remove calls (don't use defer).
		in.Close()
		out.Close()
		if err := os.Remove(in.Name()); err != nil {
			fmt.Fprintln(os.Stderr, "Cannot remove", in.Name(), err)
		}
	},
}

func init() {
	rootCmd.AddCommand(assetfsCmd)

	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.Debug, "debug", generateCmdOpts.Debug, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.Dev, "dev", generateCmdOpts.Dev, "Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated source code's package.")
	assetfsCmd.PersistentFlags().StringVar(&generateCmdOpts.Tags, "tags", generateCmdOpts.Tags, "Optional set of build tags to include.")
	assetfsCmd.PersistentFlags().StringVar(&generateCmdOpts.Prefix, "prefix", generateCmdOpts.Prefix, "Optional path prefix to strip off asset names.")
	assetfsCmd.PersistentFlags().StringVar(&generateCmdOpts.Package, "pkg", generateCmdOpts.Package, "Package name to use in the generated code.")
	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoMemCopy, "nomemcopy", generateCmdOpts.NoMemCopy, "Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.")
	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoCompress, "nocompress", generateCmdOpts.NoCompress, "Assets will *not* be GZIP compressed when this flag is specified.")
	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.NoMetadata, "nometadata", generateCmdOpts.NoMetadata, "Assets will not preserve size, mode, and modtime info.")
	assetfsCmd.PersistentFlags().BoolVar(&generateCmdOpts.HttpFileSystem, "fs", generateCmdOpts.HttpFileSystem, "Whether generate instance http.FileSystem interface code.")
	assetfsCmd.PersistentFlags().UintVar(&generateCmdOpts.Mode, "mode", generateCmdOpts.Mode, "Optional file mode override for all files.")
	assetfsCmd.PersistentFlags().Int64Var(&generateCmdOpts.ModTime, "modtime", generateCmdOpts.ModTime, "Optional modification unix timestamp override for all files.")
	assetfsCmd.PersistentFlags().StringVar(&generateCmdOpts.Output, "o", generateCmdOpts.Output, "Optional name of the output file to be generated.")
}

func isDebug() bool {
	return generateCmdOpts.Debug
}

func getOutputFile() string {
	return generateCmdOpts.Output
}

func getBhojpurFile() (*os.File, *os.File, error) {
	var defaultFolder = "./"
	var defaultName = defaultFolder + "bhojpur.go"
	outputLoc := defaultFolder + getOutputFile()

	_, err := copyBhojpurFile(outputLoc, defaultName)
	if err != nil {
		fmt.Errorf("copying file", outputLoc, defaultName)
		return &os.File{}, &os.File{}, err
	}

	tempFile, err := os.Open((defaultName))
	if err != nil {
		fmt.Errorf("opening file", defaultName)
		return &os.File{}, &os.File{}, err
	}

	outputFile, err := os.Create(outputLoc)
	if err != nil {
		fmt.Errorf("creating file", outputLoc)
		return &os.File{}, &os.File{}, err
	}

	return outputFile, tempFile, nil
}

func copyBhojpurFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		fmt.Errorf("checking src file", src)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		fmt.Errorf("opening src file", src)
		return 0, err
	}

	destination, err := os.Create(dst)
	if err != nil {
		fmt.Errorf("creating dst file", src)
		return 0, err
	}
	nBytes, err := io.Copy(destination, source)

	// Close files (don't use defer).
	source.Close()
	destination.Close()

	return nBytes, err
}
