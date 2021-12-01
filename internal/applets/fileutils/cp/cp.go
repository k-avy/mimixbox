//
// mimixbox/internal/applets/fileutils/cp/cp.go
//
// Copyright 2021 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cp

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	mb "github.com/nao1215/mimixbox/internal/lib"

	"github.com/jessevdk/go-flags"
)

const cmdName string = "cp"

const version = "1.0.3"

var osExit = os.Exit

// Exit code
const (
	ExitSuccess int = iota // 0
	ExitFailuer
)

type options struct {
	Force       bool `short:"f" long:"force" description:"If file exists, forcibly overwrite it"`
	Interactive bool `short:"i" long:"interactive" description:"Ask every time if you want to remove"`
	Recursive   bool `short:"r" long:"recursive" description:"Recursively copy directories"`
	Version     bool `short:"v" long:"version" description:"Show cp command version"`
}

func Run() (int, error) {
	var opts options
	var args []string
	var err error

	if args, err = parseArgs(&opts); err != nil {
		return ExitFailuer, nil
	}

	err = cp(args, opts)
	if err != nil {
		return ExitFailuer, err
	}

	return ExitSuccess, nil
}

func cp(files []string, opts options) error {
	dest := os.ExpandEnv(files[len(files)-1])

	for _, src := range files[:len(files)-1] {
		s := os.ExpandEnv(src)
		if !mb.Exists(s) {
			return errors.New(s + " does not exist")
		}

		if mb.IsSamePath(s, dest) {
			return errors.New(s + " and " + dest + " is same.")
		}

		if mb.IsFile(s) {
			if err := cpFile(s, dest, opts); err != nil {
				return err
			}
		} else {
			if err := cpDir(s, dest, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

func cpFile(src string, dest string, opts options) error {
	if mb.IsFile(dest) && mb.IsSameFileName(src, dest) && opts.Interactive {
		if !mb.Question("Overwrite " + dest) {
			return nil // Skip this file
		}
	}
	if mb.IsDir(dest) {
		dest = filepath.Join(dest, path.Base(src))
	}

	return mb.Copy(src, dest)
}

func cpDir(src string, dest string, opts options) error {
	if !opts.Recursive {
		return errors.New("--recursive is not specified: omitting directory: " + src)
	}

	srcDirs, srcFiles, err := mb.Walk(src)
	if err != nil {
		return err
	}

	// Make destination directory
	for _, dir := range srcDirs {
		dir = path.Join(dest, dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	for _, src := range srcFiles {
		destFile := path.Join(dest, src)
		err := mb.Copy(src, destFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseArgs(opts *options) ([]string, error) {
	p := initParser(opts)

	args, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if opts.Version {
		mb.ShowVersion(cmdName, version)
		osExit(ExitSuccess)
	}

	if !isValidArgNr(args) {
		showHelp(p)
		osExit(ExitFailuer)
	}
	return args, nil
}

func initParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS] SOURCE DESTINATION"

	return parser
}

func isValidArgNr(args []string) bool {
	return len(args) >= 2
}

func showHelp(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
}
