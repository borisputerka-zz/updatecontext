package cli

import (
	"fmt"
	"github.com/borisputerka/updatecontext/pkg/logger"
	"github.com/borisputerka/updatecontext/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// RootCmd function defines plugin function run
func RootCmd() *cobra.Command {
	o := &plugin.ConfigFlags{}
	cmd := &cobra.Command{
		Use:           "updatecontext",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger()
			log.Info("")

			err := viper.BindPFlags(cmd.Flags())
			if err != nil {
				return errors.Cause(err)
			}

			err = o.Complete()
			if err != nil {
				return errors.Cause(err)
			}

			err = o.RunPlugin()
			if err != nil {
				return errors.Cause(err)
			}

			log.Info("")

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

// InitAndExecute function is called in main
func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
