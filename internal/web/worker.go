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

func (a *app) worker() {
	for id := range a.queue {
		err := a.handleCall(id)
		if err != nil {
			a.errorLog.Print(err)
		}
	}
}

func (a *app) handleCall(id int64) error {
	call, err := a.callRepository.Find(id)
	if err != nil {
		return err
	}

	if call.Status != models.Pending {
		return fmt.Errorf("worker: invalid call model %v", id)
	}

	err = a.callRepository.Update(call, &models.Call{
		Status: models.InProgress,
	})
	if err != nil {
		return err
	}

	command, err := a.commandRepository.Find(call.CommandID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.commandTimeout)*time.Second)
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
		err = a.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: "command timed out",
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		err = a.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: err.Error(),
		})
		if err != nil {
			return err
		}
	} else {
		err = a.callRepository.Update(call, &models.Call{
			Status: models.Completed,
			Output: string(out),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
