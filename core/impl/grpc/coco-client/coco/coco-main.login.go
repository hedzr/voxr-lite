/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func postReq() {
	url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func doLogin() (token, did string, uid uint64, err error) {
	var (
		b          []byte
		req        *http.Request
		pickedUser *v10.LoginReq
	)

	for k, v := range DemoLoginRequests {
		if v == false {
			pickedUser = k
			break
		}
	}

	if pickedUser == nil {
		err = errors.New("slot full, no useable user for new login.")
		return
	}

	// picked.
	DemoLoginRequests[pickedUser] = true

	b, err = json.Marshal(pickedUser)

	url := "http://localhost:7111/v1/api/login"
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	// b, err = base64.StdEncoding.DecodeString(string(body))
	r := &v10.Result{}
	err = json.Unmarshal(body, r)
	if len(r.Data) > 0 {
		tgt := &v10.UserInfoToken{}
		if err := ptypes.UnmarshalAny(r.Data[0], tgt); err != nil {
			fmt.Errorf("ERROR: %v\n", err)
		}

		token = tgt.Token
		did = tgt.DeviceId
		uid = uint64(tgt.UserInfo.Id)

		DemoLoginRequests[pickedUser] = true
	} else {
		logrus.Errorf("LOGIN FAILED: %v", pickedUser)
		DemoLoginRequests[pickedUser] = false // return the user to slots.
	}
	return
}
