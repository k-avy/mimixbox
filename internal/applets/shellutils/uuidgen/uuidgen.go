//
// mimixbox/internal/applets/shellutils/uuidgen/uuidgen.go
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
package uuidgen

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	mb "github.com/nao1215/mimixbox/internal/lib"
)

const cmdName string = "uuidgen"

const version = "1.0.0"

var osExit = os.Exit

type options struct {
	Version bool `short:"v" long:"version" description:"Show uuidgen command version"`
}

func Run() (int, error) {
	var opts options
	var err error

	if _, err = parseArgs(&opts); err != nil {
		return mb.ExitFailure, nil
	}
	return uuidgen()
}

func parseArgs(opts *options) ([]string, error) {
	p := initParser(opts)

	args, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if opts.Version {
		mb.ShowVersion(cmdName, version)
		osExit(mb.ExitSuccess)
	}
	return args, nil
}

func initParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS] "

	return parser
}

// Generate UUID version.4
// [e.g. ver 1.0]
//  xxxxxxxx-xxxx-1xxx-xxxx-xxxxxxxxxxxx
//                ^
//                |
//              always "1"
// [e.g. ver 4.0]
//  xxxxxxxx-xxxx-4xxx-xxxx-xxxxxxxxxxxx
//                ^
//                |
//              always "4"
func uuidgen() (int, error) {
	fp, err := os.Open("/proc/sys/kernel/random/uuid")
	if err != nil {
		return mb.ExitFailure, err
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 38) // UUID's format sample: 333f899e-5dbf-41ec-8c42-e3b3ddbe2e68
	line, _, err := reader.ReadLine()
	if err != nil {
		return mb.ExitFailure, err
	}

	fmt.Fprintln(os.Stdout, string(line))
	return mb.ExitSuccess, nil
}
