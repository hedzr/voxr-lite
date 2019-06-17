/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/models"
	"testing"
)

func TestRoles(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	// ref1(t, db)

	user := &models.User{LoginName: "jy"}
	var users []*models.User
	if err := db.Engine().Find(&users, user); err != nil {
		t.Fatalf("Err: %v", err)
	} else {
		if len(users) == 0 {
			t.Logf("user not found")
		} else {
			t.Logf("users found: %v", users[0])
		}
	}

	role := &models.Role{Id: 15}
	roleTo := &models.Role{Id: 15, RoleName: "r5", RoleDesc: "r5 role"}
	var roles []*models.Role
	if err := db.Engine().Find(&roles, role); err != nil {
		t.Fatalf("Err: %v", err)
	} else {
		if len(roles) == 0 {
			if rows, err := db.Engine().Insert(role); err != nil {
				t.Fatalf("Err: %v", err)
			} else {
				t.Logf("%v row(s) inserted.", rows)
			}
		} else {
			t.Logf("roles found: %v", roles[0])
			if rows, err := db.Engine().Update(roleTo, role); rows != 1 || err != nil {
				t.Fatalf("Err: %v | rows = %v", err, rows)
			} else {
				t.Logf("%v row(s) updated.", rows)
			}
		}
	}
}
