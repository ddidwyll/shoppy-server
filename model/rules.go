package model

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"shoppy/cnst"
	"shoppy/db"
)

// MTIME

func (t *MTime) Get() {
	t.Users, _ = db.ReadOne(cnst.USER, "upd")
	t.Products, _ = db.ReadOne(cnst.PRDT, "upd")
	t.Orders, _ = db.ReadOne(cnst.ORDR, "upd")
}

/* USER */

func (u *User) id() string {
	return u.Id
}

func (u *User) get(id string) bool {
	g, err := GetOne(cnst.USER, id)
	n, ok := g.(*User)
	*u = *n
	if err != nil || !ok {
		return false
	}
	return true
}

func (u *User) validate(user string) error {
	u.Id = normalize(u.Id)
	u.Name = normalize(u.Name)
	u.Phone = normalize(u.Phone)
	if u.Role == 0 {
		u.Role = cnst.SOLD
	}
	if u.Id == "" ||
		u.Name == "" {
		return cnst.Error(cnst.EMPTY, cnst.USER)
	}
	c := new(User)
	if !c.get(user) ||
		((c.Id != u.Id) && c.Role != cnst.ADMN) {
		return cnst.Error(cnst.DENY, cnst.USER)
	}
	if !db.Exist(cnst.PASS, u.Id) {
		str := "{\"id\":\"" + u.Id + "\",\"pass\":\"NA\"}"
		db.Insert(cnst.PASS, u.Id, str)
	}
	return nil
}

func (u *User) pattern(user string) string {
	c := new(User)
	if !c.get(user) || c.Role == 0 {
		return "nothing"
	}
	if c.Role == cnst.CLNT {
		return "A:" + c.Id
	}
	return "A:*"
}

/* PASSWORD */

func (p *Password) id() string {
	return p.Id
}

func (p *Password) get(id string) bool {
	g, err := GetOne(cnst.PASS, id)
	n, ok := g.(*Password)
	*p = *n
	if err != nil || !ok {
		return false
	}
	return true
}

func (p *Password) validate(user string) error {
	if p.Id == "" || p.Pass == "" {
		return cnst.Error(cnst.EMPTY, cnst.PASS)
	}
	c, u := new(User), new(User)
	if !c.get(user) || !u.get(p.Id) {
		return cnst.Error(cnst.NTFND, cnst.USER)
	}
	if !((u.Id == c.Id) ||
		(c.Role == cnst.ADMN)) {
		return cnst.Error(cnst.DENY, cnst.PASS)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(p.Pass), 10)
	p.Pass = string(hash)
	return nil
}

func (p *Password) Login() (*User, error) {
	u := new(User)
	if p.Id == "" || p.Pass == "" {
		return u, cnst.Error(cnst.EMPTY, cnst.PASS)
	}
	pass := []byte(p.Pass)
	if !p.get(p.Id) || !u.get(p.Id) {
		return u, cnst.Error(cnst.NTFND, cnst.USER)
	}
	err := bcrypt.CompareHashAndPassword([]byte(p.Pass), pass)
	if err != nil {
		return u, cnst.Error(cnst.NTINS, cnst.PASS)
	}
	return u, nil
}

func (p *Password) pattern(user string) string {
	return "nothing"
}

/* PRODUCT */

func (p *Product) id() string {
	return p.Id
}

func (p *Product) get(id string) bool {
	g, err := GetOne(cnst.PRDT, id)
	n, ok := g.(*Product)
	*p = *n
	if err != nil || !ok {
		return false
	}
	return true
}

func (p *Product) validate(user string) error {
	p.Id = normalize(p.Id)
	p.Type = normalize(p.Type)
	p.Full = normalize(p.Full)
	p.Unit = normalize(p.Unit)
	if p.Id == "" || p.Price == 0 {
		log.Println(p)
		return cnst.Error(cnst.EMPTY, cnst.PRDT)
	}
	u := new(User)
	if !u.get(user) || u.Role < cnst.SOLD {
		return cnst.Error(cnst.DENY, cnst.PRDT)
	}
	return nil
}

func (p *Product) pattern(user string) string {
	//u := new(User)
	//if !u.get(user) || u.Role == 0 {
	//return "nothing"
	//}
	return "A:*"
}

/* ORDER */

func (o *Order) id() string {
	return o.Id
}

func (o *Order) get(id string) bool {
	g, err := GetOne(cnst.ORDR, id)
	n, ok := g.(*Order)
	*o = *n
	if err != nil || !ok {
		return false
	}
	return true
}

func (o *Order) validate(user string) error {
	o.Who = user
	u := new(User)
	if !u.get(user) ||
		u.Role < cnst.SOLD {
		return cnst.Error(cnst.DENY, cnst.ORDR)
	}
	return nil
}

func (o *Order) pattern(user string) string {
	u := new(User)
	if !u.get(user) || u.Role == 0 {
		return "nothing"
	}
	return "A:*"
}
