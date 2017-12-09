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
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/staffano/tcb/workspace"
)

const defaultWorkspaceName = "tcb_workspace"

var home string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tcb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	var (
		defaultWorkspace string
		err              error
	)
	home, err = homedir.Dir()
	if err != nil {
		log.Fatalf("Could not determine homedir")
	}

	cobra.OnInitialize(initConfig)
	defaultWorkspace = path.Join(home, defaultWorkspaceName)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&(workspace.Wd), "ws", defaultWorkspace, "workspace (default is "+defaultWorkspace+")")
	RootCmd.PersistentFlags().StringP("builder.repo.url", "", "https://github.com/staffano/meta-crosstools.git", "Repository to use for builder.")
	viper.BindPFlag("builder.repo.url", RootCmd.PersistentFlags().Lookup("builder.repo.url"))
	RootCmd.PersistentFlags().BoolP("keep-sources", "", false, "If set, git pull will not be called for the source directory")
	viper.BindPFlag("keep-sources", RootCmd.PersistentFlags().Lookup("keep-sources"))
	RootCmd.PersistentFlags().BoolP("dryrun", "", false, "If set, build commands will not be executed, but printed to stdout instead.")
	viper.BindPFlag("dryrun", RootCmd.PersistentFlags().Lookup("dryrun"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if workspace.Wd != "" {
		// Use config file from the flag.
		viper.SetConfigName("config")
		viper.AddConfigPath(workspace.Wd)
	}

	workspace.InitWorkspace()

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Read config file %s", viper.ConfigFileUsed())
	} else {
		log.Printf("Not using config file")
	}
}
