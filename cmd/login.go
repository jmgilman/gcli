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
	"github.com/jmgilman/gcli/ui"
	"github.com/jmgilman/gcli/vault/auth"
	"github.com/jmgilman/gcli/vault/client"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var server string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates against Vault, retrieving a new token and persisting it to ~/.vault-token",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		vaultClient, err := client.NewDefaultClient()
		if err != nil {
			ui.ErrorThenExit("Error trying to load Vault client configuration", err)
		}

		if err := vaultClient.SetConfigValues(server, ""); err != nil {
			ui.ErrorThenExit("Error setting Vault server: ", err)
		}
		fmt.Println("login called")

		// Verify the vault is in a usable state
		status, err := vaultClient.Available()
		if err != nil {
			ui.ErrorThenExit("Error trying to check vault status", err)
		}

		if !status {
			fmt.Println("The vault is either sealed or not initialized - cannot continue")
			os.Exit(1)
		}

		login(vaultClient)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringVarP(&server, "server", "s", "", "Vault server to authenticate against")
}

func login(vaultClient *client.VaultClient) {
	// Ask which authentication type they would like to use
	prompt := ui.NewSelectPrompt("Please choose an authentication method:", auth.GetAuthNames())
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error getting authentication method:", err)
		os.Exit(1)
	}

	// Collect authentication details for the selected method
	authType := auth.Types[result]()
	details, err := ui.GetAuthDetails(authType, ui.NewPrompt)

	// Login with the collected details
	if err := vaultClient.Login(authType, details); err != nil {
		fmt.Println("Error logging in:", err)
		os.Exit(1)
	}

	home, err := homedir.Dir()
	if err != nil {
		ui.ErrorThenExit("Error getting user home directory", err)
	}

	tokenPath := filepath.Join(home, ".vault-token")
	if err := ioutil.WriteFile(tokenPath, []byte(vaultClient.Token()), 0644); err != nil {
		ui.ErrorThenExit("Error persising token to ~/.vault-token", err)
	}

	fmt.Println("Authentication successful! Token saved to", tokenPath)
}
