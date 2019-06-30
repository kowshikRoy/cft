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

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Args:  cobra.ExactArgs(2),
	Short: "Sets an individual value in a cftconfig file",
	Long: `Sets an individual value in a cftconfig file

	PROPERTY_NAME is a name of the property.

	PROPERTY_VALUE is the new value you wish to set.
	`,
	Example: `# Set the user of the account
	cft config set user tourist

	# Set the directory for the codeforces contest
	cft config set workdir /Users/tourist/code/codeforces
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set(args[0], args[1])
		err := viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		fmt.Println("Config Saved in ", viper.ConfigFileUsed())
		return nil
	},
}

func init() {
	configCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
