/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/sirupsen/logrus"
)

type UserDao struct{}

func (s *UserDao) Add(in *models.User) (id uint64, err error) {
	var (
		ok   bool
		rows int64
		rmap []map[string]interface{}
	)

	if len(in.Nickname) == 0 {
		if len(in.LoginName) != 0 {
			in.Nickname = in.LoginName
		} else if len(in.Mobile) != 0 {
			in.Nickname = in.Mobile
		} else if len(in.Email) != 0 {
			in.Nickname = in.Email
		} else {
			err = exception.New(exception.InvalidParams)
			return
		}
	}

	tmpl := &models.User{Nickname: in.Nickname}
	if ok, err = dbe.DBE.Engine().Exist(tmpl); err != nil {
		return
	}

	if !ok {
		if rows, err = dbe.DBE.Engine().Insert(in); err != nil {
			return
		}
		if rows != 1 {
			err = exception.New(exception.DaoCantInsertError)
			return
		}
		id = uint64(in.Id)
	} else {
		if rmap, err = dbe.DBE.Engine().QueryInterface("SELECT id FROM t_user WHERE nickname = ?", in.Nickname); err != nil {
			return
		}
		if rmap == nil {
			err = exception.New(exception.DaoError)
			return
		}
		id = uint64(rmap[0]["id"].(int64))
	}
	return
}

func (s *UserDao) RemoveById(id uint64) (rows int64, err error) {
	tmpl := &models.User{Id: int64(id)}
	if rows, err = dbe.DBE.Engine().Delete(tmpl); err != nil {
		return
	}
	return
}

func (s *UserDao) Remove(tmpl *models.User) (rows int64, err error) {
	if tmpl.Id == 0 && (len(tmpl.LoginName) == 0 || len(tmpl.Nickname) == 0 || len(tmpl.Mobile) == 0 || len(tmpl.Email) == 0) {
		err = exception.New(exception.InvalidParams)
		return
	}

	if rows, err = dbe.DBE.Engine().Delete(tmpl); err != nil {
		return
	}
	return
}

func (s *UserDao) Update(in *models.User) (ret *models.User, rows int64, err error) {
	if in.Id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}

	if rows, err = dbe.DBE.Engine().Update(in); err != nil {
		return
	}
	ret = in
	return
}

func (s *UserDao) GetById(id uint64) (ret *models.User, err error) {
	if id == 0 {
		err = exception.New(exception.InvalidParams)
		return
	}

	tmpl := &models.User{Id: int64(id)}
	if ret, err = s.Get(tmpl); err != nil {
		return
	}
	return
}

func (s *UserDao) Get(tmpl *models.User) (ret *models.User, err error) {
	var ok bool
	if ok, err = dbe.DBE.Engine().Get(tmpl); err != nil {
		return
	}
	if !ok {
		err = exception.New(exception.DaoNotFoundError)
	}
	ret = tmpl
	return
}

func (s *UserDao) List(tmpl *models.User) (ret []*models.User, err error) {
	if err = dbe.DBE.Engine().Find(&ret, tmpl); err != nil {
		return
	}
	return
}

//
//

/**
** 用户登录(密码方式/Mobile/Phone)
**/
func UserLoginUsePass(userInfo *v10.UserInfo) (ret *v10.UserInfo, err error) {
	if len(userInfo.Pass) == 0 {
		return
	}

	tmp := &models.User{
		LoginName: userInfo.Nickname,
		Mobile:    userInfo.Phone,
		Email:     userInfo.Email,
		Type:      0,
	}
	yes, err := dbe.DBE.Engine().Where("type&1=? and status&1=0", 0).Get(tmp)
	if yes && err == nil && tmp.IsPasswordMatched(userInfo.Pass) {
		return tmp.ToProto(), err
	}

	err = exception.New(exception.UserPasswdMissed)
	// if len(userInfo.Phone) > 0 { // 手机登录，如果没有，就注册，优先手机登录
	// 	u, _ := GetUserInfoByPhoneAndPass(userInfo.Phone, userInfo.Pass)
	// 	// if u == nil {
	// 	// 	uReg,_ := UserRegister(userInfo)
	// 	// 	return uReg
	// 	// }else {
	// 	// 	return u
	// 	// }
	// 	return u
	// }
	// if len(userInfo.Email) > 0 {
	// 	u, _ := GetUserInfoByEmail(userInfo.Email, userInfo.Pass)
	// 	return u
	// 	// if u == nil {
	// 	// 	uReg,_ := UserRegister(userInfo)
	// 	// 	return uReg
	// 	// }else {
	// 	// 	return u
	// 	// }
	// }
	return
}

