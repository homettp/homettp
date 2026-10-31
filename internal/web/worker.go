package web

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/homettp/homettp/internal/models"
)

const (
	payloadEnvName = "HOMETTP_PAYLOAD"
	waitDelay      = time.Second
)

func (a *app) worker(ctx context.Context) {
	for id := range a.queue {
		var err error

		if a.isQueueClosed() {
			err = a.interruptCall(id)
		} else {
			err = a.handleCall(ctx, id)
		}

		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			a.errorLog.Print(err)
		}
	}
}

func (a *app) isQueueClosed() bool {
	a.queueMu.RLock()
	defer a.queueMu.RUnlock()

	return a.queueClosed
}

func (a *app) interruptCall(id int64) error {
	call, err := a.callRepository.Find(id)
	if errors.Is(err, models.ErrNoRecord) {
		return nil
	}

	if err != nil {
		return err
	}

	if call.Status != models.Pending {
		return nil
	}

	err = a.callRepository.Update(call, &models.Call{
		Status: models.Failed,
		Output: "command interrupted",
	})
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}

	return nil
}

func (a *app) enqueue(id int64) error {
	a.queueMu.RLock()
	defer a.queueMu.RUnlock()

	if a.queueClosed {
		return ErrQueueClosed
	}

	select {
	case a.queue <- id:
		return nil
	default:
		return ErrQueueFull
	}
}

func (a *app) closeQueue() {
	a.queueMu.Lock()
	defer a.queueMu.Unlock()

	if a.queueClosed {
		return
	}

	a.queueClosed = true
	close(a.queue)
}

func (a *app) waitWorkers(ctx context.Context, cancelWorkers context.CancelFunc) {
	finished := make(chan struct{})

	go func() {
		a.workers.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case <-ctx.Done():
		a.infoLog.Print("Canceling remaining calls")
		cancelWorkers()
		<-finished
	}
}

func (a *app) failInterruptedCalls() error {
	commands, err := a.commandRepository.FindAll()
	if err != nil {
		return err
	}

	for _, command := range commands {
		calls, err := a.callRepository.FindAllByCommand(command)
		if err != nil {
			return err
		}

		for _, call := range calls {
			if !call.Status.IsRunning() {
				continue
			}

			err = a.callRepository.Update(call, &models.Call{
				Status: models.Failed,
				Output: "command interrupted",
			})
			if err != nil && !errors.Is(err, models.ErrNoRecord) {
				return err
			}
		}
	}

	return nil
}

func (a *app) handleCall(ctx context.Context, id int64) error {
	call, err := a.callRepository.Find(id)
	if err != nil {
		return err
	}

	if call.Status != models.Pending {
		return fmt.Errorf("%w: %d", ErrInvalidCall, id)
	}

	command, err := a.commandRepository.Find(call.CommandID)
	if err != nil {
		return errors.Join(err, a.callRepository.Update(call, &models.Call{
			Status: models.Failed,
			Output: err.Error(),
		}))
	}

	err = a.callRepository.Update(call, &models.Call{
		Status: models.InProgress,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.appConfig.CommandTimeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, shellName, shellArg, strings.ReplaceAll(command.Value, models.PayloadVariable, payloadArgument(call.Payload)))
	cmd.Env = append(os.Environ(), payloadEnvName+"="+call.Payload)
	cmd.WaitDelay = waitDelay
	configureProcess(cmd)

	out, err := cmd.Output()

	if errors.Is(err, exec.ErrWaitDelay) && ctx.Err() == nil {
		err = nil
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
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
			Output: failedOutput(out, err),
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
