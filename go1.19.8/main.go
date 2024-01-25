// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.19.8 command runs the go command from Go 1.19.8.
//
// To install, run:
//
//	$ go install github.com/btwiuse/dl/go1.19.8@latest
//	$ go1.19.8 download
//
// And then use the go1.19.8 command as if it were your normal go
// command.
//
// See the release notes at https://go.dev/doc/devel/release#go1.19.minor.
//
// File bugs at https://go.dev/issue/new.
package main

import "github.com/btwiuse/dl/version"

func main() {
	version.Run("go1.19.8")
}
