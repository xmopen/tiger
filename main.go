// Tiger main
package main

import "github.com/xmopen/tiger/internal/server"

func main() {
	// 启动
	svr := server.New(server.WithTCPProtocol(8848))
	if err := svr.ListenAndServerWithSignal(); err != nil {
		panic(err)
	}
}
