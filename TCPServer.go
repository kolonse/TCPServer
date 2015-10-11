package TCPServer

import (
	"errors"
	"github.com/kolonse/function"
	"github.com/kolonse/logs"
	"net"
	"reflect"
)

type TCPServer struct {
	Addr      string
	logger    *logs.BeeLogger
	newConnCB *function.Function
	recvCB    *function.Function
	writeCB   *function.Function
}

func (ts *TCPServer) Register(opt ...interface{}) error {
	if len(opt) == 0 {
		return errors.New("参数为空")
	}

	if reflect.TypeOf(opt[0]).Kind() != reflect.String {
		// 如果参数不正确 那么相当于设置失败
		return errors.New("第一个参数必须为字符串 标识要设置的目标")
	}
	switch opt[0] {
	case "newConnCB":
		if len(opt) > 1 {
			ts.newConnCB = opt[1].(*function.Function)
		}
	case "recvCB":
		if len(opt) > 1 {
			ts.recvCB = opt[1].(*function.Function)
		}
	case "writeCB":
		if len(opt) > 1 {
			ts.writeCB = opt[1].(*function.Function)
		}
	case "logger":
		if len(opt) > 1 {
			ts.logger = opt[1].(*logs.BeeLogger)
		}
	}

	return nil
}

// TCP 服务开始接口
func (ts *TCPServer) Server() error {
	// 启动 TCP 服务
	listener, err := net.Listen("tcp", ts.Addr)
	if err != nil {
		return err
	}
	ts.logger.Info("server addr:%s start success!", ts.Addr)
	for {
		//等待客户端接入
		conn, err := listener.Accept()
		if nil != err {
			return err
		}
		tcpConn := NewTCPConn(conn, ts.recvCB, ts.writeCB)
		if ts.newConnCB != nil {
			ts.newConnCB.Call(tcpConn)
		}
		go tcpConn.Server()
	}
}

// 默认的TCP 服务
var DefaultTCPServer *TCPServer

func NewTCPServer(addr string) *TCPServer {
	logger := logs.NewLogger(10000)
	err := logger.SetLogger("console", "")
	if err != nil {
		panic(err.Error())
	}
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(3)

	return &TCPServer{
		Addr:   addr,
		logger: logger,
	}
}

func init() {
	DefaultTCPServer = NewTCPServer("0.0.0.0:9999")
}
