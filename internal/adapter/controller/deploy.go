package controller

import (
	"fmt"
	"io"
	"os"
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
	dataDir := path.Join(core.Cfg.GetDataDir(), "output")
	outputDir := path.Join(core.Cfg.GetCacheDir(), "output")

	files, err := os.ReadDir(dataDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		src, err := os.Open(path.Join(dataDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer src.Close()
		dest, err := os.Create(path.Join(outputDir, file.Name()))
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer dest.Close()
		_, err = io.Copy(dest, src)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	log, err := exec.Command("vercel", "--cwd", outputDir, "--prod").Output()
	if err != nil {
		fmt.Printf("err: %v\ndeploy log:\n%v", err, log)
		return nil
	}
	fmt.Println("ok.")
	return nil
}
