/*
 * Copyright © 2019 Hedzr Yeh.
 */

package exception

import (
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"strconv"
	"strings"
)

const UnknownRequest = "1001::Unknown Request"
const UnknownError = "1002::Unknown Error"
const InvalidParams = "1005::Invalid Parameters"
const InvalidOutputs = "1006::Invalid Outputs"
const DaoError = "1100::Generic DAO error"
const DaoRecordNotExists = "1101::Record not exists"
const DaoCantInsertError = "1102::DAO error: CANNOT insert."
const DaoCantUpdateError = "1103::DAO error: CANNOT update."
const DaoCantRemoveError = "1104::DAO error: CANNOT delete."
const CannotUpdateError = "1102::Cannot update"
const CannotInsertError = "1103::Cannot insert"
const DaoNotFoundError = "1105::DAO error: CANNOT FOUND."

const TokenErr = "2001::invalid token"
const SignErr = "2002::invalid sign"

const UserExist = "3001::user exist"
const UserNotExist = "3002::user not exist"
const UserPasswdMissed = "3008::user passwd missed"
const UserUnknownError = "3009::user unknown error"
const ContactNotExist = "3100::contact not exist"
const ContactGroupNotExist = "3120::contact group not exist"
const ContactRelationNotExist = "3140::contact relation not exist"

type myError struct {
	errCode int64
	errMsg  string
}

func (e *myError) Error() string {
	return "[" + strconv.Itoa(int(e.errCode)) + "][" + e.errMsg + "]"
}

func UnwrapErr(errStr string) *myError {
	errCode, errr := strconv.ParseInt(strings.Split(errStr, "::")[0], 10, 64)
	if errr != nil {
		return nil
	}
	var errMsg string = strings.Split(errStr, "::")[1]
	return &myError{errCode: errCode, errMsg: errMsg}
}

type XmError struct {
	Code       int64
	Msg        string
	InnerError error
}

func (e *XmError) Error() string {
	if e.InnerError == nil {
		return fmt.Sprintf("[%v] %v", e.Code, e.Msg)
	} else {
		return fmt.Sprintf("[%v] %v - cause: %v", e.Code, e.Msg, e.InnerError)
	}
}

func New(msg string) *XmError {
	a := strings.Split(msg, "::")
	errCode, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		return nil
	}
	var errMsg string = a[1]
	return &XmError{errCode, errMsg, nil}
}

func New2(msg string) (*XmError, v10.Err) {
	a := strings.Split(msg, "::")
	errCode, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		return nil, 9999
	}
	var errMsg string = a[1]
	return &XmError{errCode, errMsg, nil}, v10.Err(errCode)
}

func NewWith(msg string, innerErr error) *XmError {
	a := strings.Split(msg, "::")
	errCode, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		errCode = 9999
	}
	var errMsg string = msg
	if len(a) > 1 {
		errMsg = a[1]
	}
	return &XmError{errCode, errMsg, innerErr}
}

func NewError2(msg string, innerErr error) (*XmError, v10.Err) {
	a := strings.Split(msg, "::")
	errCode, err := strconv.ParseInt(a[0], 10, 64)
	if err != nil {
		errCode = 9999
	}
	var errMsg string = msg
	if len(a) > 1 {
		errMsg = a[1]
	}
	return &XmError{errCode, errMsg, innerErr}, v10.Err(errCode)
}
