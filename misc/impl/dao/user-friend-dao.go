/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/sirupsen/logrus"
	"strings"
)

const UserRelationTableName = "t_user_relation"

/**
** 添加好友
** return userInfo
**/
func AddFriend(relation *v10.Relation) (ret *v10.FriendUserInfo, err error) {
	var res *v10.Contact
	res, err = AddContactFromUserInternal(relation.Uid, relation.FriId, models.CGUnsorted, relation.IsGroup, relation.Relationship, relation.Remarkname, relation.Tags, relation.Remark)
	if err != nil {
		logrus.Errorf("Err: %v", err)
		return
	}

	ret = &v10.FriendUserInfo{
		UserInfo: res.User,
		Relation: (&models.ContactRelation{}).FromProto(res.Relation).ToOldRelation(res.Relation.Cid+res.Relation.Gid, res.Cb.Uid),
	}

	_, err = AddFriend2(relation)
	return
}

func AddFriend2(relation *v10.Relation) (ret *v10.FriendUserInfo, err error) {
	sql := "insert into " + UserRelationTableName + " (uid,fri_id) values (?,?)"
	rs, err := dbe.DBE.Engine().Exec(sql, relation.Uid, relation.FriId)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}
	row, errr := rs.RowsAffected()
	if errr != nil {
		return nil, errr
	}
	if row > 0 {
		fri, err := GetUserInfoByAutoId(relation.FriId)
		if err != nil {
			return nil, err
		}
		friUserInfo := v10.FriendUserInfo{}
		friUserInfo.UserInfo = fri
		friUserInfo.Relation = relation
		return &friUserInfo, nil
	}
	return nil, nil
}

/**
** 修改朋友信息(增加备注，拉黑名单等)
**/
func UpdateFriend(relation *v10.Relation) (ret *v10.Relation, err error) {
	var res *v10.UpdateContactReply
	var req *v10.UpdateContactReq = &v10.UpdateContactReq{
		ProtoOp: v10.Op_UpdateContact, Seq: 1,
		UidOwner:  relation.Uid,
		UidFriend: relation.FriId,
		Cb:        &v10.ContactBase{
			// TODO contact 表的字段还需要调整：需要精简一部分，转移一部分到其它地方去
			// 当前的联系人表被设计为全局唯一，这是不对的，应该是每个用户有一张联系人表
		},
		Relation: &v10.ContactRelationBase{
			RelationShip: relation.Relationship,
			RemarkName:   relation.Remarkname,
			Remarks:      relation.Remark,
			RemarkTags:   relation.Tags,
		},
		Group: &v10.ContactGroupBase{
			Name: models.CGUnsorted,
		},
	}
	res, err = UpdateContact(req)
	if err != nil {
		logrus.Errorf("Err: %v", err)
		return
	} else {
		logrus.Debugf("updated ok: %v", *res)
		ret = relation
	}

	_, err = UpdateFriend2(relation)
	return
}

func UpdateFriend2(relation *v10.Relation) (*v10.Relation, error) {
	sql := "update " + UserRelationTableName + " set fri_remarkname=?,fri_relationship=?,tags=?,remark=? where uid=? and fri_id=?"
	rs, err := dbe.DBE.Engine().Exec(sql, relation.Remarkname, relation.Relationship, relation.Tags, relation.Remark, relation.Uid, relation.FriId)
	if err != nil {
		return nil, err
	}
	row, errr := rs.RowsAffected()
	if errr != nil {
		return nil, errr
	}
	if row > 0 {
		return relation, nil
	}
	return nil, nil
}

/*
* 删除好友
 */
func DeleteFriend(relation *v10.Relation) bool {
	cgid, cid, _, err := contactIdsBy(relation.Uid, relation.FriId)
	if err != nil {
		logrus.Errorf("Err: %v", err)
		return false
	}

	if res, err := RemoveContact(&v10.RemoveContactReq{ProtoOp: v10.Op_RemoveContact, Seq: 1, UidOwner: relation.Uid, CgId: cgid, CId: cid}); err != nil {
		logrus.Errorf("Err: %v", err)
		return false
	} else {
		logrus.Debugf("delete ok: %v", *res)
	}

	sql := "update " + UserRelationTableName + " set isdel=1 where uid=? and fri_id=?"
	rs, err := dbe.DBE.Engine().Exec(sql, relation.Uid, relation.FriId)
	if err != nil {
		return false
	}
	row, errr := rs.RowsAffected()
	if errr != nil {
		return false
	}
	if row > 0 {
		return true
	}
	return false
}

/*
* 获取好友列表
 */
func GetFriendList(autoId int64, uid string) ([]*v10.FriendUserInfo, error) {
	var (
		rel = &models.ContactRelation{UidOwner: autoId}
		res *v10.ContactGroups
		err error
	)

	res, err = ListContacts(autoId)
	if err != nil {
		return nil, err
	}

	friends := []*v10.FriendUserInfo{}
	for _, g := range res.Groups {
		if strings.EqualFold(g.Cg.Name, models.CGUnsorted) {
			for _, c := range g.Contacts {
				x := &v10.FriendUserInfo{}
				x.Relation = rel.FromProto(c.Relation).ToOldRelation(0, c.Cb.Uid)
				x.UserInfo = &v10.UserInfo{
					Id: c.Cb.Uid, Realname: c.Cb.FullName, Nickname: c.Cb.Nickname,
					Phone: c.Cb.Tel, Avatar: c.Cb.Avatar,
					// TODO 没有传递足够正确的UserInfo，没有在这里发出 dao.get_from_t_user 的数据库查询。
				}
				friends = append(friends, x)
			}
		}
	}

	// var sql = "SELECT " +
	// 	"i.id," +
	// 	"i.unique_id," +
	// 	"i.type," +
	// 	"i.nickname," +
	// 	"i.mobile," +
	// 	"i.avatar," +
	// 	"i.sex," +
	// 	"0," +
	// 	"r.isGroup,r.fri_remarkname,fri_relationship,r.tags,r.remark " +
	// 	"from " +
	// 	"(SELECT * from t_user_relation where isdel=0 and uid=?) r " +
	// 	"LEFT JOIN t_user i " +
	// 	"on r.fri_id = i.id"
	//
	// var value = autoId
	//
	// rows, err := config.MyDb.Query(sql, value)
	//
	// if err != nil {
	// 	return nil, err
	// }
	//
	// friends := []*v10.FriendUserInfo{}
	//
	// for rows.Next() {
	// 	fri_ui := v10.FriendUserInfo{}
	// 	fri := v10.UserInfo{}
	// 	rela := v10.Relation{}
	// 	err := rows.Scan(&fri.Id, &fri.Uid, &fri.Type, &fri.Nickname, &fri.Phone, &fri.Avatar, &fri.Sex, &fri.Age,
	// 		&rela.IsGroup, &rela.Remarkname, &rela.Relationship, &rela.Tags, &rela.Remark)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fri_ui.UserInfo = &fri
	// 	fri_ui.Relation = &rela
	// 	friends = append(friends, &fri_ui)
	// }
	// defer rows.Close()
	return friends, nil
}
