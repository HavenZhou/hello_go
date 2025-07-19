package main

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main2() {
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
	//db.AutoMigrate(&Language{})
	//db.AutoMigrate(&User{})

	// 创建数据
	// user := User{}
	// db.Create(&user)

	// language1 := Language{Name: "English"}

	// language2 := Language{Name: "Chinese"}
	// db.Create(&[]Language{language1, language2})

	// user1 := User{Languages: []Language{language1, language2}}
	// db.Create(&user1)

	// var usersTwo []User
	// errTwo := db.Debug().Where(&User{user.Languages: "jinzhu"}, "name", "Age").Find(&usersTwo)
	// jsonData, _ := json.MarshalIndent(usersTwo, "", "  ")
	// fmt.Println(string(jsonData))
	// jsonDataError, _ := json.MarshalIndent(errTwo, "", "  ")
	// fmt.Printf("jsonDataError= %+v\n", jsonDataError)

	// 查询数据
	// user := User{Model: gorm.Model{ID: 3}}
	// error := db.Preload("Languages").Find(&user).Error
	// if error != nil {
	// 	panic(error)
	// }
	// jsonData, _ := json.MarshalIndent(user, "", "  ")
	// fmt.Println(string(jsonData))
	// fmt.Println("======================")

	// Count all languages
	// user1 := User{Model: gorm.Model{ID: 3}}
	// count1 := db.Model(&user1).Association("Languages").Count()
	// jsonData1, _ := json.MarshalIndent(count1, "", "  ")
	// fmt.Println(string(jsonData1))
	// fmt.Println("======================")

	// // Count with conditions
	// user2 := User{Model: gorm.Model{ID: 3}}
	// codes := []string{"zh-CN", "en-US", "ja-JP"}
	// count2 := db.Model(&user2).Where("Name IN ?", codes).Association("Languages").Count()
	// jsonData2, _ := json.MarshalIndent(count2, "", "  ")
	// fmt.Println(string(jsonData2))
	// fmt.Println("======================")

	// 更新语句
	var user = User{Model: gorm.Model{ID: 3}}
	err1 := db.Debug().Preload(clause.Associations).Find(&user).Error
	if err1 != nil {
		panic(err1)
	}

	user.Languages[0].Name = "李四1111111"
	jsonData, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println("=============update 更新操作=========")
	// err2 := db.Debug().Updates(&user).Error // Name不生效
	err2 := db.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error
	if err2 != nil {
		panic(err2)
	}

}

// User 拥有并属于多种 language，`user_languages` 是连接表
type User struct {
	gorm.Model
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
	Name string
}
