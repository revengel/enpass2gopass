// main package
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger  = logrus.New()
	version string
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "0201-150405"
	customFormatter.FullTimestamp = false

	logger.SetFormatter(customFormatter)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.WarnLevel)
}

func main() {
	var err error
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	a := &app{
		ctx:    ctx,
		logger: logger,
	}

	defer a.Close()

	rootCmd := &cobra.Command{
		Use:   "enpass2gopass",
		Short: `enpass to gopass importer`,
	}

	rootCmd.PersistentFlags().StringP("log-level", "", logrus.InfoLevel.String(), "log level")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "enable debug log level")

	versionCmd := &cobra.Command{
		Use:     "version",
		Short:   "show version",
		Aliases: []string{"ver"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(getVersion().String())
		},
	}

	importCmd := &cobra.Command{
		Use:     "import",
		Short:   "Import command",
		Aliases: []string{},
		PreRunE: a.Before,
		RunE:    a.Import,
	}

	importCmd.PersistentFlags().StringP("prefix", "", "", "destination storage path prefix")
	importCmd.PersistentFlags().BoolP("dry-run", "", false, "do not make changes")
	importCmd.PersistentFlags().StringP("source-provider", "", EnpassJsonSourceType, "source provider")
	importCmd.PersistentFlags().StringP("source-enpass-json-path", "", "", "source enpass json path")
	importCmd.PersistentFlags().StringP("destination-provider", "", GopassDestinationType, "destination provider")

	rootCmd.AddCommand(versionCmd, importCmd)

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		logger.Fatalf("cannot run root command: %s", err.Error())
	}
}
