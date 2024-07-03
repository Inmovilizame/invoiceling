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
		cmd.Printf("ID: Client | Date | Due\n")

		is := container.NewInvoiceService()
		for _, i := range is.List(filterInvoice("")) {
			cmd.Printf("%s: %s | %s | %s\n",
				i.ID,
				i.To.Name,
				i.Date.Format("2006-01-02"),
				i.Date.Add(i.Due).Format("2006-01-02"),
			)
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
