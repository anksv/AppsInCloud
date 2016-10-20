package db

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

// Model Struct
type User struct {
	Id   int
	Name string `orm:"size(100)"`
}

func init() {
	// register model
	//orm.RegisterModel(new(models.App))

	// set default database
	//postgresql: //postgre:huawei123@localhost:5432/postgres
	orm.RegisterDataBase("default", "postgres", "postgres://postgres:huawei123@localhost:5432/default?sslmode=disable", 30)

	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		fmt.Println(err)
	}

}

func InitMydb() {
	/*o := orm.NewOrm()

	user := User{Name: "slene"}

	// insert
	id, err := o.Insert(&user)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	// update
	user.Name = "astaxie"
	num, err := o.Update(&user)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

	// read one
	u := User{Id: user.Id}
	err = o.Read(&u)
	fmt.Printf("ERR: %v\n", err)

	// delete
	num, err = o.Delete(&u)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)*/
}
