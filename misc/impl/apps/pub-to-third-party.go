/*
 * Copyright © 2019 Hedzr Yeh.
 */

package apps

type (
	// VxPlug 一个第三方开发的 VxApp，是一个具备 apps.VxPlug 接口实现的实体。参见 demo/demo.go 的示意性实现。
	VxPlug interface {
		OnEvent(p Plugin, args *Args) (err error)
		OnLoad() error
		OnUnload() error
	}
)
