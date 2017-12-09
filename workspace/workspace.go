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

package workspace

import (
	"log"
	"os"
	"path"

	"github.com/staffano/tcb/utils"
)

// Wd is the path to the workspace
var Wd string

// This is the stack used to push and pop paths.
var pathStack = make([]string, 23)
var pathStackIdx = 0

// InitWorkspace assumes the Wd variable has been set as it will
// try to make sure the workspace exists.
func InitWorkspace() {

	var (
		err error
	)
	if !PathExists("") {
		if err = os.Mkdir(Wd, 0777); err != nil {
			log.Fatalf("Could not create the workspace directory: %s", Wd)
		}
		log.Printf("Created workspace directory %s", Wd)
	}
	ResetPathStack()
	log.Printf("Workspace initialized at %s", Wd)
}

// Path returns a path within the workspace
// example:
// s := workspace.Path("meta-tooldirs", "bin")
func Path(elem ...string) string {
	t := make([]string, len(elem)+1)
	for i := range elem {
		t[i+1] = elem[i]
	}
	t[0] = Wd
	return path.Join(t...)

}

// PathExists returns true if the path within the workspace exists
func PathExists(elem ...string) bool {
	truePath := Path(elem...)
	return utils.PathExists(truePath)
}

// CWD sets the current working directory as used by the operations
// in the os package.
func cwd(elem ...string) {
	truePath := Path(elem...)
	if PathExists(elem...) != true {
		log.Fatalf("Cant set CWD to a non-existing workspace path. [%s]", truePath)
	}
	if err := os.Chdir(truePath); err != nil {
		log.Fatalf("Error calling os.Chdir on path %s", truePath)
	}
	log.Printf("Chdir to %s", truePath)
}

// ResetPathStack invalidates any active push/pops and sets the current working
// directory to Wd
func ResetPathStack() {
	pathStackIdx = 0
	cwd("")
}

// Push changes the current directory and sets the argument as the
// new current directory within the workspace
func Push(elem ...string) {
	var (
		cd  string
		err error
	)
	if cd, err = os.Getwd(); err != nil {
		log.Fatalf("Could not retrieve current directory")
	}
	pathStack[pathStackIdx] = cd
	pathStackIdx++
	cwd(elem...)
}

// Pop returns to the previous current directory in the stack
func Pop() {
	if pathStackIdx == 0 {
		log.Fatalf("Trying to pop from an empty Path Stack")
	}
	pathStackIdx--
	nd := pathStack[pathStackIdx]
	if err := os.Chdir(nd); err != nil {
		log.Fatalf("Could not pop the current working dir to %s", nd)
	}
	log.Printf("Chdir to %s", nd)
}

// MakeDir creates a directory structure inside the workspace
func MakeDir(permission os.FileMode, elem ...string) {
	truePath := Path(elem...)
	if err := os.MkdirAll(truePath, permission); err != nil {
		log.Fatalf("Could not create path %v", truePath)
	}
}

// SetStamp creates/updates the stamp file for entity
func SetStamp(entity string) {
	var (
		err error
		fd  *os.File
	)
	MakeDir(0777, "stamps")
	if fd, err = os.Create(Path("stamps", entity)); err != nil {
		log.Fatalf("Error trying to create stamp file %s", (Path("stamps", entity)))
	}
	defer fd.Close()
}

// GetStamp returns true if the stamp is set
func GetStamp(entity string) bool {
	return PathExists("stamps", entity)
}

// RemoveStamp deletes the stamp file from the stamp directory.
func RemoveStamp(entity string) {
	if err := os.Remove(Path("stamps", entity)); err != nil {
		log.Fatalf("Error trying to remove stamp file %s", (Path("stamps", entity)))
	}
}

// Reset the workspace. Removing all files and also the directory
func Reset() {
	if err := os.RemoveAll(Wd); err != nil {
		log.Fatalf("Error resetting workdir %s, %v", Wd, err)
	}
}
