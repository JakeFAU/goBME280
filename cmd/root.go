package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var feedPrefix string
var units string
var address int
var baseURL string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goBME280",
	Short: "An application that sends BME280 data to adafruitIO",
	Long: `An application that sends BME280 data to adafruitIO.
  It sends temperature, humidity and pressure data to a specific feed`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goBME280.yaml)")
	rootCmd.PersistentFlags().StringVarP(&feedPrefix, "prefix", "p", "", "The feed prefix")
	rootCmd.PersistentFlags().StringVarP(&units, "units", "u", "english", "The units to use")
	rootCmd.PersistentFlags().IntVarP(&address, "address", "a", 0x77, "The address of the BME sensor")
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "b", "https://io.adafruit.com/", "The base URL")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".goBME280" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goBME280")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
