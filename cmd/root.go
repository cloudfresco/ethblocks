package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ClientAddr string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ethblocks",
	Short: "ethblocks Ethereum block explorer cli interface",
	Long:  `ethblocks Ethereum block explorer cli interface`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	v := viper.New()
	v.SetConfigFile(".env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	ClientAddr = v.GetString("ETHBLOCKS_CLIENT_ADDRESS")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
