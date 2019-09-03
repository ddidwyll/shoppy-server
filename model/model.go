package model

import (
	"shoppy/cnst"
	"shoppy/db"

	"encoding/json"
	"html"
	"log"
	"strings"
)

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  byte   `json:"role"`
}

type Password struct {
	Id   string `json:"id"`
	Pass string `json:"pass"`
}

type Product struct {
	Id      string                        `json:"id"`
	Full    string                        `json:"full"`
	Unit    string                        `json:"unit"`
	Price   float64                       `json:"price"`
	Balance float64                       `json:"balance"`
	Type    string                        `json:"type"`
	Stats   map[string]map[string]float64 `json:"stats"`
}

type Order struct {
	Id        string `json:"id"`
	Who       string `json:"who"`
	OrderVals `json:"vals"`
}

type OrderVals map[string]OrderVal

type OrderVal struct {
	Quantity float64 `json:"qt"`
	Price    float64 `json:"pr"`
}

type Modeler interface {
	id() string
	get(string) bool
	validate(string) error
	pattern(string) string
}

type MTime struct {
	Users    string `json:"users"`
	Products string `json:"products"`
	Orders   string `json:"orders"`
}

func Model(t byte) Modeler {
	switch t {
	case cnst.USER:
		return new(User)
	case cnst.PASS:
		return new(Password)
	case cnst.PRDT:
		return new(Product)
	case cnst.ORDR:
		return new(Order)
	}
	return nil
}

func GetOne(t byte, id string) (Modeler, error) {
	obj := Model(t)
	str, err := db.ReadOne(t, id)
	if str == "" {
		return obj, cnst.Error(cnst.NTFND, t)
	}
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal([]byte(str), obj)
	return obj, err
}

func Upsert(t, action byte, m Modeler, user string) error {
	if !db.Exist(t, m.id()) && action == cnst.UPD {
		return cnst.Error(cnst.NTFND, t)
	}
	if db.Exist(t, m.id()) && action == cnst.INS {
		return cnst.Error(cnst.EXIST, t)
	}
	if err := m.validate(user); err != nil {
		return err
	}
	str, err := json.Marshal(&m)
	if err != nil {
		log.Println(m)
		return err
	}
	if err := db.Insert(t, m.id(), string(str)); err != nil {
		return err
	}
	return nil
}

func Delete(t byte, id, user string) error {
	m := Model(t)
	if !m.get(id) {
		return cnst.Error(cnst.NTFND, t)
	}
	if err := m.validate(user); err != nil {
		return err
	}
	return db.Delete(t, id)
}

func GetAll(t byte, user string) (string, error) {
	m := Model(t)
	pattern := m.pattern(user)
	return db.ReadAll(t, pattern)
}

func normalize(str string) string {
	str = strings.Trim(str, " ")
	str = html.EscapeString(str)
	return str
}
