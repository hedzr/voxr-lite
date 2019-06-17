/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service_test

import (
	"context"
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
	"time"
)

// A common requirement in programs is getting the number
// of seconds, milliseconds, or nanoseconds since the
// [Unix epoch](http://en.wikipedia.org/wiki/Unix_time).
// Here's how to do it in Go.

func TestForTime(t *testing.T) {

	// Use `time.Now` with `Unix` or `UnixNano` to get
	// elapsed time since the Unix epoch in seconds or
	// nanoseconds, respectively.
	now := time.Now()
	secs := now.Unix()
	nanos := now.UnixNano()
	fmt.Println(now)

	// Note that there is no `UnixMillis`, so to get the
	// milliseconds since epoch you'll need to manually
	// divide from nanoseconds.
	millis := nanos / 1000000
	fmt.Println("secs:", secs, " + nanos:", now.Nanosecond())
	fmt.Println("millis:", millis)
	fmt.Println("nanos:", nanos)

	// You can also convert integer seconds or nanoseconds
	// since the epoch into the corresponding `time`.
	fmt.Println(time.Unix(secs, 0))
	fmt.Println(time.Unix(0, nanos))

	fmt.Println(now)
}

func test1(conn *grpc.ClientConn, in *v10.LoginReq) (token *v10.Result) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	c := v10.NewUserActionClient(conn)
	var err error
	token, err = c.Login(ctx, in)
	if err != nil {
		logrus.Warnf("querying failed: ", err)
	}
	return
}
