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
		vaultClient, err := client.NewDefaultClientWithValues(server, "", ioutil.ReadFile)
		if err != nil {
			ui.ErrorThenExit("Error creating Vault client", err)
		}

		// Verify the vault is in a usable state
		status, err := vaultClient.Available()
		if err != nil {
			ui.ErrorThenExit("Error trying to check vault status", err)
		} else if !status {
			ui.ErrorThenExit("The vault is either sealed or not initialized - cannot continue", nil)
		}

		login(vaultClient)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&server, "server", "s", "", "Vault server to authenticate against")
}

func login(vaultClient *client.VaultClient) {
	// Ask which authentication type they would like to use
	prompt := ui.NewSelectPrompt("Please choose an authentication method:", auth.GetAuthNames())
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error getting authentication method", err)
		os.Exit(1)
	}

	// Collect authentication details for the selected method
	authType := auth.Types[result]()
	details, err := ui.GetAuthDetails(authType, ui.NewPrompt)

	// Login with the collected details
	if err := vaultClient.Login(authType, details); err != nil {
		fmt.Println("Error logging in", err)
		os.Exit(1)
	}

	home, err := homedir.Dir()
	if err != nil {
		ui.ErrorThenExit("Error getting user home directory", err)
	}

	tokenPath := filepath.Join(home, ".vault-token")
	if err := ioutil.WriteFile(tokenPath, []byte(vaultClient.Token()), 0644); err != nil {
		ui.ErrorThenExit("Error saving token to ~/.vault-token", err)
	}

	fmt.Println("Authentication successful! Token saved to", tokenPath)
}
