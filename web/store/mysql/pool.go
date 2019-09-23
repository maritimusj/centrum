package mysqlStore

import (
	"sync"
)

var (
	storePool = sync.Pool{
		New: func() interface{} {
			return New()
		},
	}
)
