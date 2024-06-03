package commands

import (
	"github.com/Inmovilizame/invoiceling/internal/container"
	"github.com/spf13/cobra"
)

// clientCreateCmd represents the clientCreate command
var clientCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new client entry",
	Long: `Creates a new client entry. If id is not provided, the vat id
will be used to compose a unique id`,
	Run: func(cmd *cobra.Command, _ []string) {
		id, err := cmd.Flags().GetString("id")
		cobra.CheckErr(err)

		name, err := cmd.Flags().GetString("name")
		cobra.CheckErr(err)

		vatID, err := cmd.Flags().GetString("vat_id")
		cobra.CheckErr(err)

		address1, err := cmd.Flags().GetString("address1")
		cobra.CheckErr(err)

		address2, err := cmd.Flags().GetString("address2")
		cobra.CheckErr(err)

		cs := container.NewClientService()

		err = cs.Create(id, name, vatID, address1, address2)
		cobra.CheckErr(err)
	},
}

func init() {
	clientCmd.AddCommand(clientCreateCmd)

	clientCreateCmd.Flags().SortFlags = false
	clientCreateCmd.Flags().StringP("id", "i", "", "Provide a custom client id")
	clientCreateCmd.Flags().StringP("name", "n", "", "Client mame [req]")
	clientCreateCmd.Flags().StringP("vat_id", "v", "", "Client VAT ID [req]")
	clientCreateCmd.Flags().StringP("address1", "s", "", "Client address street info [req]")
	clientCreateCmd.Flags().StringP("address2", "c", "", "Client address region state country")

	err := clientCreateCmd.MarkFlagRequired("name")
	cobra.CheckErr(err)

	err = clientCreateCmd.MarkFlagRequired("vat_id")
	cobra.CheckErr(err)

	err = clientCreateCmd.MarkFlagRequired("address1")
	cobra.CheckErr(err)
}
