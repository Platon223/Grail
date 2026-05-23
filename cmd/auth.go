
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Platon223/Grail/internal/domain/auth"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Parent command that manages auth",	
}

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Connect your Gmail account",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("To use grail you need a Google Cloud OAuth credentials file.")
		fmt.Println("Steps:")
		fmt.Println("  1. Go to https://console.cloud.google.com")
		fmt.Println("  2. Create a project, enable Gmail API")
		fmt.Println("  3. Create OAuth 2.0 credentials (Desktop app)")
		fmt.Println("  4. Download the JSON file")
		fmt.Println("")
		fmt.Print("Path to credentials.json: ")

		var path string
		fmt.Scan(&path)

		if err := auth.SaveCredentials(path); err != nil {
			return err
		}

		config, err := auth.LoadConfig()
		if err != nil {
			return err
		}

		_, err = auth.GetClient(config)
		return err
	},
}


func init() {
	authCmd.AddCommand(authSetupCmd)
	rootCmd.AddCommand(authCmd)
}
