/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"testing"
)

func TestOrgDao_List(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	// o := &dao.OrgDao{}
	// list, err := o.List(10, 0)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("list all rows:")
	// for ix, r := range list {
	// 	t.Logf("%5d: %v", ix, *r)
	// }
	//
	// list, err = o.ListFromParentId(0)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("list all rows:")
	// for ix, r := range list {
	// 	t.Logf("%5d: %v", ix, *r)
	// }
}

func TestOrgDao_AddDel(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	var ret *models.Organization
	var rows int64

	o := &dao.OrgDao{}

	n := &models.Organization{
		Name: "AA", UniqueName: "AA", VirtualName: "AA",
	}
	ret, err = o.AddOrUpdate(n, n.Pid)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)

	if n.Id > 0 {
		t.Logf("new record inserted, id = %v", n.Id)
		rows, err = o.RemoveById(n.Id)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("new record removed, id = %v, rows = %v", n.Id, rows)
	} else {
		t.Fatal("expect the new id returned but it is 0.")
	}
}
