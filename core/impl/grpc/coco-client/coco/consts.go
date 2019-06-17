/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"github.com/hedzr/voxr-api/api/v10"
)

var (
	DemoLoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "13333333333", Pass: "123456"},
		Device: &v10.Device{SystemOs: "IOS", SystemVersion: "IOS12", Model: "iPhoneXr",
			Nickname: "Tom的手机", Unique: "sasjliuygt6789uiojkhgvfc^&rt5y67"}}

	Demo1LoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "15888888888", Pass: "111111"},
		Device: &v10.Device{SystemOs: "IOS", SystemVersion: "IOS12", Model: "iPhoneXr",
			Nickname: "Alice 的手机", Unique: "466BEC55-AC8E-42BD-B7D6-3691C9B9842F"}}
	Demo2LoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "15888888888", Pass: "111111"},
		Device: &v10.Device{SystemOs: "Windows", SystemVersion: "10", Model: "WINX",
			Nickname: "Alice 的PC", Unique: "05ACCC2A-DB42-49AC-AA44-031D6FAC198B"}}
	Demo3LoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "15888888888", Pass: "111111"},
		Device: &v10.Device{SystemOs: "macOS", SystemVersion: "11", Model: "macOS",
			Nickname: "Alice 的MAC", Unique: "23F480F1-FCA2-4F32-94F8-58E96E0F5F20"}}

	Demo4LoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "17111111111", Pass: "000000"},
		Device: &v10.Device{SystemOs: "IOS", SystemVersion: "IOS10", Model: "iPhone8",
			Nickname: "Bob 的手机", Unique: "BDBD6C43-2D52-4627-B773-D52FB495B465"}}
	Demo5LoginReq = v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "18333333333", Pass: "123456"},
		Device: &v10.Device{SystemOs: "IOS", SystemVersion: "IOS11", Model: "iPhone9",
			Nickname: "Tom 的手机", Unique: "FA8C9D83-FF44-410A-88E0-F427CCE34923"}}

	// Bob	17111111111
	// Tom	18333333333

	// 2	Y78Jko9d7x6S62MI90Ah7D9	0	Alice	15888888888
	// http://img5.duitang.com/uploads/item/201410/02/20141002212239_zWR55.jpeg
	// 09876543456789	2	18	111111	kikiol@asas.com	艾曼	1553157538332	0

	DemoLoginRequests = map[*v10.LoginReq]bool{
		&DemoLoginReq:  false,
		&Demo1LoginReq: false,
		&Demo2LoginReq: false,
		&Demo3LoginReq: false,
		&Demo4LoginReq: false,
		&Demo5LoginReq: false,
	}
)
