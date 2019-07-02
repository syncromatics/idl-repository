package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/syncromatics/idl-repository/internal/repository"
	"github.com/syncromatics/idl-repository/internal/storage"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	port            *int
	storageDiretory *string
)

func init() {
	port = rootCmd.Flags().IntP("port", "p", 80, "The port to host the server on")
	storageDiretory = rootCmd.Flags().StringP("storage", "s", "/var/idl-repository", "The storage location for modules")
}

var rootCmd = &cobra.Command{
	Use:   "idl-repository",
	Short: "idl-repository stores all sorts of idls",
	Long:  `long explanation here`,
	Run: func(cmd *cobra.Command, args []string) {
		settings := &repository.Settings{
			Port: *port,
		}

		storage, err := storage.NewFileStorage(*storageDiretory)
		if err != nil {
			panic(err)
		}
		server := repository.NewServer(settings, storage)

		ctx, cancel := context.WithCancel(context.Background())
		grp, ctx := errgroup.WithContext(ctx)

		grp.Go(server.Run(ctx))

		// Wait for SIGINT/SIGTERM
		waiter := make(chan os.Signal, 1)
		signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-waiter:
		case <-ctx.Done():
		}
		cancel()
		if err := grp.Wait(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
