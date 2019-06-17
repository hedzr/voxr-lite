/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"testing"
)

func TestContacts(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	err = db.Engine().Sync(
		&models.Contact{}, &models.ContactGroup{}, &models.ContactRelation{}, &models.ContactLinks{})
	if err != nil {
		t.Fatal(err)
	}

	var users []*models.User
	if err := db.Engine().Find(&users); err != nil {
		t.Fatalf("Err: %v", err)
	} else {
		if len(users) == 0 {
			t.Logf("user not found")
		} else {
			for ix, user := range users {
				t.Logf("users found: %3d, %v", ix, user)
				ensureUserRoles(t, db, user)
				ensureUserContacts(t, db, user)
			}
		}
	}

	user := &models.User{Nickname: "Dave"}
	if ok, err := db.Engine().Get(user); err != nil || ok == false {
		t.Fatal(err)
	}

}

func TestContactsLoop(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

}

// 测试 dao.AddContactFromUser
// 初始化联系人关联关系表组记录集
func TestAddContact(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	t.Log("Default group")
	for ix, rr := range [][]int64{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 9},
		{1, 10},
		{1, 11},
		{1, 12},

		{2, 1},
		{2, 2},
		{2, 3},
		{2, 4},
		{2, 9},
		{2, 10},
		{2, 11},

		{3, 1},
		{3, 2},
		{3, 12},
		{3, 13},
		{3, 14},
	} {
		if res, err := dao.AddContactFromUser(rr[0], rr[1], models.CGUnsorted); err != nil {
			t.Fatalf(" #%5d err: %v", ix, err)
		} else {
			t.Logf(" #%5d result: %v", ix, res)
		}
	}

	t.Log("Home group")
	for ix, rr := range [][]int64{
		{1, 15},
		{1, 16},
		{1, 17},
	} {
		if res, err := dao.AddContactFromUser(rr[0], rr[1], models.CGHome); err != nil {
			t.Fatalf(" #%5d err: %v", ix, err)
		} else {
			t.Logf(" #%5d result: %v", ix, res)
		}
	}

	t.Log("Office group")
	for ix, rr := range [][]int64{
		{1, 18},
		{1, 19},
	} {
		if res, err := dao.AddContactFromUser(rr[0], rr[1], models.CGOffice); err != nil {
			t.Fatalf(" #%5d err: %v", ix, err)
		} else {
			t.Logf(" #%5d result: %v", ix, res)
		}
	}

	t.Log("Mates group")
	for ix, rr := range [][]int64{
		{1, 23},
	} {
		if res, err := dao.AddContactFromUser(rr[0], rr[1], models.CGMates); err != nil {
			t.Fatalf(" #%5d err: %v", ix, err)
		} else {
			t.Logf(" #%5d result: %v", ix, res)
		}
	}

}

func TestListContacts(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	if res, err := dao.ListContacts(1); err != nil {
		t.Fatalf("err: %v", err)
	} else {
		t.Logf("result for user %v:", res.UidOwner)
		for ix, g := range res.Groups {
			t.Logf("%2d. group '%v' // id=%v", ix, g.Cg.Name, g.Cg.Id)
			for j, cc := range g.Contacts {
				t.Logf("  %3d. %v", j, *cc.Cb)
				t.Logf("     . %v", *cc.Relation)
			}
		}
	}
}

func TestUpdateContact(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	req := &v10.GetContactReq{ProtoOp: v10.Op_GetContact, Seq: 23, UidOwner: 1, UidFriend: 3}
	if res, err := dao.GetContact(req); err != nil {
		t.Fatalf("err: %v", err)
	} else {
		t.Logf("result for user %v:", req.UidOwner)
		t.Logf("    relation: %v:", *res.Relation)
		t.Logf("       group: %v:", *res.Group)
		t.Logf("     contact: %v:", *res.Cb)

		req1 := &v10.UpdateContactReq{ProtoOp: v10.Op_UpdateContact, Seq: 24, UidOwner: 1, UidFriend: res.Cb.Uid, Cb: res.Cb, Group: res.Group, Relation: res.Relation}
		if res1, err := dao.UpdateContact(req1); err != nil {
			t.Fatalf("err: %v", err)
		} else {
			t.Logf("update ok=%v", res1.Ok)

			req := &v10.GetContactReq{ProtoOp: v10.Op_GetContact, Seq: 23, UidOwner: 1, UidFriend: 3}
			if res, err := dao.GetContact(req); err != nil {
				t.Fatalf("err: %v", err)
			} else {
				t.Logf("result for user %v:", req.UidOwner)
				t.Logf("    relation: %v:", *res.Relation)
				t.Logf("       group: %v:", *res.Group)
				t.Logf("     contact: %v:", *res.Cb)
			}
		}
	}
}

func TestUpdateFriend(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()
	prepareOldLocalDb()
	defer localClose()

	// 1. Add Friend
	rel := &v10.Relation{Uid: 1, FriId: 6, Remarkname: "AAA", Relationship: 1, Tags: "BBB", Remark: "CCC"}
	if ret := dao.DeleteFriend(rel); ret {
		t.Logf("ret: %v", ret)
	}
	if ret, err := dao.AddFriend(rel); err != nil {
		t.Fatalf("err: %v", err)
	} else {
		t.Logf("ret: %v", ret)
	}

	// 2. Update
	rel.Tags = "BBB,CCC"
	if ret, err := dao.UpdateFriend(rel); err != nil {
		t.Fatalf("err: %v", err)
	} else {
		t.Logf("ret: %v", ret)
	}

	// 3. delete
	if ret := dao.DeleteFriend(rel); ret {
		t.Logf("ret: %v", ret)
	}

	// 4. List
}
