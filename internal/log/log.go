// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package log

import (
	"log"
	"os"
)

type Verbosity int

const (
	VerbosityVerbose Verbosity = iota + 1
	VerbosityDebug
)

var (
	level Verbosity
)

func init() {
	log.SetOutput(os.Stderr)
}

func SetLevel(l Verbosity) {
	level = l
}

func Verbose(format string, v ...any) {
	if level >= VerbosityVerbose {
		log.Printf(format, v...)
	}
}

func Debug(format string, v ...any) {
	if level >= VerbosityDebug {
		log.Printf(format, v...)
	}
}

func Fatal(format string, v ...any) {
	log.Fatalf(format, v...)
}
