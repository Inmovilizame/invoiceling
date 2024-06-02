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

	"github.com/spf13/viper"

	"github.com/Inmovilizame/invoiceling/pkg/model"

	"github.com/Inmovilizame/invoiceling/internal/container"

	"github.com/spf13/cobra"
)

// invoiceCreateCmd represents the invoiceCreate command
var invoiceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new invoice file",
	Long: `Create a new invoice file to be stored as JSON file.
	The name of the file matches invoice number.`,
	Run: func(cmd *cobra.Command, _ []string) {
		clientID, err := cmd.Flags().GetString("client")
		cobra.CheckErr(err)

		due, err := cmd.Flags().GetInt("due")
		cobra.CheckErr(err)

		invoiceID, err := cmd.Flags().GetInt("id")
		cobra.CheckErr(err)

		vat, err := cmd.Flags().GetFloat64("vat")
		cobra.CheckErr(err)

		retention, err := cmd.Flags().GetFloat64("retention")
		cobra.CheckErr(err)

		note, err := cmd.Flags().GetString("note")
		cobra.CheckErr(err)

		is := container.NewInvoiceService()
		invoice, err := is.Create(invoiceID, clientID, due, note, vat, retention)
		cobra.CheckErr(err)

		fmt.Printf("InvoiceService created: %s\n", invoice.ID)
	},
}

func init() {
	invoiceCmd.AddCommand(invoiceCreateCmd)

	defaultDue := model.DefaultDueSpan
	defaultNote := "Thank you for your business. Please add the invoice number to your payment description."
	defaultRet := viper.GetFloat64("retention")
	defaultVat := viper.GetFloat64("vat")

	invoiceCreateCmd.Flags().IntP("id", "i", 0, "InvoiceService ID")
	invoiceCreateCmd.Flags().StringP("client", "c", "", "InvoiceService client")
	invoiceCreateCmd.Flags().IntP("due", "d", defaultDue, "InvoiceService due date")
	invoiceCreateCmd.Flags().Float64P("vat", "v", defaultVat, "InvoiceService VAT")
	invoiceCreateCmd.Flags().Float64P("retention", "r", defaultRet, "InvoiceService Retention (Spanish IRPF)")
	invoiceCreateCmd.Flags().StringP("note", "n", defaultNote, "Add invoice note")

	err := invoiceCreateCmd.MarkFlagRequired("client")
	cobra.CheckErr(err)
}
