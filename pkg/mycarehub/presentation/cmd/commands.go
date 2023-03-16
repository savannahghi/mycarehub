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

	var loadOrganisationCmd = &cobra.Command{
		Use:   "loadorganisation",
		Short: "Creates the initial organisation",
		Long:  `The initial organisation is created and saved in the database.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadOrganisation(cmd.Context(), "data/organisation.json"); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

	var loadProgramCmd = &cobra.Command{
		Use:   "loadprogram",
		Short: "Creates the initial program and links it to the initial organisation",
		Long:  `The initial program is created and saved in the database.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadProgram(cmd.Context(), "data/program.json", os.Stdin); err != nil {
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

	var linkFacilityToProgramCmd = &cobra.Command{
		Use:   "linkfacilitytoprogram",
		Short: "links a facility to the initial program",
		Long:  `The facility selected will be linked to the initial program`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LinkFacilityToProgram(cmd.Context(), os.Stdout); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

	var loadSecurityQuestionsCmd = &cobra.Command{
		Use:   "loadsecurityquestions",
		Short: "Creates the system's security questions",
		Long:  `The security questions created are for the PRO and CONSUMER application`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadSecurityQuestions(cmd.Context(), "data/securityquestions.json"); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}
	var loadTermsOfServiceCmd = &cobra.Command{
		Use:   "loadtermsofservice",
		Short: "Loads terms of service from a local .txt and loads them to the database",
		Long:  `The previous terms should be invalidated. The new term should be active`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := mycarehubService.LoadTermsOfService(cmd.Context(), os.Stdin); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	}

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

	return []*cobra.Command{
		loadOrganisationCmd,
		loadProgramCmd,
		loadFacilitiesCmd,
		linkFacilityToProgramCmd,
		loadTermsOfServiceCmd,
		loadSecurityQuestionsCmd,
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
