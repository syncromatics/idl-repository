package cmd

import (
	"os"

	"github.com/syncromatics/idl-repository/pkg/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(pullCommand)
}

var pullCommand = &cobra.Command{
	Use:   "pull",
	Short: "pull it all in",
	Long:  "long stuff",
	Args: func(cmd *cobra.Command, args []string) error {
		err := initConfig()
		if err != nil {
			return errors.Wrap(err, "invalid config")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Pull(client.PullOptions{
			Configuration: configuration,
		})
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}
	},
}
