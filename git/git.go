// Copyright © 2017 Staffan Olsson <staffano@diversum.nu>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package git

import (
	"log"
	"os/exec"
)

// Clone a repository to the current directory
// example err := git.Clone("https://github.com/staffano/meta-crosstools", "meta-crosstools", "master")
func Clone(srcPath string, dstPath string, rev string) error {
	var (
		err error
	)
	log.Printf("git clone -b %s --single-branch %s %s", rev, srcPath, dstPath)
	if _, err = exec.Command("git", "clone", "-b", rev, "--single-branch",
		srcPath, dstPath).Output(); err != nil {
		log.Println("Error cloning: ", err)
		return err
	}
	return nil
}

// Pull the repository at the current workdir.
// example err := git.Pull()
func Pull() error {
	var (
		err error
	)
	log.Printf("git pull")
	if _, err = exec.Command("git", "pull").Output(); err != nil {
		log.Printf("Error pulling: %v", err)
		return err
	}
	return nil
}