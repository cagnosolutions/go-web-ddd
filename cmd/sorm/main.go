package main

import (
	"fmt"
	"github.com/cagnosolutions/go-web-ddd/pkg/sorm"
)

func main() {

	t := sorm.MakeTable(&User{})
	fmt.Println(t.CreateString())
	fmt.Println(t.DropString())
	fmt.Println(t.SelectString("*"))
	fmt.Println(t.InsertString())
	fmt.Println(t.UpdateString())
}

type User struct {
	Id           int    `sql:"id,pk"`
	FirstName    string `sql:"first_name"`
	LastName     string `sql:"last_name"`
	EmailAddress string `sql:"email_address"`
}
