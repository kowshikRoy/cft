/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/kowshikRoy/cft/model"
	"github.com/kowshikRoy/cft/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "initializes your contest",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"int"},
	Example: `
    # initialize your dev environment for contest-id 1000
    cft init 1000
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		contestID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("Not a valid contest id: %v", err)
		}
		contestDir, testDir, _ := util.CreateDev(contestID)
		standings := model.Standings{}
		if err := util.Crawl(contestID, &standings); err != nil {
			return err
		}
		if err := setupFiles(contestID, &standings.Result); err != nil {
			return err
		}
		if err := fetchTestCase(contestID, testDir, &standings.Result.Problems); err != nil {
			return err
		}
		fmt.Println("Your contest environment is ready on :", contestDir, "\nHappy Coding !!")
		return nil
	},
}

func fetchTestCase(contestID int, testDir string, problems *[]model.Problems) error {
	c := make(chan bool)
	for _, p := range *problems {
		go util.CrawlTest(c, testDir, p.Index, "https://codeforces.com/contest/"+strconv.Itoa(contestID)+"/problem/"+p.Index)
	}
	for range *problems {
		<-c
	}
	return nil
}

func setupFiles(contestID int, r *model.Result) error {
	contestDir := util.ContestDir(contestID)
	lang := viper.GetString("language")
	tpl := []byte{}
	if viper.IsSet("templates") {
		temp, err := ioutil.ReadFile(viper.GetString("templates." + lang))
		if err != nil {
			return fmt.Errorf("Couldn't read template file: %v", err)
		}
		tpl = temp
	}
	for _, problem := range r.Problems {
		ext := util.GetSourceFileExtension(viper.GetString("language"))
		if err := ioutil.WriteFile(path.Join(contestDir, problem.Index+ext), tpl, 0755); err != nil {
			return fmt.Errorf("Couldn't write the source file: %v", err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
