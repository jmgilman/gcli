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
	"encoding/base64"
	"fmt"
	"github.com/jmgilman/gcli/ui"
	"github.com/jmgilman/gcli/vault/client"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
)

const domainPrefix = "secret/ssl"

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch [Domain] [/dir/to/write/to]",
	Short: "Fetches the certificates stored for the given domain and writes them to the given directory",
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

		data, err := fetchCerts(vaultClient, args[0])
		if err != nil {
			ui.ErrorThenExit("Error fetching certificates for domain", err)
		}

		if err := writeCerts(data, args[1]); err != nil {
			ui.ErrorThenExit("Error writing certificates", err)
		}

		fmt.Println("Wrote certificates to", args[1])
	},
}

func init() {
	certCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// fetchCerts fetches and returns the certificates stored for the given domain
func fetchCerts(vaulClient *client.VaultClient, domain string) (map[string]interface{}, error) {
	path := filepath.Join(domainPrefix, domain)
	secret, err := vaulClient.Read(path)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// Validate that we got the data back we expected
	if secret == nil || secret.Data == nil {
		return map[string]interface{}{}, fmt.Errorf("no certificates stored for %s", domain)
	}

	expectedKeys := []string{"cert_stable_url", "cert_url", "certificate", "issuer_certificate", "private_key"}
	for _, key := range expectedKeys {
		if _, ok := secret.Data[key]; !ok {
			return map[string]interface{}{}, fmt.Errorf("the server returned a malformed response: %v", secret.Data)
		}
	}

	return secret.Data, nil
}

func writeCerts(certs map[string]interface{},  path string) error {
	cert, ok := certs["certificate"].(string)
	if !ok {
		return fmt.Errorf("error unmarshalling certificate from response")
	}
	certBytes, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return err
	}

	caCert, ok := certs["issuer_certificate"].(string)
	if !ok {
		return fmt.Errorf("error unmarshalling CA certificate from response")
	}
	caCertBytes, err := base64.StdEncoding.DecodeString(caCert)
	if err != nil {
		return err
	}

	privKey, ok := certs["private_key"].(string)
	if !ok {
		return fmt.Errorf("error unmarshalling private key from response")
	}
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, "certificate.pem"), certBytes, 0644); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, "ca_certificate.pem"), caCertBytes, 0644); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(path, "key.pem"), privKeyBytes, 0644); err != nil {
		return err
	}

	return nil
}