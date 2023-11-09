package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "symbiote",
	Short: "symbiote, Elysian swiss army knife",
	Long:  `Symbiote is a utility tool for managing and interacting with your cloud infrastructure. `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ")
		fmt.Println("       Symbiote       ")
		fmt.Println("    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ")
		fmt.Println("    ⠀⠀⢀⠤⠒⠒⠀⠒⠂⢄⡀⠀⠀⠀⠀⠀⠀ ")
		fmt.Println("    ⠀⡜⢀⠀⠀⠀⠀⠀⠀⢀⡈⢆⠀⠀⠀⠀⠀ ")
		fmt.Println("    ⠸⠀⠃⡄⠀⠀⠀⠀⠀⡸⢱⠀⠄⠀⠀⠀⠀ ")
		fmt.Println("    ⡇⢸⠀⠀⠄⠀⠀⠀⠜⠀⠀⠂⠄⠀⠀⠀⠀ ")
		fmt.Println("    ⠇⢨⠀⠀⠈⢚⠸⠂⠀⠀⠀⡄⡄⠀⠀⠀⠀ ")
		fmt.Println("    ⠐⢆⢂⠀⠀⠈⠀⠢⢀⣀⢎⡄⠀⠀⠀⡀⠀ ")
		fmt.Println("    ⠰⡘⣷⢶⣿⣾⣧⣿⣶⠶⡻⣠⠃⠀⠀⠘⡄ ")
		fmt.Println("    ⠀⠑⣿⣄⣿⠾⠷⡟⠥⡴⠡⠃⠀⠀⢀⠔⡀ ")
		fmt.Println("    ⠀⠀⠘⢟⣷⣠⣄⡀⠀⠈⢇⠠⠂⠉⡠⠄⠀ ")
		fmt.Println("    ⠀⠀⠀⠘⠻⣿⠶⡟⣄⡀⠀⠀⢀⠎⠀⠀⠀ ")
		fmt.Println("    ⠀⠀⠀⠀⠑⠤⠠⠤⠃⠀⠉⠉⠀⠀⠀⠀⠀ ")
		fmt.Println("    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ")
		fmt.Println("    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ")
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/symbiote/")
	err := viper.ReadInConfig()
	if err != nil {
		// TODO: Handle errors (e.g. generate config file)
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
