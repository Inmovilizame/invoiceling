package commands

import (
	"fmt"
	"github.com/Inmovilizame/invoiceling/pkg/model"

	"github.com/Inmovilizame/invoiceling/internal/container"

	"github.com/spf13/cobra"
)

// invoiceAddItemCmd represents the invoiceAddItemCmd command
var invoiceAddItemCmd = &cobra.Command{
	Use:   "item",
	Short: "Add billable item to invoice",
	Long:  `Add a single billable item to a created invoice.`,
	Run: func(cmd *cobra.Command, _ []string) {
		invoiceID, err := cmd.Flags().GetString("invoice")
		cobra.CheckErr(err)

		desc, err := cmd.Flags().GetString("desc")
		cobra.CheckErr(err)

		rate, err := cmd.Flags().GetFloat64("rate")
		cobra.CheckErr(err)

		qty, err := cmd.Flags().GetInt("quantity")
		cobra.CheckErr(err)

		vat, err := cmd.Flags().GetFloat64("vat")
		cobra.CheckErr(err)

		is := container.NewInvoiceService()
		invoice := is.Read(invoiceID)

		item := model.Item{
			Description: desc,
			Quantity:    qty,
			Rate:        rate,
			Vat:         vat,
		}

		invoice = is.AddItems(invoice, []model.Item{item})

		fmt.Printf("Invoice %s updated\n", invoice.ID)
	},
}

func init() {
	invoiceCmd.AddCommand(invoiceAddItemCmd)

	invoiceAddItemCmd.Flags().StringP("invoice", "i", "", "Invoice ID")
	invoiceAddItemCmd.Flags().StringP("desc", "d", "", "Item description")
	invoiceAddItemCmd.Flags().Float64P("rate", "r", 0.0, "Item price")
	invoiceAddItemCmd.Flags().IntP("quantity", "q", 1, "Item quantity")
	invoiceAddItemCmd.Flags().Float64P("vat", "v", 0, "Item VAT")

	err := invoiceAddItemCmd.MarkFlagRequired("invoice")
	cobra.CheckErr(err)

	err = invoiceAddItemCmd.MarkFlagRequired("desc")
	cobra.CheckErr(err)

	err = invoiceAddItemCmd.MarkFlagRequired("rate")
	cobra.CheckErr(err)
}
