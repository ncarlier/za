package output

import "github.com/ncarlier/trackr/pkg/model"

// Writer is an interface used to write metrics
type Writer interface {
	Write(hit model.PageView)
}
