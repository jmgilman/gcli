package cmd

import (
	"context"
	"fmt"
	gcert "github.com/jmgilman/gcert/proto"
	"github.com/jmgilman/gcli/rpc"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// renewCmd represents the renew command
var renewCmd = &cobra.Command{
	Use:   "renew [gcert server] [domain1] [domain 2] ...",
	Args: cobra.MinimumNArgs(2),
	Short: "Requests the gcert service to renew the given domain[s] certificates in Vault",
	Long: `Sends a request to the gcert service, asking it to renew the SSL certificates in Vault for the given domains.
It will return the paths to where the certificates were written to. You can use the fetch command to get the contents
of the certificates and write them to the local filesystem.`,
	Run: func(cmd *cobra.Command, args []string) {
		NewCertificateRequest(args[0], args[1:])
	},
}

func init() {
	certCmd.AddCommand(renewCmd)
}

func NewCertificateRequest(server string, domains []string) {
	conn, err := rpc.Dial(server, true)
	if err != nil {
		fmt.Println("Unable to connec to RPC server at", server)
		os.Exit(1)
	}

	client := gcert.NewCertificateServiceClient(conn)
	request := &gcert.CertificateRequest{
		Domains:  domains,
		Endpoint: gcert.CertificateRequest_LE_STAGING,
	}

	resp, err := client.GetCertificate(context.Background(), request)
	if err != nil || !resp.Success {
		fmt.Println("Error requesting certificate:", err)
		os.Exit(1)
	}

	fmt.Printf("New certificates saved at:\n\n%s\n", strings.Join(resp.VaultPaths, "\n"))
}