/**
** 验证码验证通过后采用这个登录
**/
func UserLoginByPhoneCode(phone string) *v10.UserInfo {
	u, _ := GetUserInfoByPhone(phone)
	if u == nil {
		userInfo := &v10.UserInfo{}
		userInfo.Phone = phone
		uReg, _ := UserRegister(userInfo)
		return uReg
	} else {
		return u
	}
}

/**
** 获取朋友列表
**/
func UserGetFriends() {

}

func getUserBy(templ *models.User) (*v10.UserInfo, error) {
	yes, err := dbe.DBE.Engine().Where("type&1=? and status&1=0", 0).Get(templ)
	if yes {
		return templ.ToProto(), err
	}
	return nil, err
}

/**
** 获取某个用户详情
**/
func GetUserInfoByUid(uid string) (*v10.UserInfo, error) {
	tmp := &models.User{UniqueId: uid}
	return getUserBy(tmp)
	// return queryUser("uid", uid, "")
}

// 通过自增长的id
func GetUserInfoByAutoId(autoId int64) (*v10.UserInfo, error) {
	tmp := &models.User{Id: autoId}
	return getUserBy(tmp)
	// return queryUser("id", strconv.FormatInt(autoId, 10), "")
}

// 验证码验证通过后直接调用该方法登录
func GetUserInfoByPhone(phone string) (*v10.UserInfo, error) {
	tmp := &models.User{Mobile: phone}
	return getUserBy(tmp)
	// return queryUser("u_phone", phone, "")
}

// // 电话加密码的方式
// func GetUserInfoByPhoneAndPass(phone string, pass string) (*v10.UserInfo, error) {
// 	tmp := &models.User{Mobile: phone,}
// 	ret, err := getUserBy(tmp)
// 	if err == nil {
// 		if tmp.IsPasswordMatched(pass) {
// 			return ret, err
// 		} else {
// 			return nil, exception.UnwrapErr(exception.UserPasswdMissed)
// 		}
// 	}
// 	return nil, err
// 	// return queryUser("u_phone", phone, pass)
// }
//
// // 邮箱加密码的方式
// func GetUserInfoByEmail(email string, pass string) (*v10.UserInfo, error) {
// 	tmp := &models.User{Email: email,}
// 	ret, err := getUserBy(tmp)
// 	if err == nil {
// 		if tmp.IsPasswordMatched(pass) {
// 			return ret, err
// 		} else {
// 			return nil, exception.UnwrapErr(exception.UserPasswdMissed)
// 		}
// 	}
// 	return nil, err
// 	// return queryUser("u_emial", email, pass)
// }

/**
** 新用户注册
**/
func UserRegister(userInfo *v10.UserInfo) (*v10.UserInfo, error) {
	var (
		tmp = &models.User{
			LoginName: userInfo.Nickname,
			Mobile:    userInfo.Phone,
			Email:     userInfo.Email,
		}
		yes, err = dbe.DBE.Engine().Exist(tmp)
		// tmp := new(models.User)
		// yes, err = dbe.DBE.Engine().Where("login_name=? OR nickname=? OR mobile=? OR (email=? AND email is not null)", obj.LoginName, obj.Nickname, obj.Mobile, obj.Email).Exist(tmp)
	)
	if yes {
		return nil, exception.UnwrapErr(exception.UserExist)
	}
	if err != nil {
		return nil, exception.UnwrapErr(err.Error())
	}

	obj := &models.User{
		Password:  userInfo.Pass,
		Nickname:  userInfo.Nickname,
		LoginName: userInfo.Nickname,
		Mobile:    userInfo.Phone,
		Email:     userInfo.Email,
		Sex:       int16(userInfo.Sex),
		GivenSn:   userInfo.Idcard,
		FullName:  userInfo.Realname,
		Type:      int(userInfo.Type),
		Avatar:    userInfo.Avatar,
		Status:    models.StatusStandardValid,
		// State, CreateTime
	}

	var rows int64
	rows, err = dbe.DBE.Engine().Insert(obj)
	if err == nil && rows == 1 {
		logrus.Debugf("obj = %v", obj)
		userInfo.Id = int64(obj.Id) // uint64
		userInfo.Uid = obj.UniqueId
		userInfo.State = int32(obj.Status)
		userInfo.Pass = ""

		err = dbe.DBE.AddRolesFor(obj, []string{models.RoleUser, models.RoleImUser})
		if err != nil {
			logrus.Error(err)
		}

		// send new user registered message

		// add main group for new user
		cg := &models.ContactGroup{UidOwner: obj.Id, Name: models.CGUnsorted}
		rows, err = dbe.DBE.Engine().Insert(cg)
		if err != nil {
			logrus.Error(err)
		}

		return userInfo, err
	} else {
		logrus.Error(err)
	}

	// uid := util.UUid()
	// userInfo.Uid = uid
	// userInfo.Createtime = int64(time.Now().UnixNano() / 1000000)
	// createTime := strconv.FormatInt(userInfo.Createtime, 10)
	// sql := "insert into " + UserTableName + " (uid,u_phone,createtime) values (?,?,?)"
	// logrus.Info("UserRegister==>sql:" + sql)
	// db := config.MyDb
	// rs, err := db.Exec(sql, uid, userInfo.Phone, createTime)
	// if err != nil {
	// 	logrus.Errorln(err)
	// 	return nil, err
	// }
	// if row, _ := rs.RowsAffected(); row > 0 {
	// 	id, _ := rs.LastInsertId()
	// 	userInfo.Id = id
	// 	return userInfo, nil
	// } else {
	// 	return nil, excep.UnwrapErr(excep.UserExist)
	// }

	return nil, err // exception.UnwrapErr(exception.UserUnknownError)
}

