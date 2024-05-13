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

		vatId, err := cmd.Flags().GetString("vatId")
		cobra.CheckErr(err)

		address1, err := cmd.Flags().GetString("address1")
		cobra.CheckErr(err)

		address2, err := cmd.Flags().GetString("address2")
		cobra.CheckErr(err)

		cs := container.NewClientService()

		err = cs.Create(id, name, vatId, address1, address2)
		cobra.CheckErr(err)
	},
}

func init() {
	clientCmd.AddCommand(clientCreateCmd)

	clientCreateCmd.Flags().SortFlags = false
	clientCreateCmd.Flags().StringP("name", "n", "", "Client mame [req]")
	clientCreateCmd.Flags().StringP("vat_id", "v", "", "Client VAT ID [req]")
	clientCreateCmd.Flags().StringP("address1", "s", "", "Client address street info [req]")
	clientCreateCmd.Flags().StringP("address2", "c", "", "Client address region state country")
	clientCreateCmd.Flags().StringP("id", "i", "", "Provide a custom client id")

	err := clientCreateCmd.MarkFlagRequired("name")
	cobra.CheckErr(err)

	err = clientCreateCmd.MarkFlagRequired("vat_id")
	cobra.CheckErr(err)

	err = clientCreateCmd.MarkFlagRequired("address1")
	cobra.CheckErr(err)
}
