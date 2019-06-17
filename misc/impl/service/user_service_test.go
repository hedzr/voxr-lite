/*
 * Copyright © 2019 Hedzr Yeh.
 */

package service_test

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/xs/mjwt"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/misc/impl/service"
	"testing"
	"time"
)

func prepareDb() (db *dbe.DB, err error) {
	// db = dbe.New()
	// err = db.OpenUrl("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	err = dbe.OpenDbConnection()
	if err != nil {
		err = dbe.DBE.OpenUrl("mysql", "dev:123456@tcp(localhost:3306)/db_im?charset=utf8&parseTime=true")
	}
	if err == nil {
		db = dbe.DBE
	}
	return
}

func TestAuthServer_RefreshToken(t *testing.T) {
	db, err := prepareDb()
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer db.Close()

	redis_op.JwtInit()

	as := &service.AuthServer{}
	ctx := context.TODO()

	var DemoLoginReq = &v10.LoginReq{UserInfo: &v10.UserInfo{Phone: "13333333333", Pass: "123456"},
		Device: &v10.Device{SystemOs: "IOS", SystemVersion: "IOS12", Model: "iPhoneXr",
			Nickname: "Tom的手机", Unique: "sasjliuygt6789uiojkhgvfc^&rt5y67"}}

	var res *v10.Result
	var token *jwt.Token
	res, err = as.Login(ctx, DemoLoginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Logf("res: %v", res)
	if res.Ok && len(res.Data) > 0 {
		uit := new(v10.UserInfoToken)
		err = ptypes.UnmarshalAny(res.Data[0], uit)
		if err != nil {
			t.Fatalf("UnmarshalAny failed: %v", err)
		}
		t.Logf("uit: %v", uit)
		t.Logf("token: %v", uit.Token)

		token, err = redis_op.JwtExtractToken(uit.Token)
		c := token.Claims.(*mjwt.ImClaims)
		tm := api.Int64SecondsToTime(c.ExpiresAt)
		t.Logf("token claims: %v, %v", tm, c)

		time.Sleep(5 * time.Second)

		res, err = as.RefreshToken(ctx, uit)
		if err != nil {
			t.Fatalf("RefreshToken failed: %v", err)
		}
		t.Logf("res: %v", res)
		if res.Ok && len(res.Data) > 0 {
			uit := new(v10.UserInfoToken)
			err = ptypes.UnmarshalAny(res.Data[0], uit)
			if err != nil {
				t.Fatalf("UnmarshalAny failed: %v", err)
			}
			t.Logf("uit: %v", uit)
			t.Logf("token: %v", uit.Token)

			token, err = redis_op.JwtExtractToken(uit.Token)
			c := token.Claims.(*mjwt.ImClaims)
			tm := api.Int64SecondsToTime(c.ExpiresAt)
			t.Logf("token claims: %v, %v", tm, c)
		}
	}
}
