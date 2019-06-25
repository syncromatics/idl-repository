package cmd

import (
	"os"

	"github.com/syncromatics/idl-repository/pkg/client"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	packageVersion *semver.Version
)

func init() {
	RootCmd.AddCommand(pushCommand)
}

var pushCommand = &cobra.Command{
	Use:   "push [version]",
	Short: "push the provides to the repository",
	Long:  "long stuff",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires version")
		}

		var err error
		packageVersion, err = semver.NewVersion(args[0])
		if err != nil {
			return errors.Wrap(err, "invalid version")
		}

		err = initConfig()
		if err != nil {
			return errors.Wrap(err, "invalid config")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Push(client.PushOptions{
			Configuration: configuration,
			Version:       packageVersion,
		})
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
	},
}
