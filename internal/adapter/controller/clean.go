package controller

import (
	"fmt"
	"os"
	"path"

	"github.com/c18t/nippo-cli/internal/core"
	"github.com/spf13/cobra"
)

type CleanController interface {
	core.Controller
}
type cleanController struct{}

func NewCleanController() CleanController {
	return &cleanController{}
}

func (c *cleanController) Exec(cmd *cobra.Command, args []string) error {
	fmt.Print("clean cache files... ")

	clearNippoCache()
	clearBuildCache()

	fmt.Println("ok.")
	return nil
}

func clearNippoCache() error {
	outputDir := path.Join(core.Cfg.GetCacheDir(), "md")
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		err = os.Remove(path.Join(outputDir, fileName))
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
