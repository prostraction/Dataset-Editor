package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	lastDir        string
	isTaskFinished bool
	taskProcent    float32
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
	if a.lastDir == "" {
		a.lastDir = "."
	}
	result, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           a.lastDir,
		Title:                      "Choose directory",
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: true})
	if len(result) > 0 {
		a.lastDir = result
		return result
	}
	return "Dialog cancelled"
}

func (a *App) IsTaskFinished() bool {
	return a.isTaskFinished
}

func (a *App) StartMergeProcess(dir_in_1 string, dir_in_2 string, dir_merged string) bool {
	a.isTaskFinished = false
	fmt.Println("Merge called")
	ProcessMerge(dir_in_1+"\\", dir_in_2+"\\", dir_merged+"\\")
	a.isTaskFinished = true
	return true
}

func (a *App) StartCropProcess(dir_in_1 string, dir_result string, x int, y int) bool {
	a.isTaskFinished = false
	fmt.Println("Crop called")
	ProcessCut(dir_in_1+"\\", dir_result+"\\", x, y)
	a.isTaskFinished = true
	return true
}

func (a *App) StartProcessBrightness(dir_in_1 string, dir_result string, factor int) bool {
	a.isTaskFinished = false
	fmt.Println("Brightness called")
	ProcessBrightness(dir_in_1+"\\", dir_result+"\\", factor)
	a.isTaskFinished = true
	return true
}
