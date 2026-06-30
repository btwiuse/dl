// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The gotip command compiles and runs the go command from the development tree.
//
// To install, run:
//
//	$ go install golang.org/dl/gotip@latest
//	$ gotip download
//
// And then use the gotip command as if it were your normal go command.
//
// To update, run "gotip download" again.
//
// By default, gotip force sets GOTOOLCHAIN=auto to avoid the GOTOOLCHAIN value
// from go env -w interfere. Users can override this behavior by setting
// GOTOOLCHAIN in their environment var setting.
//
// The version tag can be overridden by setting the GOTIP environment variable.
// For example: GOTIP=go1.28.0-go4js.1 gotip download
package main

import (
	"cmp"
	"os"

	"golang.org/dl/internal/version"
)

var (
	GOTIP  = cmp.Or(os.Getenv("GOTIP"), "go1.27.0-gotip.1")
	NOCORS = "https://no-cors.up.railway.app/"
)

func main() {
	version.RunCustomSource(GOTIP, func(v, goos, arch string) string {
		return NOCORS + "https://github.com/justwasm/go4js/releases/download/" + v + "/" + v + "." + goos + "-" + arch + ".src.min.tar.gz"
	})
}
