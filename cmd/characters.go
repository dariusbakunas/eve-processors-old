/*
Copyright © 2020 Darius Bakunas-Milanowski <bakunas@gmail.com>

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
	"github.com/dariusbakunas/eve-processors"
	"github.com/joho/godotenv"
	"log"
	"github.com/spf13/cobra"
)

// charactersCmd represents the manual command
var charactersCmd = &cobra.Command{
	Use:   "characters",
	Short: "Process all characters locally",
	Long: `Process all characters locally`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		eve_processors.ProcessCharacters()
	},
}

func init() {
	rootCmd.AddCommand(charactersCmd)
}
