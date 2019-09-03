package cnst

import (
	"errors"
	"strings"
)

const SECRET = "thanx for fish and good bye"

const STATIC = "./client/public/"

const (
	All = iota
	USER
	PASS
	PRDT
	ORDR
)

const (
	UPD = iota
	INS
	DEL
	NTUPD
	NTINS
	NTDEL
	EXIST
	NTFND
	EMPTY
	DENY
)

const (
	UNKN = iota
	CLNT
	SOLD
	MNGR
	EMPL
	ADMN
)

var maps = [][]string{
	{
		"Record successfully updated",
		"User successfully updated",
		"Password successfully updated",
		"Product has been successfully updated",
		"Order successfully updated",
	},
	{
		"Successfully added record",
		"User successfully added",
		"Password successfully added",
		"Product successfully added",
		"Order successfully added",
	},
	{
		"Recording deleted successfully",
		"User deleted successfully",
		"Password successfully deleted",
		"Product successfully deleted",
		"Order successfully deleted",
	},
	{
		"The entry was not updated.",
		"User has not been updated",
		"Password was not updated",
		"Product was not updated",
		"Order was not updated",
	},
	{
		"The entry was not added",
		"User was not added",
		"Password or login is incorrect",
		"Product was not added",
		"Order was not added",
	},
	{
		"The entry was not deleted.",
		"User has not been deleted",
		"Password was not deleted",
		"Product was not updated",
		"Order was not updated",
	},
	{
		"Entry already exists",
		"Login busy, try another one",
		"Error creating password",
		"Product already exists",
		"Order already exists",
	},
	{
		"This entry does not exist",
		"User does not exist or is blocked",
		"User does not exist or is blocked",
		"Product does not exist",
		"Order does not exist",
	},
	{
		"Fill in all required fields",
		"Fill in all required fields",
		"Fill in all required fields",
		"Fill in all required fields",
		"Fill in all required fields",
	},
	{
		"Not have permission for this action",
		"Not have permission for this action",
		"Not have permission for this action",
		"Not have permission for this action",
		"Not have permission for this action",
	},
}

const INCLOGIN = "Пароль или логин неверен"

func Message(status, model byte) string {
	return maps[status][model]
}

func Error(status, model byte) error {
	return errors.New(maps[status][model])
}

func Status(err error, status, model byte) (int, []byte) {
	code := 200
	arr := []string{Message(status, model)}
	if err != nil {
		code = 400
		arr = append(arr, err.Error())
	}
	str := "{\"message\":\"" + strings.Join(arr, " ") + "\"}"
	return code, []byte(str)
}
