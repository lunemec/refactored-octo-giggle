package cmd

import (
	"fmt"
	"os"

	"refactored-octo-giggle/pkg/api"

	log "github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config structure is populated from config file.
// This might be used as app-wide configuration.
type Config struct {
	API api.Config
}

var (
	cfgFile string
	config  Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "refactored-octo-giggle",
	Short: "GlobalWebIndex Engineering Challenge server",
	Long:  ``,
	Run:   runAPI,
}

func runAPI(cmd *cobra.Command, args []string) {
	log.Info("API Running", "addr", config.API.Addr())
	err := api.RunServer(config.API)
	if err != nil {
		log.Error("Unable to start server", "err", err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "app.toml", "config file (default is app.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatal("Error parsing config file", "path", cfgFile, "err", err)
		}
		log.Info("Using config file:", "path", viper.ConfigFileUsed())
	}
}
