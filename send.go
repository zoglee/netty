package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/golang/glog"
)

func Send(ctx context.Context, errc chan<- error, sysfd int, tcpconn *net.TCPConn) (rspc chan *Response) {
	rspc = make(chan *Response, 10) // TODO(zog): is 10 a good number?

	go func() {
		defer func() {
			glog.Infof("Send is stoped")
		}()

		for {
			select {
			case <-ctx.Done():
				glog.Infof("Done for Send")
				return

			case rsp := <-rspc:
				if err := send(tcpconn, rsp); err != nil {
					errc <- fmt.Errorf("in Send: %s", err)
					return
				}
			}
		}
	}()

	return rspc
}

func send(tcpconn *net.TCPConn, rsp *Response) (err error) {
	Len := uint32(PkgLenSize) + uint32(len(rsp.Head)) + uint32(len(rsp.Body))
	Hlen := uint16(Uint16Size) + uint16(len(rsp.Head))
	data := make([]byte, 0, int(Len)) // len:0, cap:Len; TODO(zog): cache
	buf := bytes.NewBuffer(data)      // TODO(zog): 复用
	binary.Write(buf, binary.BigEndian, Len)
	binary.Write(buf, binary.BigEndian, Hlen)
	buf.Write(rsp.Head)
	buf.Write(rsp.Body)
	if debug {
		glog.Infof("sent bytes to %s, len: %d",
			tcpconn.RemoteAddr().String(), len(buf.Bytes()))
		glog.Flush()
	}

	tcpconn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	if _, err = tcpconn.Write(buf.Bytes()); err != nil {
		return err
	}

	if debug {
		glog.Infof("sent data(len:%d): %v", buf.Len(), buf.Bytes())
		glog.Flush()
	}

	return nil
}