/**
** 修改用户信息
**/
func UpdateUserInfo(userInfo *v10.UserInfo) (ok bool, err error, obj *models.User) {
	tmpl := &models.User{
		LoginName: userInfo.Nickname,
		Mobile:    userInfo.Phone,
		Email:     userInfo.Email,
	}

	// obj := new(models.User)
	ok, err = dbe.DBE.Engine().Exist(tmpl) // Where("login_name=? OR nickname=? OR mobile=? OR (email=? AND email is not null)", tmp.LoginName, tmp.Nickname, tmp.Mobile, tmp.Email).Exist(obj)
	if ok {
		obj = &models.User{
			Password: userInfo.Pass,
			Nickname: userInfo.Nickname,
			Sex:      int16(userInfo.Sex),
			GivenSn:  userInfo.Idcard,
			FullName: userInfo.Realname,
			Type:     int(userInfo.Type),
			Avatar:   userInfo.Avatar,
			// Status:    models.StatusStandardValid,
			// Status:   1,
			// State, CreateTime
		}

		if len(obj.Password) > 0 {
			obj.NewEncodePwd()
		}

		// // obj.Mobile, obj.Email, obj.LoginName: CAN'T BE MODIFIED HERE
		var rows int64
		rows, err = dbe.DBE.Engine().Update(obj, tmpl)
		if rows != 1 || err != nil {
			logrus.Warnf("CAN'T update user info: %v | %v", err, rows)
			ok = false
		} else {
			ok = true
		}
	}
	if err != nil {
		logrus.Warnf("CAN'T find user info before updating it: %v", err)
	}

	// sql := "update " + UserTableName + " set u_type=?,u_nickname=?,u_phone=?,u_avatar=?,u_idcard=?,u_sex=?," +
	// 	"u_age=?,u_pass=?,u_email=?,u_realname=? where uid=?"
	// res, err := config.MyDb.Exec(sql, userInfo.Type, userInfo.Nickname, userInfo.Phone,
	// 	userInfo.Avatar, userInfo.Idcard, userInfo.Sex, userInfo.Age, userInfo.Pass,
	// 	userInfo.Email, userInfo.Realname, userInfo.Uid)
	// if err != nil {
	// 	logrus.Errorln(err)
	// 	return false
	// }
	// row, err := res.RowsAffected()
	// if err != nil {
	// 	logrus.Errorln(err)
	// 	return false
	// }
	// return row > 0
	return
}

// /**
// ** 查询用户
// **/
// func queryUser(where string, whereValue string, pass string) (*v10.UserInfo, error) {
// 	// var sql = "select id,uid,u_type,u_nickname,u_phone,u_avatar,u_idcard,u_sex,u_age,u_pass,u_email,u_realname,createtime" +
// 	// 	" from " + UserTableName + " where " + where + "=? "
// 	// var row *sql2.Row
// 	// if len(pass) > 0 {
// 	// 	sql += " and u_pass=?"
// 	// 	row = config.MyDb.QueryRow(sql, whereValue, pass)
// 	// } else {
// 	// 	row = config.MyDb.QueryRow(sql, whereValue)
// 	// }
// 	//
// 	// userInfoRes := v10.UserInfo{}
// 	//
// 	// err := row.Scan(&userInfoRes.Id,
// 	// 	&userInfoRes.Uid,
// 	// 	&userInfoRes.Type,
// 	// 	&userInfoRes.Nickname,
// 	// 	&userInfoRes.Phone,
// 	// 	&userInfoRes.Avatar,
// 	// 	&userInfoRes.Idcard,
// 	// 	&userInfoRes.Sex,
// 	// 	&userInfoRes.Age,
// 	// 	&userInfoRes.Pass,
// 	// 	&userInfoRes.Email,
// 	// 	&userInfoRes.Realname,
// 	// 	&userInfoRes.Createtime)
// 	// if err != nil {
// 	// 	logrus.Errorln(err)
// 	// 	return nil, err
// 	// }
// 	//
// 	// return &userInfoRes, nil
//
// 	return nil, nil
// }
