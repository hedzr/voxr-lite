/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao_test

import (
	"testing"
)

func TestDbOpen(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Errorf("open failed: %v, %v", err, db)
	}
	defer localClose()

	// Get User Testing

}
