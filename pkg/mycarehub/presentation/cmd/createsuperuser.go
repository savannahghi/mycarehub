package cmd

import (
	"context"
	"log"
	"os"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/cmd/service"
	injector "github.com/savannahghi/mycarehub/wire"
	"github.com/spf13/cobra"
)

// createsuperuserCmd represents the createsuperuser command
var createsuperuserCmd = &cobra.Command{
	Use:   "createsuperuser",
	Short: "Creates the initial user for Mycarehub",
	Long: `The initial user is assigned to the initial default organization, program and facility.
	They must be a staff user`,
	Run: func(cmd *cobra.Command, args []string) {
		useCases, err := injector.InitializeUseCases(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		mycarehubService := service.NewMyCareHubCmdInterfaces(*useCases)

		if err := mycarehubService.CreateSuperUser(cmd.Context(), os.Stdin); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(createsuperuserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createsuperuserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createsuperuserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
