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

	"github.com/kowshikRoy/cft/util"
	"github.com/kowshikRoy/cft/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes your contest",
	RunE: func(cmd *cobra.Command, args []string) error {
		contestID, err := cmd.Flags().GetInt("id")
		if err != nil {
			fmt.Errorf("%v", err)
		}
		util.CreateDir(ContestDir(contestID))
		util.CreateDir(TestDir(contestID))

		standings := model.Standings{}
		if err := util.Crawl(contestID, &standings); err != nil {
			return err
		}
		if err := setupFiles(contestID, &standings.Result); err != nil {
			return err
		}
		if err := fetchTestCase(contestID, &standings.Result.Problems); err != nil {
			return err
		}

		return nil

	},
}

func ContestDir(contestID int) string {
	woringDir := viper.GetString("workdir")
	contestDir := path.Join(woringDir, strconv.Itoa(contestID))
	return contestDir
}

func TestDir(contestID int) string {
	contestDir := ContestDir(contestID)
	testDir := path.Join(contestDir, "tests")
	return testDir

}
func fetchTestCase(contestID int, problems *[]model.Problems) error {
	c := make(chan bool)
	testDir := TestDir(contestID)
	for _, p := range *problems {
		go util.CrawlTest(c, testDir, p.Index, "https://codeforces.com/contest/"+strconv.Itoa(contestID)+"/problem/"+p.Index)
	}
	for range *problems {
		<-c
	}
	return nil
}

func setupFiles(contestID int, r *model.Result) error {
	contestDir := ContestDir(contestID)
	tpl := []byte{}
	if viper.IsSet("templateFile") {
		temp, err := ioutil.ReadFile(viper.GetString("templateFile"))
		if err != nil {
			return fmt.Errorf("Couldn't read template file: %v", err)
		}
		tpl = temp
	}
	for _, problem := range r.Problems {
		ext := util.FileExtension[viper.GetString("lang")]
		if err := ioutil.WriteFile(path.Join(contestDir, problem.Index+ext), tpl, 0755); err != nil {
			fmt.Errorf("Couldn't write the source file: %v", err)
		}
	}

	return nil
}

func init() {
	initCmd.Flags().Int("id", 0, "contest id (https://codeforces.com/contest/<contest id>")
	initCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
