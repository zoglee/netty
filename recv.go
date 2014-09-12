package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/golang/glog"
)

func Recv(ctx context.Context, errc chan<- error, sysfd int, tcpconn *net.TCPConn) (reqc chan *Request) {
	reqc = make(chan *Request, 10) // TODO(zog): is 10 a good number?

	go func() {
		defer func() {
			glog.Infof("Recv is stoped")
		}()

		var req *Request
		var err error
		for {
			select {
			case <-ctx.Done():
				glog.Infof("Done for Recv")
				return

			default:
				if req, err = recv(tcpconn); err != nil {
					errc <- fmt.Errorf("in Recv: %s", err)
					return
				}
				reqc <- req
			}
		}
	}()

	return reqc
}

func recv(tcpconn *net.TCPConn) (req *Request, err error) {
	// TODO(zog): 用 bytes.Buffer 优化
	var Len uint32
	var Hlen uint16

	tcpconn.SetDeadline(time.Now().Add(100 * time.Millisecond)) // TODO(zog): is 100ms good?
	if err = binary.Read(tcpconn, binary.BigEndian, &Len); err != nil {
		return nil, err
	}
	if int(Len) < MinAllowedPkgLen || int(Len) > MaxAllowedPkgLen {
		err = fmt.Errorf("invalid package len: Len=%d(min:%d, max:%d)",
			Len, MinAllowedPkgLen, MaxAllowedPkgLen)
		return nil, err
	}

	data := make([]byte, int(Len)-Uint32Size) // TODO(zog): cache

	tcpconn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	if n, rerr := io.ReadFull(tcpconn, data); err != nil {
		err = fmt.Errorf("read package data n:%d(expect:%d), err: %s",
			n, len(data), rerr)
		return nil, err
	}

	Hlen = (uint16(data[1]) | uint16(data[0])<<8)
	if (int(Hlen) < MinAllowedPkgTotalHeadLen) ||
		(int(Hlen) > MaxAllowedPkgTotalHeadLen) ||
		(int(Hlen) > len(data)-MinAllowedPkgBodyLen) {
		err = fmt.Errorf("invalid Hlen: %d, (min:%d, max:%d), over expect max: %d",
			Hlen, MinAllowedPkgTotalHeadLen, MaxAllowedPkgTotalHeadLen, len(data))
		return nil, err
	}

	if debug {
		glog.Infof("Len:%d, Hlen:%d, len(data):%d, data:%v",
			Len, Hlen, len(data), data)
		glog.Flush()
	}

	return &Request{
		Head: data[2 : 2+(Hlen-2)],
		Body: data[2+(Hlen-2):],
	}, nil
}
