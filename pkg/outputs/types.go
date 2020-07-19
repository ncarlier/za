package outputs

import (
	"github.com/ncarlier/trackr/pkg/model"
)

// Output writer
type Output interface {
	// Connect to the Output
	Connect() error
	// Close any connections to the Output
	Close() error
	// Write page views to the Output
	Write(views []*model.PageView) error
}
