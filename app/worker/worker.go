// Simple in-memory worker & queue
package worker

import (
	"context"
	"fmt"

	"github.com/fatih/color"

	"github.com/pocketbase/pocketbase/core"
)

type Work interface {
	Name() string
	Execute(context.Context) error
}

var started = false
var queue = make(chan Work, 100)
var ctx = context.Background()

var app core.App

func SetApplication(application core.App) {
	app = application
}

func StartWorker(application core.App) {
	if started {
		return
	}
	started = true
	SetApplication(application)

	bold := color.New(color.Bold).Add(color.FgGreen)
	red := color.New(color.FgRed)

	bold.Println("> Worker started")
	go func() {
		fmt.Println("  - Waiting for jobs...")
		for {
			select {
			case <-ctx.Done():
				bold.Println("> Worker stopped")
				return
			case work := <-queue:
				bold.Println("> Job started:", work.Name())
				if err := work.Execute(context.WithoutCancel(ctx)); err != nil {
					red.Println("> Worker error:", err)
				}
			}
		}
	}()
}

func StopWorker() {
	ctx.Done()
}

func SubmitAndWait(work Work) {
	queue <- work
}

func SubmitAndForget(work Work) {
	go func() {
		queue <- work
	}()
}
