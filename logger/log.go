// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package logger

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func NewLogger() (l *log.Logger) {
	prog := filepath.Base(os.Args[0])
	pid := strconv.Itoa(os.Getpid())

	l = log.New(os.Stderr, prog+"["+pid+"] ", log.LstdFlags)
	return
}
