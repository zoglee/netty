package main

import (
	"net"
)

func SysfdByTcpConn(tcpConn *net.TCPConn) (int, error) {
	file, err := tcpConn.File()
	if err != nil {
		return -1, err
	}
	defer file.Close()
	return int(file.Fd()), nil
}

func exportInnerReqChan(sysfd int, tcpconn *net.TCPConn, inreqc chan *Request) {
	fdmap[sysfd] = &FdContext{
		fd:      sysfd,
		tcpconn: tcpconn,
		inreqc:  inreqc,
		player:  nil,
	}
}
