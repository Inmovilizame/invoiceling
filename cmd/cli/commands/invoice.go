package commands

import (
	"strings"

	"github.com/Inmovilizame/invoiceling/internal/container"
	"github.com/Inmovilizame/invoiceling/internal/repository"
	"github.com/Inmovilizame/invoiceling/pkg/model"
	"github.com/spf13/cobra"
)

// invoiceCmd represents the invoice command
var invoiceCmd = &cobra.Command{
	Use:   "invoice",
	Short: "invoice commands",
	Long:  `invoice related commands`,
	Run: func(cmd *cobra.Command, _ []string) {
		is := container.NewInvoiceService()
		for _, i := range is.List(filterInvoice("")) {
			cmd.Printf("Invoice: %s | %s \n", i.ID, i.Date.Format("2006-01-02"))
		}
	},
}

func init() {
	rootCmd.AddCommand(invoiceCmd)
}

func filterInvoice(filter string) repository.Filter[*model.Invoice] {
	return func(i *model.Invoice) bool {
		if filter == "" {
			return true
		}

		if strings.Contains(i.ID, filter) ||
			strings.Contains(i.To.Name, filter) ||
			strings.Contains(i.To.VatID, filter) ||
			strings.Contains(i.To.Address1, filter) ||
			strings.Contains(i.To.Address2, filter) ||
			strings.Contains(strings.Join(i.Notes.ToSlice(), ":"), filter) {
			return true
		}

		return false
	}
}
