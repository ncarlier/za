package outputs

import (
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/model"
)

const maxEntriesChanSize = 5000

// Manager used to handle outputs worker
type Manager struct {
	outputs []model.Output
}

// NewOutputsManager create an outputs manager
func NewOutputsManager(outputs []model.Output) (*Manager, error) {
	manager := Manager{
		outputs: make([]model.Output, len(outputs)),
	}

	for idx, output := range outputs {
		if err := output.Connect(); err != nil {
			logger.Error.Printf("unable to connect to the output writer: %v", err)
		} else {
			manager.outputs[idx] = output
		}
	}

	return &manager, nil
}

// SendPageView page view to all outputs
func (m *Manager) SendPageView(view model.PageView) {
	for _, out := range m.outputs {
		if err := out.SendPageView(view); err != nil {
			logger.Error.Println("unable to send page view to the output:", err)
		}
	}
}

// Shutdown stop the manager
func (m *Manager) Shutdown() {
	for _, out := range m.outputs {
		if err := out.Close(); err != nil {
			logger.Error.Println("unable to close the output:", err)
		}
	}
}
