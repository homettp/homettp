//go:build unix

package web

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/homettp/homettp/internal/config"
	"github.com/homettp/homettp/internal/models"
)

type memoryCallRepository struct {
	models.CallRepository
	calls map[int64]*models.Call
}

func (mcr *memoryCallRepository) Find(id int64) (*models.Call, error) {
	call, ok := mcr.calls[id]
	if !ok {
		return nil, models.ErrNoRecord
	}

	copied := *call

	return &copied, nil
}

func (mcr *memoryCallRepository) Update(call, newCall *models.Call) error {
	call.Status = newCall.Status
	call.Output = newCall.Output
	mcr.calls[call.ID] = call

	return nil
}

type memoryCommandRepository struct {
	models.CommandRepository
	command *models.Command
}

func (mcr *memoryCommandRepository) Find(int) (*models.Command, error) {
	return mcr.command, nil
}

func runTestCall(t *testing.T, timeout int, value, payload string) *models.Call {
	t.Helper()

	calls := &memoryCallRepository{calls: map[int64]*models.Call{
		1: {ID: 1, CommandID: 1, Status: models.Pending, Payload: payload},
	}}

	a := &app{
		appConfig:         &config.Config{CommandTimeout: timeout},
		callRepository:    calls,
		commandRepository: &memoryCommandRepository{command: &models.Command{ID: 1, Value: value}},
	}

	err := a.handleCall(t.Context(), 1)
	if err != nil {
		t.Fatal(err)
	}

	return calls.calls[1]
}

func TestHandleCallPayloadIsNotParsed(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "pwned")
	payload := `x; touch ` + marker + ` $(touch ` + marker + `) "'`

	call := runTestCall(t, 5, "printf %s %p", payload)

	if call.Status != models.Completed || call.Output != payload {
		t.Fatalf("expected completed with payload output, got %s %q", call.Status, call.Output)
	}

	if _, err := os.Stat(marker); err == nil {
		t.Fatal("payload was executed by the shell")
	}
}

func TestHandleCallTimeoutKillsChildren(t *testing.T) {
	start := time.Now()
	call := runTestCall(t, 1, "sleep 30; echo done", "")

	if call.Status != models.Failed || call.Output != "command timed out" {
		t.Fatalf("expected timed out failure, got %s %q", call.Status, call.Output)
	}

	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("expected the timeout to stop the command, took %v", elapsed)
	}
}

func TestHandleCallEmptyPayloadIsNoArgument(t *testing.T) {
	call := runTestCall(t, 5, `printf '[%s]' %p; echo " $#"; set -- %p; echo $#`, "")

	if call.Status != models.Completed || call.Output != "[] 0\n0\n" {
		t.Fatalf("expected no argument, got %s %q", call.Status, call.Output)
	}
}

func TestHandleCallBackgroundChildCompletes(t *testing.T) {
	call := runTestCall(t, 5, "sleep 3 & echo started", "")

	if call.Status != models.Completed || call.Output != "started\n" {
		t.Fatalf("expected completed with output, got %s %q", call.Status, call.Output)
	}
}

func TestHandleCallFailedKeepsOutput(t *testing.T) {
	call := runTestCall(t, 5, "echo unreachable; echo oops >&2; exit 3", "")

	if call.Status != models.Failed || call.Output != "unreachable\noops\nexit status 3" {
		t.Fatalf("expected failed with output, got %s %q", call.Status, call.Output)
	}
}

func TestWorkerInterruptsQueuedCallsAfterClose(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "ran")

	calls := &memoryCallRepository{calls: map[int64]*models.Call{
		1: {ID: 1, CommandID: 1, Status: models.Pending},
	}}

	a := &app{
		appConfig:         &config.Config{CommandTimeout: 5},
		callRepository:    calls,
		commandRepository: &memoryCommandRepository{command: &models.Command{ID: 1, Value: "touch " + marker}},
		queue:             make(chan int64, 1),
	}

	a.queue <- 1
	a.closeQueue()
	a.worker(t.Context())

	call := calls.calls[1]
	if call.Status != models.Failed || call.Output != "command interrupted" {
		t.Fatalf("expected interrupted failure, got %s %q", call.Status, call.Output)
	}

	if _, err := os.Stat(marker); err == nil {
		t.Fatal("queued call ran after the queue was closed")
	}
}
