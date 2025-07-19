package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main3() {
	// 创建连接
	db := connectData()

	// // 创建表
	// if err := db.AutoMigrate(&Student{}); err != nil {
	// 	panic(err)
	// }

	// // 题目1：题目1：基本CRUD操作
	// q1(db)

	// 题目2：转账事务处理
	// if err := db.AutoMigrate(&Account{}, &Transaction{}); err != nil {
	// 	panic(err)
	// }

	// 初始化测试账户
	// if err := initTestAccounts(db); err != nil {
	// 	log.Fatalf("初始化账户失败: %v", err)
	// }

	// 执行转账：从账户1向账户2转账100元
	if err := q2(db, 1, 2, 100); err != nil {
		log.Printf("转账失败: %v", err)
	} else {
		log.Println("转账成功")
	}

}

type Student struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"` // 明确指定主键和自增
	Name  string `gorm:"type:varchar(100)"`
	Age   uint   `gorm:"type:int"`
	Grade string `gorm:"type:varchar(50)"`
}

func connectData() *gorm.DB {
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
	return db
}

// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 要求 ：
// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
func q1(db *gorm.DB) {
	newStudent := Student{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&newStudent)

	if result.Error != nil {
		panic(result.Error)
	}

	// 查询
	var students []Student
	result1 := db.Where("age > ?", 18).Find(&students)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}

	for _, student := range students {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n",
			student.ID, student.Name, student.Age, student.Grade)
	}

	// 更新
	result2 := db.Model(&Student{}).Where("name =?", "张三").Update("grade", "四年级")
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	fmt.Printf("更新了%d条记录\n", result2.RowsAffected)
	// student, _ := json.MarshalIndent(newStudent, "", "  ")
	// fmt.Printf("插入成功,student= %+v\n", student)

	// 删除
	result3 := db.Where("age < ?", 15).Delete(&Student{})
	if result3.Error != nil {
		log.Fatal(result3.Error)
	}
	fmt.Printf("删除了%d条记录\n", result3.RowsAffected)
}

// 题目2：事务语句
//
//	假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）
//	和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
//
// 要求 ：
//
//	编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
//	在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
//	并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
type Account struct {
	ID      uint `gorm:"primaryKey"`
	Balance float64
}

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

func q2(db *gorm.DB, fromId uint, toId uint, amount float64) error {
	// 开启事务
	tx := db.Begin()

	// 检查错误
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %v", tx.Error)
	}

	// 确保在panic时回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1.检查转出账户是否存在及余额是否足够
	var fromAccount Account
	if err := tx.First(&fromAccount, fromId).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转出账户(d%)失败，%v", fromId, err)
	}
	fromAccountJson, _ := json.MarshalIndent(fromAccount, "", "  ")
	fmt.Println(string(fromAccountJson))

	if fromAccount.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("账户(d%)余额不足，当前余额：%.2f", fromId, fromAccount.Balance)
	}

	// 2.检查转入账户是否存在
	var ToAccount Account
	if err := tx.First(&ToAccount, toId).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转入账户(%d)失败: %v", toId, err)
	}

	// 3.扣除转出账户余额
	if err := tx.Model(&Account{}).
		Where("id =? AND balance >= ?", fromId, amount). // 乐观锁
		Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("增加转入账户(%d)金额失败: %v", toId, err)
	}

	// 4.增加转入金额
	if err := tx.Model(&Account{}).
		Where("id = ?", toId).
		Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("增加转入账户(%d)金额失败: %v", toId, err)
	}

	// 5记录交易信息
	if err := tx.Create(&Transaction{FromAccountID: fromId, ToAccountID: toId, Amount: amount}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录交易失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// 初始化测试账户
func initTestAccounts(db *gorm.DB) error {
	// 清空表
	if err := db.Exec("TRUNCATE TABLE accounts").Error; err != nil {
		return err
	}
	if err := db.Exec("TRUNCATE TABLE transactions").Error; err != nil {
		return err
	}

	// 创建测试账户
	accounts := []Account{
		{ID: 1, Balance: 500},
		{ID: 2, Balance: 300},
	}
	for _, acc := range accounts {
		if err := db.Create(&acc).Error; err != nil {
			return err
		}
	}
	return nil
}
