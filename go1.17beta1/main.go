// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.17beta1 command runs the go command from Go 1.17beta1.
//
// To install, run:
//
//     $ go get golang.org/dl/go1.17beta1
//     $ go1.17beta1 download
//
// And then use the go1.17beta1 command as if it were your normal go
// command.
//
// See the release notes at https://tip.golang.org/doc/go1.17
//
// File bugs at https://golang.org/issues/new
package main

import "golang.org/dl/internal/version"

func main() {
	version.Run("go1.17beta1")
}
