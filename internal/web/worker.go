package web

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/homettp/homettp/internal/models"
)

func (app *app) worker() {
	for id := range app.queue {
		err := app.handleCall(id)
		if err != nil {
			app.errorLog.Print(err)
		}
	}
}

func (app *app) handleCall(id int64) error {
	call, err := app.callRepository.Find(id)
	if err != nil {
		return err
	}

	if call.Status != models.Pending {
		return fmt.Errorf("worker: invalid call model %v", id)
	}

	err = app.callRepository.Update(call, &models.Call{
		Status: models.InProgress,
	})
	if err != nil {
		return err
	}

	command, err := app.commandRepository.Find(call.CommandID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.commandTimeout)*time.Second)
	defer cancel()

	name := "/bin/sh"
	arg := "-c"

	if runtime.GOOS == "windows" {
		name = "cmd"
		arg = "/C"
	}

	cmd := exec.CommandContext(ctx, name, arg, strings.ReplaceAll(command.Value, models.PayloadVariable, call.Payload))
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		err = app.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: "command timed out",
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		err = app.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: err.Error(),
		})
		if err != nil {
			return err
		}
	} else {
		err = app.callRepository.Update(call, &models.Call{
			Status: models.Completed,
			Output: string(out),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
