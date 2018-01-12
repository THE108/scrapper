package result

import (
	"sync"

	"scrapper/info"
)

type resultHolder struct {
	mu     sync.Mutex
	result map[string]info.PageInfo
}

func NewResultHolder() *resultHolder {
	return &resultHolder{
		result: make(map[string]info.PageInfo),
	}
}

func (rh *resultHolder) Add(url string, pageInfo info.PageInfo) {
	rh.mu.Lock()
	rh.result[url] = pageInfo
	rh.mu.Unlock()
}

func (rh *resultHolder) Get() map[string]info.PageInfo {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	return rh.result
}
