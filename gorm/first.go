package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main1() {
	// 配置数据库连接参数
	username := "root"
	password := "147258"
	host := "127.0.0.1"
	port := 3306
	Dbname := "hello_gorm"
	timeout := "10s"

	// 拼接dsn参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	// 连接成功
	fmt.Println("数据库连接成功")

	// 创建表
	// autoMigrateErro := db.AutoMigrate(&User{})
	// fmt.Println("autoMigrateErro= %v\n", autoMigrateErro)

	// if autoMigrateErro != nil {
	// 	panic("迁移表error=" + autoMigrateErro.Error())
	// }

	//user := User{Name: "张三", Age: 18, Birthday: time.Now()}
	// user := User{Name: "Jinzhu", Age: 18, Birthday: new(time.Time)}
	// result := db.Create(&user)
	// fmt.Println("result = ", result)

	// users := []*User{
	// 	{Name: "Jinzhu", Age: 18, Birthday: time.Now()},
	// 	{Name: "Jackson", Age: 19, Birthday: time.Now()},
	// }
	// result := db.Create(&users)
	// fmt.Println("影响行数:", result.RowsAffected)
	// fmt.Println("错误信息:", result.Error)

	// var userFirst User
	// errFirst := db.Debug().First(&userFirst).Error
	// fmt.Printf("userFirst= %+v\n", userFirst)
	// fmt.Printf("errFirst= %+v\n", errFirst)
	// fmt.Printf("---------------------------------")

	var usersTwo []User
	errTwo := db.Debug().Where(&User1{Name: "jinzhu"}, "name", "Age").Find(&usersTwo)
	jsonData, _ := json.MarshalIndent(usersTwo, "", "  ")
	fmt.Println(string(jsonData))
	jsonDataError, _ := json.MarshalIndent(errTwo, "", "  ")
	fmt.Printf("jsonDataError= %+v\n", jsonDataError)
}

type Product1 struct {
	gorm.Model
	Code  string
	Price uint
}

type User1 struct {
	ID           uint           // Standard field for the primary key
	Name         string         // A regular string field
	Email        *string        // A pointer to a string, allowing for null values
	Age          uint8          // An unsigned 8-bit integer
	Birthday     time.Time      // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // Automatically managed by GORM for creation time
	UpdatedAt    time.Time      // Automatically managed by GORM for update time
	ignored      string         // fields that aren't exported are ignored
}
