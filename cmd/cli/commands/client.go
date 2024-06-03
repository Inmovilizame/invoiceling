package commands

import (
	"strings"

	"github.com/Inmovilizame/invoiceling/internal/container"
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		filter, err := cmd.Flags().GetString("filter")
		cobra.CheckErr(err)

		cs := container.NewClientService()

		for _, c := range cs.List(filterClient(filter)) {
			cmd.Printf("Client: %s | %s \n", c.Name, c.VatID)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringP("filter", "f", "", "Filter clients by name")
}

func filterClient(filter string) repository.Filter[*model.Client] {
	return func(c *model.Client) bool {
		if filter == "" {
			return true
		}

		if strings.Contains(c.ID, filter) ||
			strings.Contains(c.Name, filter) ||
			strings.Contains(c.VatID, filter) ||
			strings.Contains(c.Address1, filter) ||
			strings.Contains(c.Address2, filter) {
			return true
		}

		return false
	}
}
