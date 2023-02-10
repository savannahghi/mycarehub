package cmd

import (
	"context"
	"log"
	"os"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/cmd/service"
	injector "github.com/savannahghi/mycarehub/wire"
	"github.com/spf13/cobra"
)

func mycarehubCommands() []*cobra.Command {
	useCases, err := injector.InitializeUseCases(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	mycarehubService := service.NewMyCareHubCmdInterfaces(*useCases)

	var createsuperuserCmd = &cobra.Command{
		Use:   "createsuperuser",
		Short: "Creates the initial user for Mycarehub",
		Long: `The initial user is assigned to the initial default organization, program and facility.
			They must be a staff user`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.CreateSuperUser(cmd.Context(), os.Stdin); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

	var loadFacilitiesCmd = &cobra.Command{
		Use:   "loadfacilities",
		Short: "Creates the initial facilities Mycarehub",
		Long:  `The initial facilities are created and saved in the database.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadFacilities(cmd.Context(), "data/facilities.csv"); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

	var loadOrganisatioAndProgramCmd = &cobra.Command{
		Use:   "loadorganisationandprogram",
		Short: "Creates the initial organisation and a program associated with it",
		Long:  `The organisation is first created then a program is associated with that organisation`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadOrganisatioAndProgram(cmd.Context(), "data/organisation.json", "data/program.json"); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

	return []*cobra.Command{
		loadOrganisatioAndProgramCmd,
		loadFacilitiesCmd,
		createsuperuserCmd,
	}

}

func init() {
	commands := mycarehubCommands()
	for _, command := range commands {
		rootCmd.AddCommand(command)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createsuperuserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createsuperuserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
