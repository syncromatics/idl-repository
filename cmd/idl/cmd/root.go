package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/syncromatics/idl-repository/pkg/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configLocation string
	configuration  *config.Configuration
)

func init() {
	RootCmd.PersistentFlags().StringVar(&configLocation, "config", "./idl.yaml", "The location of the idl configuration yaml file")
}

var RootCmd = &cobra.Command{
	Use:   "idl",
	Short: "idl stores and fetches all sorts of idls",
	Long:  `long explanation here`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func initConfig() error {
	configuration = new(config.Configuration)

	dat, err := os.Open(configLocation)
	if os.IsNotExist(err) {
		return errors.New("The idl.yaml file does not exist, run 'idl init [name] [repository]' to create it")
	}

	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	err = configuration.UnMarshal(bufio.NewReader(dat))
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	return nil
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
