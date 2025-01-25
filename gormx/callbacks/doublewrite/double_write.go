package doublewrite

import (
	"go.uber.org/atomic"
	"gorm.io/gorm"
)

type Callback struct {
	pattern *atomic.String
}

func (c *Callback) doubleWriteCreate() func(db *gorm.DB) {
	return func(db *gorm.DB) {
	}
}
