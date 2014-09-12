package main

import (
	"fmt"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/golang/glog"
)

func Update(ctx context.Context, errc chan<- error, sysfd int, rspc chan<- *Response) {
	go func() {
		defer func() {
			glog.Infof("Update is stoped")
		}()

		for {
			select {
			case <-ctx.Done():
				glog.Infof("Done for Update")
				return

			default:
				if err := update(sysfd, rspc); err != nil {
					errc <- fmt.Errorf("in Update: %s", err)
					return
				}
			}
		}
	}()
}

func update(sysfd int, rspc chan<- *Response) (err error) {
	time.Sleep(1 * time.Second)
	glog.Infof("Update, time:%v", time.Now())
	return nil
}
