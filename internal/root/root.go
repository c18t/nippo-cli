package root

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

type RootParams struct {
	ConfigFile    string
	Version       bool
	LicenseNotice bool
}

var Version string
var Params RootParams

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	if Params.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(Params.ConfigFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		defaultConfigDir := path.Join(home, ".config")
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" || !path.IsAbs(configDir) {
			configDir = defaultConfigDir
		}
		configPath := path.Join(configDir, "nippo")

		// Search config in home directory with name ".nippo" (without extension).
		viper.AddConfigPath(configPath)
		viper.SetConfigType("toml")
		viper.SetConfigName(".nippo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func CreateCmdFunc() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// show nippo-cli version
		if Params.Version {
			if Version != "" {
				// go build -ldflags "-X 'main.version=vx.x.x'"
				fmt.Println(Version)
			} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
				// go install version tag
				fmt.Println(buildInfo.Main.Version)
			} else {
				// unknown version
				fmt.Println("(unknown)")
			}
			return nil
		}

		// show license notice
		if Params.LicenseNotice {
			fmt.Println("set --license-notice")
			return nil
		}

		// show help
		cmd.Help()
		return nil
	}
}
