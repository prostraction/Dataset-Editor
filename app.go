package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SetDirectoryDialog() string {
	result, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           ".",
		Title:                      "Choose directory",
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: true})
	if len(result) > 0 {
		return result
	}
	return "Dialog cancelled"
}

func (a *App) StartMergeProcess(dir_in_1 string, dir_in_2 string, dir_merged string) bool {
	ProcessMerge(dir_in_1+"\\", dir_in_2+"\\", dir_merged+"\\")
	return true
}

func (a *App) StartCropProcess(dir_in_1 string, dir_result string, x int, y int) bool {
	ProcessCut(dir_in_1+"\\", dir_result+"\\", x, y)
	return true
}
