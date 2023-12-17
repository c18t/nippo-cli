package root

import (
	"fmt"
	"runtime/debug"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
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
	core.Cfg = &core.Config{}
	core.Cfg.LoadConfig(Params.ConfigFile)
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
