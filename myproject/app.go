package main

import (
	"context"
	"fmt"
	"time"
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

	go a.ConvertFile()

	fmt.Println("do stuff")
}

// Greet returns a greeting for the given name

func (a *App) ConvertFiles(inputPath string, outputPath string) {
	go a.doStuff(inputPath, outputPath)
	time.Sleep(1 * time.Second)
	a.PercentBar()
}
