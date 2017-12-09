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

package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/staffano/tcb/docker"
	"github.com/staffano/tcb/workspace"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Really cleans everything...",
	Run:   Clean,
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}

// Clean results and intermediate files
func Clean(cmd *cobra.Command, targets []string) {
	log.Printf("builder.Clean(%v)", targets)
	if len(targets) == 0 {
		docker.Prune("bb-build-vol")
		workspace.Reset()
		return
	}
	for _, t := range targets {
		switch ut := strings.ToUpper(t); ut {
		case "STAMPS":
			os.RemoveAll(workspace.Path("stamps"))
		case "RESULTS":
			os.RemoveAll(workspace.Path("results"))
		case "DOCKER":
			docker.Prune("bb-build-vol")
		case "ALL":
			docker.Prune("bb-build-vol")
			workspace.Reset()
		}
	}
}
