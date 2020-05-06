package cmd

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/jmgilman/gcli/ui"
	"github.com/jmgilman/gcli/vault/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"time"
)

var expiration bool

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
			domainString, ok := domain.(string)
			if !ok {
				ui.ErrorThenExit("Error parsing domain as a string", nil)
			}
			if expiration {
				expDate, err := getExpirationDate(vaultClient, domainString)
				if err != nil {
					ui.ErrorThenExit("Error getting certificate expiration date", err)
				}

				expDatePST, err := timeToPST(expDate)
				if err != nil {
					ui.ErrorThenExit("Error parsing expiration date to PST", err)
				}

				fmt.Printf("%s     %v", domainString, expDatePST.Format("January 02, 2006"))
			} else {
				fmt.Println(domainString)
			}
		}
	},
}

func init() {
	certCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&expiration, "expiration", "e", false, "Show the expiration dates of certificates (may take longer)")
}

func getExpirationDate(vaultClient *client.VaultClient, domain string) (time.Time, error) {
	secret, err := vaultClient.Read(filepath.Join(domainPrefix, domain))
	if err != nil {
		return time.Now(), err
	}

	certString, ok := secret.Data["certificate"].(string)
	if !ok {
		return time.Now(), fmt.Errorf("error unmarshalling certificate from server response")
	}

	certBytes, err := base64.StdEncoding.DecodeString(certString)
	if err != nil {
		return time.Now(), err
	}

	pemContents, _ := pem.Decode(certBytes)
	certificate, err := x509.ParseCertificate(pemContents.Bytes)
	if err != nil {
		return time.Now(), err
	}

	return certificate.NotAfter, nil
}

func timeToPST(t time.Time) (time.Time, error) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.Now(), err
	}

	return t.In(location), nil
}