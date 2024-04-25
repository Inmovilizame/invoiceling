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
	dirs           = []string{"client", "config", "invoice", "pdf"}
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialiaze a folder with default configurations",
	Long: `Initialiaze the current folder with default configurations. This can be
	changed manually by editing the generated config file.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			err := os.Mkdir(dir, 0755)
			cobra.CheckErr(err)
		}

		fmt.Println("Generating default configuration file...")
		defaultConfig()
		err := viper.WriteConfigAs(fmt.Sprintf("./config/config.%s", format))
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("format", "f", "yaml", "Configuration file format: yaml, json, toml")
}

func defaultConfig() {
	for _, dir := range dirs {
		if dir == "config" { // Config dir is not configurable
			continue
		}
		viper.SetDefault(
			fmt.Sprintf("dirs.%s", dir),
			fmt.Sprintf("./%s", dir),
		)
	}

	viper.SetDefault("currency", "EUR")
	viper.SetDefault("logo", ".config/logo.png")

	viper.SetDefault("freelance.company", "Your Company Name")
	viper.SetDefault("freelance.name", "Your Full Name")
	viper.SetDefault("freelance.email", "your.email@example.com")
	viper.SetDefault("freelance.phone", "+99 123456789")
	viper.SetDefault("freelance.vat", "CC12345678A")
	viper.SetDefault("freelance.address1", "Your Street Address")
	viper.SetDefault("freelance.address2", "City, ST, Zip Code")

	viper.SetDefault("payment.holder", "Bank account holder")
	viper.SetDefault("payment.iban", "CC00 1234 1234 12 1234567890")
	viper.SetDefault("payment.swift", "ABCDDEFFXXX")
}
