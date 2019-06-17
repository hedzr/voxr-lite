/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dbe

import (
	"github.com/go-xorm/core"
	"github.com/hedzr/voxr-api/models"
	"github.com/sirupsen/logrus"
	"strings"
)

func (db *DB) AddRolesFor(user *models.User, roleNames []string) (err error) {
	for _, name := range roleNames {
		role := new(models.Role)
		ok := false
		ok, err = db.Engine().Where("role_name=?", name).Get(role)
		if ok && err == nil {
			// logrus.Debugf("  role: %v", role)
			userRole := &models.UserRole{}
			ok, err = db.Engine().Where("user_id=? and role_id=?", user.Id, role.Id).Get(userRole)
			if err == nil && !ok {
				userRole = &models.UserRole{UserId: user.Id, RoleId: role.Id}
				_, err = db.Engine().ID(core.PK{1, 2}).Insert(userRole)
			}
		}
	}
	return err
}

func (db *DB) EnsureUserRoles(loginName string, roleNames []string) (err error) {
	err = db.AddRoles(loginName, roleNames)
	return
}

func (db *DB) AddRoles(loginName string, roleNames []string) (err error) {
	err = db.Engine().Where("login_name=?", loginName).Iterate(new(models.User), func(i int, bean interface{}) error {
		user := bean.(*models.User)
		// logrus.Debugf("user: %v", user)
		return db.AddRolesFor(user, roleNames)
	})
	return
}

func (db *DB) EnsureUserRole(loginName, roleName string) (err error) {
	roleNames := strings.Split(roleName, ",")
	return db.AddRoles(loginName, roleNames)
}

func (db *DB) LoadRoles(r *models.User) (roles []*models.Role, err error) {
	// err = db.Engine().Table(&UserRole{}).Alias("m").
	// 	Join("LEFT", []string{"t_user","u"}, "m.user_id=u.id").
	// 	Join("LEFT", []string{"t_role","r"}, "m.role_id=r.id").
	// 	Where("user_id=?", r.Id).
	// 	Find(&urs)

	err = db.Engine().Table(&models.Role{}).Alias("r").
		Join("LEFT", []string{(&models.UserRole{}).TableName(), "m"}, "m.role_id=r.id").
		Join("LEFT", []string{(&models.User{}).TableName(), "u"}, "m.user_id=u.id").
		Where("u.id=?", r.Id).Find(&roles)
	return
}

func (db *DB) UserDrop(loginName string) (err error) {
	var rows int64
	rows, err = db.Engine().Unscoped().Where("login_name=?", loginName).Delete(new(models.User))
	logrus.Debugf("drop %v rows for user.loginName=%v", rows, loginName)
	return
}
