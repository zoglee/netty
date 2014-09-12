package main

import (
	"fmt"

	"code.google.com/p/go.net/context"
	"github.com/golang/glog"
)

func Work(ctx context.Context, errc chan<- error, sysfd int, reqc <-chan *Request, inreqc <-chan *Request, rspc chan<- *Response) {
	go func() {
		defer func() {
			glog.Infof("Work is stoped")
		}()

		for {
			var rsp *Response
			var err error
			select {
			case <-ctx.Done():
				glog.Infof("Done for Work")
				return

			case req := <-reqc:
				if rsp, err = work(req); err != nil {
					errc <- fmt.Errorf("in Work req channel: %s", err)
					return
				}
				rspc <- rsp

			case req := <-inreqc:
				var rsp *Response
				if rsp, err = work(req); err != nil {
					errc <- fmt.Errorf("in Work inner req channel: %s", err)
					return
				}
				rspc <- rsp
			}
		}
	}()
}

func work(req *Request) (rsp *Response, err error) {
	return &Response{
		Head: req.Head,
		Body: req.Body,
	}, nil
}
