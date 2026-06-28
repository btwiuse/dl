// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go4js command runs the go command from the justwasm/go4js project.
//
// To install, run:
//
//	$ go install golang.org/dl/go4js@latest
//	$ go4js download
//
// And then use the go4js command as if it were your normal go
// command.
package main

import "golang.org/dl/internal/version"

func main() {
	version.RunCustom("go1.27.0-go4js.1", func(v, goos, arch string) string {
		return "https://github.com/justwasm/go4js/releases/download/" + v + "/" + v + "." + goos + "-" + arch + ".tar.gz"
	})
}
