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

package docker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var proxyVars = [...]string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy",
	"FTP_PROXY", "ftp_proxy", "NO_PROXY", "no_proxy"}

var dockerFile = `FROM ubuntu
RUN apt update -y && apt upgrade -y
RUN apt install -y build-essential gnat-5 git locales python3 wget m4 gawk unzip nano texinfo
RUN locale-gen en_US.UTF-8
RUN git clone https://github.com/openembedded/bitbake.git && cd /bitbake && git checkout 1.32
ENV LANG en_US.UTF-8
ENV PATH /bitbake/bin:$PATH
ENV PYTHONPATH /bitbake/lib:$PYTHONPATH
RUN mkdir -p /build/conf
VOLUME /meta-crosstools
VOLUME /build/tmp
RUN printf "BBPATH = \"${TOPDIR}\"\nBBFILES ?= \"\"\nBBLAYERS ?= \"/meta-crosstools\"\n" >> /build/conf/bblayers.conf
WORKDIR /build
CMD ["bitbake", "--help"]
#ENTRYPOINT [ "/build.sh"]
`

func getProxyArgs(argName string) []string {
	var res []string

	// Append environment proxy variables if they are needed
	for _, s := range proxyVars {
		if val := os.Getenv(s); val != "" {
			res = append(res, argName, fmt.Sprintf("%s=%s", s, val))
		}
	}
	return res
}

func getVolumeArgs(buildTmpVol, resultDir, metaCrosstoolsDir, localConfPath string) []string {
	var res []string

	res = append(res, "-v", fmt.Sprintf("%s:/build/tmp", buildTmpVol))
	res = append(res, "-v", fmt.Sprintf("%s:/build/RESULT", resultDir))
	res = append(res, "-v", fmt.Sprintf("%s:/meta-crosstools/", metaCrosstoolsDir))
	res = append(res, "--mount", fmt.Sprintf("type=bind,source=%s,target=/build/conf/local.conf,readonly", localConfPath))

	return res
}

// BuildImage creates the image we will base our container on
func BuildImage() {

	args := []string{"build"}
	args = append(args, getProxyArgs("--build-arg")...)
	args = append(args, "-t", "meta_crosstools_bitbake", "-")

	cmd := exec.Command("docker", args...)

	// Connect dockerFile to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, dockerFile)
	}()

	handleCmdOutput(cmd, "docker build")
}

// Prune all tcb related docker stuff
func Prune(buildTmpVol string) {
	cmd := exec.Command("docker", "volume", "rm", "-f", buildTmpVol)
	handleCmdOutput(cmd, "docker volume rm")

	cmd = exec.Command("docker", "image", "rm", "-f", "meta_crosstools_bitbake")
	handleCmdOutput(cmd, "docker image rm meta_crosstools_bitbake")
}

// Execute set up the container and executes it with a bash command
func Execute(buildTmpVol, resultDir, metaCrosstoolsDir, localConfPath string, arguments ...string) {

	args := []string{"run", "-i", "--rm"}
	args = append(args, getProxyArgs("--env")...)
	args = append(args, getVolumeArgs(buildTmpVol, resultDir, metaCrosstoolsDir, localConfPath)...)
	args = append(args, "meta_crosstools_bitbake")
	args = append(args, arguments...)

	if viper.GetBool("dryrun") {
		fmt.Printf("docker %s", strings.Trim(fmt.Sprintf("%v", args), "[]"))
	} else {
		cmd := exec.Command("docker", args...)
		handleCmdOutput(cmd, "docker run")
	}
}

func handleCmdOutput(cmd *exec.Cmd, prefix string) error {

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
		os.Exit(1)
	}

	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("[%s] ERR %s| %s\n", time.Now().Format(time.StampMilli), prefix, errScanner.Text())
		}
	}()

	cmdStdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	stdOutScanner := bufio.NewScanner(cmdStdOutReader)
	go func() {
		for stdOutScanner.Scan() {
			fmt.Printf("[%s] OUT %s| %s\n", time.Now().Format(time.StampMilli), prefix, stdOutScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
