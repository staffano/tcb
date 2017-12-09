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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/staffano/tcb/builder"
	"github.com/staffano/tcb/docker"
	"github.com/staffano/tcb/workspace"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build toolchain(s)",
	Run:   Build,
}

func init() {
	RootCmd.AddCommand(buildCmd)
}

// Build builds the targets specified. If no target och "all" target is
// specified then all known targets will be built.
// If preconfigured is set, then the targets specified in the builder
// will be built
func Build(cmd *cobra.Command, targets []string) {
	log.Printf("builder.Build(%v)", targets)
	if !viper.GetBool("keep-sources") {
		builder.CheckoutMetaCrosstools()
	}

	if len(targets) == 0 || targets[0] == "all" {
		allTargets := builder.GetAllTargets()
		for _, t := range allTargets {
			BuildTarget(t)
		}
	} else {
		BuildTarget(targets[0])
	}
}

// BuildTarget builds one target
func BuildTarget(target string) {
	stamp := target + ".build"
	// Skip if already built
	if workspace.GetStamp(stamp) {
		return
	}

	// Then we have to build any dependencies
	deps := builder.GetDependencies(target)
	for _, t := range deps {
		BuildTarget(t)
	}

	builder.SetTarget(target)

	// Make sure docker image is built
	docker.BuildImage()

	// Run bitbake dockerized
	workspace.MakeDir(0777, "results")
	docker.Execute("bb-build-vol", workspace.Path("results"), workspace.Path("meta-crosstools"),
		workspace.Path("build", "conf", "local.conf"), "bitbake", "image")
	workspace.SetStamp(stamp)
}
