package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/syncromatics/idl-repository/pkg/config"
)

func init() {
	RootCmd.AddCommand(initCommand)
}

var initCommand = &cobra.Command{
	Use:   "init [name] [repository]",
	Short: "inits the config",
	Long:  "long stuff",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		configuration = new(config.Configuration)
		configuration.Name = args[0]
		configuration.Repository = args[1]

		f, err := os.Create(configLocation)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = configuration.Marshal(f)
		if err != nil {
			panic(err)
		}
	},
}
