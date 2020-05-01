/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request [gcert server] [domain1] [domain 2] ...",
	Args: cobra.MinimumNArgs(2),
	Short: "Requests the gcert service to renew the given domain's certificate in Vault",
	Long: `Sends a request to the gcert service, asking it to renew the SSL certificates in Vault for the given domains.
It will return the paths to where the certificates were written to. You can use the fetch command to get the contents
of a certificate or the write command to write all certificates to the local filesystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("request called")
	},
}

func init() {
	certCmd.AddCommand(requestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// requestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// requestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
