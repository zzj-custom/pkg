package pAntsPool

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

var (
	size              int = 128
	asyncTaskPoolOnce sync.Once
	asyncTaskPool     *ants.Pool
)

func InitAsyncTaskPool(s int) (*ants.Pool, error) {
	if asyncTaskPool != nil {
		return nil, fmt.Errorf("async task pool already inited")
	}
	if s > 0 {
		size = s
	}
	return AsyncTaskPool(), nil
}

func AsyncTaskPool() *ants.Pool {
	asyncTaskPoolOnce.Do(func() {
		opts := ants.WithNonblocking(true)
		asyncTaskPool, _ = ants.NewPool(size, opts)
	})
	return asyncTaskPool
}
