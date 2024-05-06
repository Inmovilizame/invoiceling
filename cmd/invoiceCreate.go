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

	"github.com/Inmovilizame/invoiceling/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// invoiceCreateCmd represents the invoiceCreate command
var invoiceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new invoice file",
	Long: `Create a new invoice file to be stored as JSON file.
	The name of the file matches invoice number.`,
	Run: createInvoiceCmdFunc,
}

func init() {
	invoiceCmd.AddCommand(invoiceCreateCmd)
	invoiceCreateCmd.Flags().IntP("id", "i", 0, "Invoice ID")
	invoiceCreateCmd.Flags().StringP("client", "c", "client1", "Invoice client")
	invoiceCreateCmd.Flags().IntP("due", "d", 30, "Invoice due date")
}

func createInvoiceCmdFunc(cmd *cobra.Command, _ []string) {
	clientID, err := cmd.Flags().GetString("client")
	cobra.CheckErr(err)

	due, err := cmd.Flags().GetInt("due")
	cobra.CheckErr(err)

	invoiceID, err := cmd.Flags().GetInt("id")
	cobra.CheckErr(err)

	fs := service.NewFreelancer()
	me := fs.Get()

	cs := service.NewClientFS(viper.GetString("dirs.client"))
	client := cs.Read(clientID)

	is := service.NewInvoiceFS(
		viper.GetString("invoice.currency"),
		viper.GetString("invoice.id_format"),
		viper.GetString("invoice.logo"),
		viper.GetString("dirs.invoice"),
	)
	invoice, err := is.Create(invoiceID, me, client, due, service.DF_YMD)
	cobra.CheckErr(err)

	fmt.Printf("Invoice created: %s\n", invoice.ID)
}
