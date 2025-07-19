package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Dog struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Cat struct {
	ID   int
	Name string
	Toys []Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
	ID        int
	Name      string
	OwnerID   int
	OwnerType string
}

func main() {
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
	if err := db.AutoMigrate(&Toy{}, &Dog{}, &Cat{}); err != nil {
		panic(err)
	}

	// 创建数据
	error := db.Create(&Dog{
		Name: "wangcai",
		Toys: []Toy{
			{Name: "磨牙棒"},
			{Name: "飞盘"},
		},
	}).Error
	if error != nil {
		panic(error)
	}
}
