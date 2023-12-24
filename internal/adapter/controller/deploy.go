package controller

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type DeployController interface {
	core.Controller
}

type deployController struct{}

func NewDeployController() DeployController {
	return &deployController{}
}

func (c *deployController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Print("deploy to vercel... ")
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")
	log, err := exec.Command("vercel --cwd " + outputDir + " --prod").Output()
	if err != nil {
		fmt.Printf("err: %v\ndeploy log:\n%v", err, log)
		return nil
	}
	fmt.Println("ok.")
	return nil
}
