// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.9beta2 command runs the go command from Go 1.9beta2.
//
// To install, run:
//
//	$ go install github.com/btwiuse/dl/go1.9beta2@latest
//	$ go1.9beta2 download
//
// And then use the go1.9beta2 command as if it were your normal go
// command.
//
// See the release notes at https://tip.golang.org/doc/go1.9
//
// File bugs at https://golang.org/issues/new
package main

import "github.com/btwiuse/dl/version"

func main() {
	version.Run("go1.9beta2")
}
