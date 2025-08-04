package commands

import (
	"fmt"

	"github.com/Inmovilizame/invoiceling/internal/container"
	"github.com/Inmovilizame/invoiceling/pkg/i18n"

	"github.com/spf13/cobra"
)

// pdfCmd represents the pdf command
var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		invoiceID, err := cmd.Flags().GetString("invoice")
		cobra.CheckErr(err)

		draft, err := cmd.Flags().GetBool("draft")
		cobra.CheckErr(err)

		renderer, err := cmd.Flags().GetString("renderer")
		cobra.CheckErr(err)

		languageStr, err := cmd.Flags().GetString("language")
		cobra.CheckErr(err)

		language, err := i18n.ParseLanguage(languageStr)
		cobra.CheckErr(err)

		doc, err := container.NewDocumentService(renderer, draft, language)
		cobra.CheckErr(err)

		is := container.NewInvoiceService()
		invoice := is.Read(invoiceID)

		err = doc.Render(invoice)
		cobra.CheckErr(err)

		fmt.Println("Generated PDF for:", invoiceID)
	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)

	pdfCmd.Flags().StringP("invoice", "i", "", "Invoice id to render")
	pdfCmd.Flags().BoolP("draft", "d", false, "Generate draft PFD")
	pdfCmd.Flags().StringP("renderer", "r", "Basic", "Generate draft PFD")
	pdfCmd.Flags().StringP("language", "l", "en", "Language for the PDF (en, es)")
}
