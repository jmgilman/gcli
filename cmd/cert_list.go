package cmd

import (
	"fmt"
	"github.com/jmgilman/gcli/ui"
	"github.com/jmgilman/gcli/vault/client"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the domains which have certificates stored in Vault",
	Long: ``,
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

		domains, err := vaultClient.List(domainPrefix)
		if err != nil {
			ui.ErrorThenExit("Error fetching list of domains", err)
		}

		fmt.Println("Domains:")
		for _, domain := range domains {
			fmt.Println(domain)
		}
	},
}

func init() {
	certCmd.AddCommand(listCmd)
}