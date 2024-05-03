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
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// invoiceCreateCmd represents the invoiceCreate command
var invoiceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new invoice file",
	Long: `Create a new invoice file to be stored as JSON file.
	The name of the file matches invoice number.`,
	Run: createCmdFunc,
}

func init() {
	invoiceCmd.AddCommand(invoiceCreateCmd)
	invoiceCreateCmd.Flags().StringP("client", "c", "client1", "Client to create invoice for")
}

func createCmdFunc(cmd *cobra.Command, args []string) {
	today := time.Now()

	fmt.Printf("F%s-%03d", today.Format("06"), 1)
	fmt.Println("")

	//me := model.LoadFreelancer()
	//
	//clientID, err := cmd.Flags().GetString("client")
	//cobra.CheckErr(err)
	//
	//client, err := model.LoadClient(getClientSrcFile(clientID))
	//cobra.CheckErr(err)
	//
	//invoice := model.NewInvoice(&me, &client, 30, model.DF_YMD)
	//err = invoice.Save(viper.GetString("dirs.invoice"))
	//cobra.CheckErr(err)
	//
	//fmt.Printf("Invoice created: %s\n", invoice.ID)
}

func getClientSrcFile(clientID string) string {
	return filepath.Join(viper.GetString("dirs.client"), clientID+".json")
}
