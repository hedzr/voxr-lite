/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/util"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
)

type FriendServer struct{}

func (friser *FriendServer) AddFriend(ctx context.Context, relation *v10.Relation) (*v10.Result, error) {
	var result = util.MyBaseResult{}.FailResult()
	friendInfo, err := dao.AddFriend(relation)
	if err != nil {
		return result, err
	}
	data, err := ptypes.MarshalAny(friendInfo)
	if err != nil {
		return result, err
	}
	result = util.MyBaseResult{}.SuccessResult([]*any.Any{data})
	return result, nil
}
func (friser *FriendServer) UpdateFriend(ctx context.Context, relation *v10.Relation) (*v10.Result, error) {
	var result = util.MyBaseResult{}.FailResult()
	friendInfo, err := dao.UpdateFriend(relation)
	if err != nil {
		return result, err
	}
	data, err := ptypes.MarshalAny(friendInfo)
	if err != nil {
		return result, err
	}
	result = util.MyBaseResult{}.SuccessResult([]*any.Any{data})
	return result, nil
}
func (friser *FriendServer) DeleteFriend(ctx context.Context, relation *v10.Relation) (*v10.Result, error) {
	var result = util.MyBaseResult{}.FailResult()
	isSuccess := dao.DeleteFriend(relation)
	if isSuccess {
		result = util.MyBaseResult{}.SuccessResult(nil)
	}
	return result, nil
}
func (friser *FriendServer) GetFriendList(ctx context.Context, uids *v10.UserId) (*v10.Result, error) {
	fmt.Printf("uids:%v", uids)
	var result = util.MyBaseResult{}.FailResult()
	friends, err := dao.GetFriendList(uids.Id, uids.Uid)
	fmt.Printf("friends:%v", friends)
	if err != nil {
		return result, err
	}

	any := []*any.Any{}

	for _, fri := range friends {
		data, err := ptypes.MarshalAny(fri)
		if err != nil {
			return result, nil
		}
		any = append(any, data)
	}

	if err != nil {
		return result, err
	}
	result = util.MyBaseResult{}.SuccessResult(any)
	return result, nil
}
