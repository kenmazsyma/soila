package api

import (
	"errors"
	"fmt"
	_ "github.com/kenmazsyma/soila/peer/blockchain"
	"github.com/kenmazsyma/soila/peer/db"
)

type Person struct {
}

type Person_Create struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Pass    string `json:"pass"`
	Profile string `json:"profile"`
}

func NewPerson() *Person {
	return &Person{}
}

func (self *Person) Create(param Person_Create, rslt *string) (err error) {
	err = db.Person_Create(param.Id, param.Name, param.Pass, param.Profile)
	if err != nil {
		fmt.Printf("DBError:%s\n", err.Error())
		return errors.New("作成に失敗しました。")
	}
	(*rslt) = "value"
	return nil
}
