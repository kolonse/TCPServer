package TCPServer

import (
	"github.com/kolonse/function"
	"net"
)

type TCPConn struct {
	net.Conn
	recvCB  *function.Function
	writeCB *function.Function
}

// TCP 连接服务接口
func (tc *TCPConn) Server() {
	buff := make([]byte, 10000)
	for {
		n, err := tc.Read(buff)
		if err != nil {
			// 调用对端关闭连接通知函数
			if tc.recvCB != nil {
				tc.recvCB.Call(tc, nil, err)
			}
			break
		}

		// 调用收到数据处理函数
		if tc.recvCB != nil && tc.recvCB.IsValid() {
			tc.recvCB.Call(tc, buff[:n], nil)
		}
	}
	tc.Close()
}

func NewTCPConn(conn net.Conn, recvCB *function.Function, writeCB *function.Function) *TCPConn {
	return &TCPConn{
		Conn:    conn,
		recvCB:  recvCB,
		writeCB: writeCB,
	}
}
