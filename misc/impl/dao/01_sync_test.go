/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"github.com/hedzr/voxr-api/models"
	"testing"
)

func TestSyncSchemeToDB(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer localClose()

	err = db.Engine().Sync(
		&models.Contact{}, &models.ContactGroup{}, &models.ContactRelation{}, &models.ContactLinks{},
		&models.UserDevice{})
	if err != nil {
		t.Fatal(err)
	}
}
