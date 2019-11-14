// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The gotip command compiles and runs the go command from the development tree.
//
// To install, run:
//
//     $ go get golang.org/dl/gotip
//     $ gotip download
//
// And then use the gotip command as if it were your normal go command.
//
// To update, run "gotip download" again.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	log.SetFlags(0)

	root, err := goroot("gotip")
	if err != nil {
		log.Fatalf("gotip: %v", err)
	}

	if len(os.Args) == 2 && os.Args[1] == "download" {
		if err := installTip(root); err != nil {
			log.Fatalf("gotip: %v", err)
		}
		log.Printf("Success. You may now run 'gotip'!")
		os.Exit(0)
	}

	gobin := filepath.Join(root, "bin", "go"+exe())
	if _, err := os.Stat(gobin); err != nil {
		log.Fatalf("gotip: not downloaded. Run 'gotip download' to install to %v", root)
	}

	cmd := exec.Command(gobin, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	newPath := filepath.Join(root, "bin")
	if p := os.Getenv("PATH"); p != "" {
		newPath += string(filepath.ListSeparator) + p
	}
	cmd.Env = dedupEnv(caseInsensitiveEnv, append(os.Environ(), "GOROOT="+root, "PATH="+newPath))
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// TODO: return the same exit status maybe.
			os.Exit(1)
		}
		log.Fatalf("gotip: failed to execute %v: %v", gobin, err)
	}
	os.Exit(0)
}

func installTip(root string) error {
	git := func(args ...string) error {
		cmd := exec.Command("git", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = root
		return cmd.Run()
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		if err := os.MkdirAll(root, 0755); err != nil {
			return fmt.Errorf("failed to create repository: %v", err)
		}
		if err := git("clone", "--depth=1", "https://go.googlesource.com/go", root); err != nil {
			return fmt.Errorf("failed to clone git repository: %v", err)
		}
	} else {
		log.Printf("Updating the go development tree...")
		if err := git("fetch", "origin"); err != nil {
			return fmt.Errorf("failed to fetch git repository updates: %v", err)
		}
	}

	// Use checkout and a detached HEAD, because it will refuse to overwrite
	// local changes, and warn if commits are being left behind, but will not
	// mind if master is force-pushed upstream.
	if err := git("-c", "advice.detachedHead=false", "checkout", "origin/master"); err != nil {
		return fmt.Errorf("failed to checkout git repository: %v", err)
	}
	// It shouldn't be the case, but in practice sometimes binary artifacts
	// generated by earlier Go versions interfere with the build.
	//
	// Ask the user what to do about them if they are not gitignored. They might
	// be artifacts that used to be ignored in previous versions, or precious
	// uncommitted source files.
	if err := git("clean", "-i", "-d"); err != nil {
		return fmt.Errorf("failed to cleanup git repository: %v", err)
	}
	// Wipe away probably boring ignored files without bothering the user.
	if err := git("clean", "-q", "-f", "-d", "-X"); err != nil {
		return fmt.Errorf("failed to cleanup git repository: %v", err)
	}

	cmd := exec.Command(filepath.Join(root, "src", makeScript()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(root, "src")
	if runtime.GOOS == "windows" {
		// Workaround make.bat not autodetecting GOROOT_BOOTSTRAP. Issue 28641.
		goroot, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			return fmt.Errorf("failed to detect an existing go installation for bootstrap: %v", err)
		}
		cmd.Env = append(os.Environ(), "GOROOT_BOOTSTRAP="+strings.TrimSpace(string(goroot)))
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build go: %v", err)
	}

	return nil
}

func makeScript() string {
	switch runtime.GOOS {
	case "plan9":
		return "make.rc"
	case "windows":
		return "make.bat"
	default:
		return "make.bash"
	}
}

const caseInsensitiveEnv = runtime.GOOS == "windows"

func exe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func goroot(version string) (string, error) {
	home, err := homedir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(home, "sdk", version), nil
}

func homedir() (string, error) {
	// This could be replaced with os.UserHomeDir, but it was introduced too
	// recently, and we want this to work with go as packaged by Linux
	// distributions. Note that user.Current is not enough as it does not
	// prioritize $HOME. See also Issue 26463.
	switch runtime.GOOS {
	case "plan9":
		return "", fmt.Errorf("%q not yet supported", runtime.GOOS)
	case "windows":
		if dir := os.Getenv("USERPROFILE"); dir != "" {
			return dir, nil
		}
		return "", errors.New("can't find user home directory; %USERPROFILE% is empty")
	default:
		if dir := os.Getenv("HOME"); dir != "" {
			return dir, nil
		}
		if u, err := user.Current(); err == nil && u.HomeDir != "" {
			return u.HomeDir, nil
		}
		return "", errors.New("can't find user home directory; $HOME is empty")
	}
}

// dedupEnv returns a copy of env with any duplicates removed, in favor of
// later values.
// Items are expected to be on the normal environment "key=value" form.
// If caseInsensitive is true, the case of keys is ignored.
//
// This function is unnecessary when the binary is
// built with Go 1.9+, but keep it around for now until Go 1.8
// is no longer seen in the wild in common distros.
//
// This is copied verbatim from golang.org/x/build/envutil.Dedup at CL 10301
// (commit a91ae26).
func dedupEnv(caseInsensitive bool, env []string) []string {
	out := make([]string, 0, len(env))
	saw := map[string]int{} // to index in the array
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 1 {
			out = append(out, kv)
			continue
		}
		k := kv[:eq]
		if caseInsensitive {
			k = strings.ToLower(k)
		}
		if dupIdx, isDup := saw[k]; isDup {
			out[dupIdx] = kv
		} else {
			saw[k] = len(out)
			out = append(out, kv)
		}
	}
	return out
}
