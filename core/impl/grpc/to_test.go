/*
 * Copyright © 2019 Hedzr Yeh.
 */

package grpc_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api/v10"
	"testing"
)

// 20 33
// 8 20 24 33 34

func TestPBMsg2Byte(t *testing.T) {

	var msg = &v10.SendMsgReq{ProtoOp: v10.Op_SendMsg, Seq: 33,
		Body: &v10.SaveMessageRequest{
			GroupId: 1, FromUser: 2, ToUser: 3, MsgContent: "aavv",
			AtUserList: []uint64{5, 6, 7},
			AtAll:      true,
			MsgType:    9,
		},
	}

	b, _ := proto.Marshal(msg)
	t.Logf("b = %v", b)

}
