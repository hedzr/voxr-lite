/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/labstack/echo"
)

func Init(callback func(bi *BuildInf)) (ready bool) {
	ready = false

	// e.POST(common.GetApiPrefix()+"/login", h.Login)
	// e.POST(common.GetApiPrefix()+"/refresh-token", h.RefreshToken)

	for _, bi := range aBuildInfoList {
		// e.POST(common.GetApiPrefix()+bi.Entry, buildEchoHandlerFunc(bi))
		callback(bi)
	}

	ready = true
	return
}

var (
	//
	// 定义 RESTful 路由映射表，将 RESTful 请求映射到 grpc 调用
	//
	aBuildInfoList = []*BuildInf{
		//
		// core
		//

		{"/msg/send", api.GrpcCore, api.GrpcCorePackageName, api.CoreActionName,
			"SendMsg", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.SendMsgReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.SendMsgReply)
			}, true, nil},

		//
		// auth, user
		//

		{"/login/v11", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"LoginV11", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AuthReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.AuthReply)
			}, false, func(ret proto.Message) {
				PF().OnLoginV11Ok(ret.(*v10.AuthReply))
			}},
		{"/token/refresh/v11", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"RefreshTokenV11", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AuthReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.AuthReply)
			}, false, func(ret proto.Message) {
				PF().OnRefreshTokenV11Ok(ret.(*v10.AuthReply))
			}},
		{"/register/v11", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"RegisterV11", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AuthReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.AuthReply)
			}, false, func(ret proto.Message) {
				PF().OnRegisterV11Ok(ret.(*v10.AuthReply))
			}},

		// &BuildInf{"/login", api.GrpcAuth, api.GrpcAuthPackageName, UserActionName, LoginMethod, func(c echo.Context) (in interface{}, err error) { in = new(user.LoginReq); err = c.Bind(in); return }},
		// &BuildInf{"/refresh-token", api.GrpcAuth, api.GrpcAuthPackageName, UserActionName, RefreshTokenMethod, func(c echo.Context) (in interface{}, err error) { in = new(user.UserInfoToken); err = c.Bind(in); return }},
		{"/register", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"Register", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AuthReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.UserInfoToken)
			}, true, nil},
		{"/refresh-user-info", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"RefreshUserInfo", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UserInfoToken)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.UserInfoToken)
			}, true, nil},
		{"/update-user-info", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"UpdateUserInfo", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UserInfoToken)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.UserInfoToken)
			}, true, nil},
		// &BuildInf{"/refresh-token", api.GrpcAuth, api.GrpcAuthPackageName, UserActionName,
		// "RefreshToken", func(c echo.Context) (in proto.Message, err error) {
		// 	in = new(user.UserInfoToken)
		// 	err = c.Bind(in)
		// 	return
		// }, nil},
		{"/validate-token", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"ValidateToken", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UserInfoToken)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.TokenValidate)
			}, true, nil},
		{"/validate-token-string", api.GrpcAuth, api.GrpcAuthPackageName, api.UserActionName,
			"ValidateTokenString", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.Token)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.TokenValidate)
			}, true, nil},

		//
		//
		//

		{"/contact/get", api.GrpcUser, api.GrpcUserPackageName, api.UserContactActionName,
			"GetContact", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GetContactReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GetContactReply)
			}, false, nil},
		{"/contact/add", api.GrpcUser, api.GrpcUserPackageName, api.UserContactActionName,
			"AddContact", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AddContactReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.AddContactReply)
			}, false, nil},
		{"/contact/remove", api.GrpcUser, api.GrpcUserPackageName, api.UserContactActionName,
			"RemoveContact", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.RemoveContactReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.RemoveContactReply)
			}, false, nil},
		{"/contact/update", api.GrpcUser, api.GrpcUserPackageName, api.UserContactActionName,
			"UpdateContact", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UpdateContactReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.UpdateContactReply)
			}, false, nil},
		{"/contact/list", api.GrpcUser, api.GrpcUserPackageName, api.UserContactActionName,
			"ListContacts", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.ListContactsReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.ListContactsReply)
			}, false, nil},

		//
		//
		//

		{"/friend/add", api.GrpcUser, api.GrpcUserPackageName, api.FriendActionName,
			"AddFriend", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.Relation)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.FriendUserInfo)
			}, true, nil},

		{"/friend/update", api.GrpcUser, api.GrpcUserPackageName, api.FriendActionName,
			"UpdateFriend", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.Relation)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.Relation) // TODO The return type for UpdateFirend be not sure!
			}, true, nil},

		{"/friend/delete", api.GrpcUser, api.GrpcUserPackageName, api.FriendActionName,
			"DeleteFriend", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.Relation)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				// return new(user.UserInfo)
				return nil
			}, true, nil},

		{"/friend/list", api.GrpcUser, api.GrpcUserPackageName, api.FriendActionName,
			"GetFriendList", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UserId)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.FriendUserInfo)
				// return nil
			}, true, nil},

		//
		// storage
		//

		{"/message/save", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"SaveMessage", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.SaveMessageRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.SaveMessageResponse)
				// return nil
			}, true, nil},

		{"/message/get", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"GetMessage", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GetMessageRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GetMessageResponse)
				// return nil
			}, true, nil},

		{"/message/history", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"GetMessageHistory", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GetMessageHistoryRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GetMessageHistoryResponse)
				// return nil
			}, true, nil},

		{"/message/ack", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"AckMessage", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.AckMessageRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.AckMessageResponse)
				// return nil
			}, true, nil},

		{"/message/update", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"UpdateMessage", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.UpdateMessageRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.UpdateMessageResponse)
				// return nil
			}, true, nil},

		{"/message/unread", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"GetNotReadConversationSet", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GetNotReadConversationSetRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GetNotReadConversationSetResponse)
				// return nil
			}, true, nil},

		//
		//
		//

		{"/conversation/list", api.GrpcStorage, api.GrpcStoragePackageName, api.MessageActionName,
			"GetOffLineConversationSet", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GetOffLineConversationSetRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GetOffLineConversationSetResponse)
				// return nil
			}, true, nil},

		//
		//
		//

		{"/group/create", api.GrpcStorage, api.GrpcStoragePackageName, api.GroupActionName,
			"CreateGroup", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GroupCreateRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GroupCreateResponse)
				// return nil
			}, true, nil},

		{"/group/update", api.GrpcStorage, api.GrpcStoragePackageName, api.GroupActionName,
			"UpdateGroup", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GroupUpdateRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GroupUpdateResponse)
				// return nil
			}, true, nil},

		{"/group/get", api.GrpcStorage, api.GrpcStoragePackageName, api.GroupActionName,
			"GetGroup", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GroupGetRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GroupGetResponse)
				// return nil
			}, true, nil},

		{"/group/remove", api.GrpcStorage, api.GrpcStoragePackageName, api.GroupActionName,
			"RemoveGroup", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.GroupRemoveRequest)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.GroupRemoveResponse)
				// return nil
			}, true, nil},

		//
		// circle
		//

		{"/circle/get", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"GetCircle", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},

		{"/circle/remove", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"RemoveCircle", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},

		{"/circle/update", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"UpdateCircle", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},

		{"/circle/send", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"SendCircle", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},

		{"/circle/list", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"ListCircles", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},

		// 使用 PB 数据包进行图片上传
		{"/circle/upload-image", api.GrpcCircle, api.GrpcCirclePackageName, api.ImCircleActionName,
			"UploadImage", func(c echo.Context) (in proto.Message, err error) {
				in = new(v10.CircleAllReq)
				err = c.Bind(in)
				return
			}, func() (out proto.Message) {
				return new(v10.CircleAllReply)
				// return nil
			}, true, nil},
	}
)
