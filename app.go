package main

import (
	"context"
	"fmt"

	"dataset/internal/database"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	db             database.Database
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

func (a *App) GetProcentValue() float32 {
	return a.taskProcent
}

func (a *App) IsTaskFinished() bool {
	return a.isTaskFinished
}

func (a *App) StartMergeProcess(dir_in_1 string, dir_in_2 string, dir_merged string) bool {
	a.taskProcent = 0.00
	a.isTaskFinished = false
	fmt.Println("Merge called")
	a.ProcessMerge(dir_in_1+"\\", dir_in_2+"\\", dir_merged+"\\")
	a.isTaskFinished = true
	return true
}

func (a *App) StartCropProcess(dir_in_1 string, dir_result string, x int, y int) bool {
	a.taskProcent = 0.00
	a.isTaskFinished = false
	fmt.Println("Crop called")
	a.ProcessCut(dir_in_1+"\\", dir_result+"\\", x, y)
	a.isTaskFinished = true
	return true
}

func (a *App) StartBrightnessProcess(dir_in_1 string, dir_result string, factor float64) bool {
	a.taskProcent = 0.00
	a.isTaskFinished = false
	fmt.Println("Brightness called")
	a.ProcessBrightness(dir_in_1+"\\", dir_result+"\\", factor)
	a.isTaskFinished = true
	return true
}

// TO DO: add URI to args
func (a *App) StartDotsToDbProccess(dir_in_1 string) bool {
	a.taskProcent = 0.00
	a.isTaskFinished = false
	fmt.Println("Dots to DB called")
	var err error
	if !a.db.IsInit {
		err = a.db.Init()
	}
	if err != nil {
		fmt.Println(err.Error())
		a.isTaskFinished = true
		return true
	}
	a.ProcessDotsToDB(dir_in_1 + "\\")
	a.isTaskFinished = true
	return true
}
