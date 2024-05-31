/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
package commands

import (
	"fmt"

	"github.com/Inmovilizame/invoiceling/internal/container"

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

		doc, err := container.NewDocumentService(draft)
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
}
