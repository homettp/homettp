package web

import (
	"errors"
	"fmt"

	"github.com/homettp/homettp/internal/models"
)

func (app *App) worker() {
	for id := range app.queue {
		err := app.handleCall(id)
		if err != nil {
			app.errorLog.Print(err)
		}
	}
}

func (app *App) handleCall(id int64) error {
	call, err := app.callRepository.Find(id)
	if err != nil {
		return err
	}

	if call.Status != models.Pending {
		return errors.New(fmt.Sprintf("worker: invalid call model %v", id))
	}

	// TODO: Handle call

	return nil
}
