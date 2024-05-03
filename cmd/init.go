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
	"os"
	"slices"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var (
	allowedFormats = []string{"yaml", "yml", "json"}
	dirs           = []string{"client", "invoice", "pdf", "static"}

	defaultMask = os.FileMode(0755)
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialiaze a folder with default configurations",
	Long: `Initialiaze the current folder with default configurations. This can be
	changed manually by editing the generated config file.`,
	Run: func(cmd *cobra.Command, _ []string) {
		if viper.ConfigFileUsed() != "" {
			fmt.Println("Directory already initialized, exiting...")
			return
		}

		format := cmd.Flag("format").Value.String()
		if !slices.Contains(allowedFormats, format) {
			cobra.CheckErr(fmt.Errorf("format option '%s' not allowed", format))
		}

		fmt.Println("Generating folder structure...")
		for _, dir := range dirs {
			err := os.Mkdir(dir, defaultMask)
			cobra.CheckErr(err)
		}

		fmt.Println("Generating default configuration file...")
		defaultConfig()
		err := viper.WriteConfigAs(fmt.Sprintf("./config.%s", format))
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("format", "f", "yaml", "Configuration file format: yaml, json, toml")
}

func defaultConfig() {
	for _, dir := range dirs {
		viper.SetDefault(
			fmt.Sprintf("dirs.%s", dir),
			fmt.Sprintf("./%s", dir),
		)
	}

	viper.SetDefault("invoice.currency", "EUR")
	viper.SetDefault("invoice.logo", "./static/logo.png")
	viper.SetDefault("invoice.id_format", "F%s-%03d")

	viper.SetDefault("freelancer.company", "Your Company Name")
	viper.SetDefault("freelancer.name", "Your Full Name")
	viper.SetDefault("freelancer.email", "your.email@example.com")
	viper.SetDefault("freelancer.phone", "+99 123456789")
	viper.SetDefault("freelancer.vat_id", "CC12345678A")
	viper.SetDefault("freelancer.address1", "Your Street Address")
	viper.SetDefault("freelancer.address2", "City, ST, Zip Code")

	viper.SetDefault("payment.holder", "Bank account holder")
	viper.SetDefault("payment.iban", "CC00 1234 1234 12 1234567890")
	viper.SetDefault("payment.swift", "ABCDDEFFXXX")
}