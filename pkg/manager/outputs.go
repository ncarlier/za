package manager

import (
	"sync"
	"time"

	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/outputs"
)

const maxEntriesChanSize = 5000

// Outputs is the outputs worker
type Outputs struct {
	config    *config.Config
	quit      chan struct{}
	entries   chan model.PageView
	outputs   []outputs.Output
	waitGroup sync.WaitGroup
}

// NewOutputsManager create an ouputs manager
func NewOutputsManager(conf *config.Config) (*Outputs, error) {
	manager := Outputs{
		config:  conf,
		quit:    make(chan struct{}),
		entries: make(chan model.PageView, maxEntriesChanSize),
		outputs: make([]outputs.Output, len(conf.Outputs)),
	}

	for idx, output := range conf.Outputs {
		out := *output
		if err := out.Connect(); err != nil {
			logger.Error.Println("unable to connect to the output writer: %v", err)
		} else {
			manager.outputs[idx] = out
		}
	}

	manager.waitGroup.Add(len(manager.outputs))
	go manager.run()

	return &manager, nil
}

// Send page view to all outputs
func (m *Outputs) Send(view model.PageView) {
	m.entries <- view
}

// Shutdown stop the manager
func (m *Outputs) Shutdown() {
	close(m.quit)
	m.waitGroup.Wait()
}

func (m *Outputs) closeOutputs() {
	for _, out := range m.outputs {
		if err := out.Close(); err != nil {
			logger.Error.Println("unable to close output:", err)
		}
		m.waitGroup.Done()
	}
}

func (m *Outputs) run() {
	var batch []*model.PageView
	batchSize := 0
	maxWait := time.NewTimer(m.config.Global.BatchInterval.Duration)

	defer func() {
		if batchSize > 0 {
			m.write(batch)
		}
		m.closeOutputs()
	}()

	for {
		select {
		case <-m.quit:
			return
		case entry := <-m.entries:
			batch = append(batch, &entry)
			batchSize++
			if batchSize >= m.config.Global.BatchSize {
				m.write(batch)
				batch = []*model.PageView{}
				batchSize = 0
				maxWait.Reset(m.config.Global.BatchInterval.Duration)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				m.write(batch)
				batch = []*model.PageView{}
				batchSize = 0
			}
			maxWait.Reset(m.config.Global.BatchInterval.Duration)
		}
	}
}

func (m *Outputs) write(entries []*model.PageView) {
	for _, out := range m.outputs {
		if err := out.Write(entries); err != nil {
			logger.Error.Println("unable to write entries:", err)
		}
	}
}
