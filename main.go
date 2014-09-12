package main

import (
	"flag"
	"net"
	"runtime"

	"code.google.com/p/go.net/context"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	laddr, _ := net.ResolveTCPAddr("tcp", ":"+port) // TODO(zog): flag
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		glog.Fatalf("net.ListenTCP: %v", err)
	}

	ctx := context.Background()

	glog.Infof("server start listen on: %s", laddr.String())
	for {
		tcpconn, err := l.AcceptTCP()
		if err != nil {
			glog.Errorf("Accept error, server stop (listen on:%s), err: %s",
				laddr.String(), err)
			break
		}

		go HandleConnection(ctx, tcpconn)
	}
	glog.Infof("server stop listen on: %s)", laddr.String())
}

func HandleConnection(ctx context.Context, tcpconn *net.TCPConn) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		tcpconn.Close()
		if err := recover(); err != nil {
			glog.Errorf("tcpconn panic: %s", err)
		}
		cancel()
	}()

	sysfd, err := SysfdByTcpConn(tcpconn)
	if err != nil {
		glog.Errorf("SysfdByTcpConn err: %s, addr: %s, fd=%d",
			err, tcpconn.RemoteAddr().String(), sysfd)
		return
	}

	glog.Infof("Start handle conn, addr: %s, fd=%d", tcpconn.RemoteAddr().String(), sysfd)

	errc := make(chan error)
	inreqc := make(chan *Request, 10) // TODO(zog): is 10 a good number?
	reqc := Recv(ctx, errc, sysfd, tcpconn)
	rspc := Send(ctx, errc, sysfd, tcpconn)
	Work(ctx, errc, sysfd, reqc, inreqc, rspc)
	Update(ctx, errc, sysfd, rspc)

	exportInnerReqChan(sysfd, tcpconn, inreqc)

	if err := <-errc; err != nil {
		glog.Errorf("HandleConnection err: %v, addr: %s", err, tcpconn.RemoteAddr().String())
	}

	glog.Infof("Stop handle conn, addr: %s, fd=%d", tcpconn.RemoteAddr().String(), sysfd)
}
