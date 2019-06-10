package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(nodeSyncCmd)
}

var cfg *viper.Viper

func initConfig() {
	pwd, _ := os.Getwd()
	cfg = viper.New()
	cfg.SetEnvPrefix("wwsync")
	cfg.SetDefault("config_path", pwd)
	cfg.SetDefault("system", "")
	cfg.SetDefault("sqs.timeout", 20)
	cfg.AutomaticEnv()
}

var rootCmd = &cobra.Command{
	Use:   "warewulf-sync",
	Short: "Warewulf Sync is a tool for building and syncronizing a warewulf configuration",
	Long: `A tool for generating and synchronizing a warewulf database from either declarative
			configuration files and/or another api or database.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
