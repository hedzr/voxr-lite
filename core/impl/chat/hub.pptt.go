/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

import (
	"fmt"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-lite/core/impl/grpc"
)

func (h *Hub) ppttDoProcessMsg(bin *BinaryMsg) {
	// logrus.Debugf("pptt: (%v,%v) -> %v", bin.from.userId, bin.from.deviceId, bin.body)
	if bin.body[0] == 0xa5 {
		// fn(bin.from, bin.body[1:])

		var body = bin.body[1:]
		var lead = int(body[0])
		// var op = v10.Op(int(bin.body[2])) // 预取一字节以判断 MainOp，注意现在约定 MainOp 编号值不大于 128，因此可以以一个字节来完成检测
		var op int64
		var ate int
		op, ate = api.DecodeZigZagInt(body[1:])

		if lead == 8 { // pb tag = 1
			if err := grpc.Instance.PooledInvoke(v10.Op(op), bin.from, body); err != nil {
				h.writeBack(bin.from, fmt.Sprintf("    [ws] [WARN] Unknown Data Diagram. (op=%v,ate=%v). %v", op, ate, err))
			}
		} else {
			h.writeBack(bin.from, "    [ws] [WARN] Unsupport Data Diagram.")
		}
	}
}

func (h *Hub) writeBack(from *WsClient, msg string) {
	from._writeBack(msg)
	// _ = from.conn.WriteMessage(1, []byte(msg))
}

func (h *Hub) handshake(from *WsClient, body []byte) {
	//
}
