package db

import (
	"github.com/tidwall/buntdb"
	"shoppy/cnst"

	"log"
	"strconv"
	"strings"
	"time"
)

var (
	Users, usersErr         = buntdb.Open("users.db")
	Passwords, passwordsErr = buntdb.Open("passwords.db")
	Products, productsErr   = buntdb.Open("products.db")
	Orders, ordersErr       = buntdb.Open("orders.db")
)

func shooseDB(name byte) *buntdb.DB {
	switch name {
	case cnst.USER:
		return Users
	case cnst.PASS:
		return Passwords
	case cnst.PRDT:
		return Products
	case cnst.ORDR:
		return Orders
	}
	return nil
}

func Exist(db byte, key string) bool {
	if str, _ := ReadOne(db, key); str == "" {
		return false
	}
	return true
}

func lastUpdate(tx *buntdb.Tx) {
	now := time.Now().Unix()
	str := strconv.FormatInt(now, 10)
	tx.Set("A:upd", str, nil)
}

func Insert(db byte, key, val string) (err error) {
	shooseDB(db).Update(
		func(tx *buntdb.Tx) error {
			_, _, err = tx.Set("A:"+key, val, nil)
			log.Println(val)
			lastUpdate(tx)
			return err
		})
	return err
}

func ReadOne(db byte, key string) (str string, err error) {
	shooseDB(db).View(
		func(tx *buntdb.Tx) error {
			str, err = tx.Get("A:" + key)
			return err
		})
	return str, err
}

func ReadAll(db byte, pattern string) (str string, err error) {
	var arr []string
	err = shooseDB(db).View(
		func(tx *buntdb.Tx) error {
			if db != cnst.ORDR {
				return tx.AscendKeys(
					pattern,
					func(key, value string) bool {
						s := "\"" + key[2:] + "\":" + value
						arr = append(arr, s)
						return true
					})
			} else {
				monthAgo := time.Now().Unix() - 2592000
				return tx.DescendKeys(pattern, func(key, value string) bool {
					s := "\"" + key[2:] + "\":" + value
					k, err := strconv.ParseInt(key[2:], 10, 64)
					log.Println("order this month: ", monthAgo, k, k > monthAgo)
					if err != nil || k == 0 {
  					log.Println("pass")
  					return true
					}
					if k > monthAgo {
						arr = append(arr, s)
						return true
					}
					return false
				})
			}
		})
	str = "{" + strings.Join(arr, ",") + "}"
	return str, err
}

func Delete(db byte, key string) error {
	err := shooseDB(db).Update(
		func(tx *buntdb.Tx) error {
			val, err := tx.Delete("A:" + key)
			if err != nil {
				return err
			}
			_, _, err = tx.Set("D:"+key, val, nil)
			lastUpdate(tx)
			return err
		})
	return err
}

func main() {
	if usersErr != nil {
		log.Fatal(usersErr)
	}
	if passwordsErr != nil {
		log.Fatal(passwordsErr)
	}
	if productsErr != nil {
		log.Fatal(productsErr)
	}
	if ordersErr != nil {
		log.Fatal(ordersErr)
	}
	defer Users.Close()
	defer Passwords.Close()
	defer Products.Close()
	defer Orders.Close()
}
