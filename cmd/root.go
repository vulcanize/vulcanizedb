package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vulcanize/vulcanizedb/pkg/config"
)

var (
	cfgFile        string
	databaseConfig config.Database
	ipc            string
)

var rootCmd = &cobra.Command{
	Use:              "vulcanizedb",
	PersistentPreRun: database,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func database(cmd *cobra.Command, args []string) {
	ipc = viper.GetString("client.ipcpath")
	databaseConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
	}
	viper.Set("database.config", databaseConfig)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "environment/public.toml", "config file location")
	rootCmd.PersistentFlags().String("database-name", "vulcanize_public", "database name")
	rootCmd.PersistentFlags().Int("database-port", 5432, "database port")
	rootCmd.PersistentFlags().String("database-hostname", "localhost", "database hostname")
	rootCmd.PersistentFlags().String("client-ipcPath", "", "location of geth.ipc file")

	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("database-name"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("database-port"))
	viper.BindPFlag("database.hostname", rootCmd.PersistentFlags().Lookup("database-hostname"))
	viper.BindPFlag("client.ipcPath", rootCmd.PersistentFlags().Lookup("client-ipcPath"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".vulcanizedb")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	}
}
