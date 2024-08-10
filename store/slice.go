package store

import (
	"github.com/opentoys/ledisdb/store/driver"
)

type Slice interface {
	driver.ISlice
}
