package client

import "sync"

var (
	lock         = &sync.Mutex{}
	initClientID = 0
)

// generateClientID 初始化客户端ID
func generateClientID() int {
	lock.Lock()
	defer lock.Unlock()
	initClientID++
	return initClientID
}
