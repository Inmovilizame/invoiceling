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
package cmd

import (
	"fmt"

	"github.com/Inmovilizame/invoiceling/assets"
	"github.com/Inmovilizame/invoiceling/internal/pdf"
	"github.com/Inmovilizame/invoiceling/pkg/model"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		freelancer := model.Freelancer{
			Company:  "Your Company Name",
			Name:     "Your Full Name",
			Email:    "your.email@example.com",
			Phone:    "+99 123456789",
			VatID:    "CC12345678A",
			Address1: "Your Street Address",
			Address2: "City, ST, Zip Code",
		}

		client := model.Client{
			Name:     "Company Inc.",
			VatID:    "CC12345678A",
			Address1: "Company Street Address",
			Address2: "City, ST, Zip Code, Country",
		}

		invoice := model.NewInvoice(&freelancer, &client, 30, model.DF_YMD)

		interFont, err := assets.FS.ReadFile("fonts/Inter.ttf")
		cobra.CheckErr(err)

		interBoldFont, err := assets.FS.ReadFile("fonts/Inter-Bold.ttf")
		cobra.CheckErr(err)

		fonts := map[string][]byte{
			"Inter":      interFont,
			"Inter-Bold": interBoldFont,
		}

		doc, err := pdf.NewGoPdf(fonts)
		cobra.CheckErr(err)

		doc.Logo(invoice.Logo)
		doc.InvoiceInfo(invoice.ID, invoice.Date, invoice.Due)

		//pdf.WriteTitle(doc.pdfObject, invoice.Title, invoice.ID, invoice.Date)
		//pdf.WriteBillTo(doc.pdfObject, invoice.To)
		//pdf.WriteHeaderRow(doc.pdfObject)
		//
		//subtotal := 0.0
		//for i := range invoice.Items {
		//	q := 1
		//	if len(invoice.Quantities) > i {
		//		q = invoice.Quantities[i]
		//	}
		//
		//	r := 0.0
		//	if len(invoice.Rates) > i {
		//		r = invoice.Rates[i]
		//	}
		//
		//	pdf.WriteRow(doc.pdfObject, invoice.Items[i], q, r, invoice.Currency)
		//	subtotal += float64(q) * r
		//}
		//
		//if invoice.Note != "" {
		//	pdf.WriteNotes(doc.pdfObject, invoice.Note)
		//}
		//
		//pdf.WriteTotals(doc.pdfObject, subtotal, subtotal*invoice.Tax, subtotal*invoice.Discount, invoice.Currency)
		//
		//if invoice.Due != "" {
		//	pdf.WriteDueDate(doc.pdfObject, invoice.Due)
		//}
		//pdf.WriteFooter(doc.pdfObject, invoice.ID)

		output := "output/" + invoice.ID + ".pdf"
		err = doc.SaveTo(output)
		cobra.CheckErr(err)

		fmt.Printf("Generated %s\n", output)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)

	pdfCmd.Flags().StringP("invoice", "i", "", "Invoice id to render")
}
