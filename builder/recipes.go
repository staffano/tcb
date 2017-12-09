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

package builder

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"github.com/staffano/tcb/git"
	"github.com/staffano/tcb/utils"
	"github.com/staffano/tcb/workspace"
)

func init() {
	viper.SetDefault("builder.repo.url", "ssh://staffan@homeland:/srv/repos/git/meta-crosstools.git")
	viper.SetDefault("builder.repo.rev", "master")
}

const metaCrosstools = "meta-crosstools"

// CheckoutMetaCrosstools installs the meta-crosstools repository from
// github.com/staffano/meta-crosstools
func CheckoutMetaCrosstools() {

	// Check if the directory exists
	metaCrosstoolsDir := workspace.Path(metaCrosstools)
	ctExists := workspace.PathExists(metaCrosstools)
	log.Printf("%s exists: %t", metaCrosstoolsDir, ctExists)
	if ctExists {
		workspace.Push(metaCrosstools)
		git.Pull()
		workspace.Pop()
	} else {
		git.Clone(viper.GetString("builder.repo.url"), metaCrosstoolsDir, viper.GetString("builder.repo.rev"))
	}
}

// GetAllTargets retrieves all targets defined inside meta-crosstools/conf/toolchain
func GetAllTargets() []string {
	var result []string
	confDir := workspace.Path("meta-crosstools", "conf", "toolchains")
	files, err := ioutil.ReadDir(confDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".conf") {
			result = append(result, strings.TrimSuffix(file.Name(), ".conf"))
		}
	}
	return result
}

// GetDependencies returns a list of other targets this target depends on
func GetDependencies(target string) []string {
	src := workspace.Path("meta-crosstools", "conf", "toolchains", target+".conf")
	content, _ := ioutil.ReadFile(src)
	re := regexp.MustCompile(`^#\s* ///\s*DEPENDENCIES=(\S*)+`)
	matches := re.FindSubmatch(content)
	if matches == nil {
		return nil
	}
	fmt.Println(matches[1])
	return strings.Split(string(matches[1][:]), ",")
}

// SetTarget initializes the local.conf for the target
func SetTarget(target string) {
	// Now, lets build this specific target
	workspace.MakeDir(0777, "build", "conf")

	// Create a copy of $wsp/meta-crosstools/conf/toolchains/<target>.conf to
	// $wsp/build/conf/local.conf
	src := workspace.Path("meta-crosstools", "conf", "toolchains", target+".conf")
	dst := workspace.Path("build", "conf", "local.conf")
	os.RemoveAll(dst)
	utils.CopyFile(src, dst)
	log.Printf("Copied %s to %s", src, dst)

	// Append some specifics to local.conf
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(`
BB_NUMBER_THREADS = "12"
		
MAKE_JX := "-j12"
	`); err != nil {
		panic(err)
	}

}
