package outputs

import (
	"log/slog"

	"github.com/ncarlier/za/pkg/events"
)

const maxEntriesChanSize = 5000

// Manager used to handle outputs worker
type Manager struct {
	outputs []Output
}

// NewOutputsManager create an outputs manager
func NewOutputsManager(outputs []Output) (*Manager, error) {
	manager := Manager{
		outputs: make([]Output, len(outputs)),
	}

	for idx, output := range outputs {
		if err := output.Connect(); err != nil {
			slog.Error("unable to connect to the output writer", "error", err)
		} else {
			manager.outputs[idx] = output
		}
	}

	return &manager, nil
}

// SendEvent sent event to all outputs
func (m *Manager) SendEvent(event events.Event) {
	if event == nil {
		return
	}
	for _, out := range m.outputs {
		if err := out.SendEvent(event); err != nil {
			slog.Error("unable to send event to the output", "error", err)
		}
	}
}

// Shutdown stop the manager
func (m *Manager) Shutdown() {
	for _, out := range m.outputs {
		if err := out.Close(); err != nil {
			slog.Error("unable to close the output", "error", err)
		}
	}
}
