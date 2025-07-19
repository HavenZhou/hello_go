package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	connectGormData()
	fmt.Println("连接数据库成功")
	MigrateModels(db)
	// 初始化数据
	initTestData()

	getUserPostsWithComments(1)
	fmt.Printf("--------------------------------------------")

	getMostCommentedPost()

	HookTest()

}

// 进阶gorm
// 题目1：模型定义
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
// UserBlog 用户模型
type UserBlog struct {
	gorm.Model
	Username  string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `gorm:"type:varchar(100);not null"`
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多关系：一个用户有多篇文章
	PostCount int    `gorm:"default:0"`         // 新增：用户文章计数
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(200);not null"`
	Content  string    `gorm:"type:text;not null"`
	UserID   uint      `gorm:"index;not null"`    // 外键，关联UserBlog
	UserBlog UserBlog  `gorm:"foreignKey:UserID"` // 属于关系：文章属于用户
	Comments []Comment `gorm:"foreignKey:PostID"` // 一对多关系：一篇文章有多个评论
	// 新增：评论状态
	CommentStatus string `gorm:"type:varchar(20);default:'有评论'"`
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content  string   `gorm:"type:text;not null"`
	PostID   uint     `gorm:"index;not null"`    // 外键，关联Post
	Post     Post     `gorm:"foreignKey:PostID"` // 属于关系：评论属于文章
	UserID   uint     `gorm:"index;not null"`    // 外键，关联UserBlog（评论作者）
	UserBlog UserBlog `gorm:"foreignKey:UserID"` // 属于关系：评论属于用户
}

// 数据库连接
var db *gorm.DB

func connectGormData() {
	// 配置数据库连接参数
	username := "root"
	password := "147258"
	host := "127.0.0.1"
	port := 3306
	Dbname := "hello_gorm"
	timeout := "10s"

	// 拼接dsn参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, err=" + err.Error())
	} else {
		fmt.Printf("数据库连接成功db = %v\n", db)
	}
}

// MigrateModels 迁移模型到数据库
func MigrateModels(db *gorm.DB) error {
	// 注册要迁移的模型
	models := []interface{}{
		&UserBlog{},
		&Post{},
		&Comment{},
	}

	// 执行迁移
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	} else {
		fmt.Println("建表成功")
	}

	return nil
}

// 题目2：关联查询
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
// 数据初始化
func initTestData() {
	// 清空现有数据
	// db.Exec("DELETE FROM comments")
	// db.Exec("DELETE FROM posts")
	// db.Exec("DELETE FROM user_blogs")

	// 创建测试用户
	users := []UserBlog{
		{Username: "user1", Email: "user1@example.com", Password: "pass1"},
		{Username: "user2", Email: "user2@example.com", Password: "pass2"},
	}
	db.Create(&users)

	// 创建测试文章
	posts := []Post{
		{Title: "First Post", Content: "Content of first post", UserID: users[0].ID},
		{Title: "Second Post", Content: "Content of second post", UserID: users[0].ID},
		{Title: "Third Post", Content: "Content of third post", UserID: users[1].ID},
	}
	db.Create(&posts)

	// 创建测试评论
	comments := []Comment{
		{Content: "Great post!", PostID: posts[0].ID, UserID: users[1].ID},
		{Content: "Nice article", PostID: posts[0].ID, UserID: users[1].ID},
		{Content: "Interesting", PostID: posts[1].ID, UserID: users[1].ID},
		{Content: "Well written", PostID: posts[2].ID, UserID: users[0].ID},
		{Content: "Thanks for sharing", PostID: posts[2].ID, UserID: users[0].ID},
		{Content: "I learned a lot", PostID: posts[2].ID, UserID: users[0].ID},
	}
	db.Create(&comments)

	fmt.Println("测试数据初始化完成！")
}

// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func getUserPostsWithComments(userID uint) {
	var user UserBlog
	err := db.Preload("Posts.Comments.UserBlog"). // 预加载用户的所有文章及文章的评论(包括评论作者)
							Preload("Posts.UserBlog"). // 预加载每篇文章的作者信息
							First(&user, userID).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用户文章 %v\n,文章的评论信息如下:", user.Username)
	for _, post := range user.Posts {
		fmt.Printf("文章标题: %s\n", post.Title)
		fmt.Printf("文章内容: %s\n", post.Content)
		fmt.Println("评论:")
		for _, comment := range post.Comments {
			fmt.Printf("- %s (评论人: %s)\n", comment.Content, comment.UserBlog.Username)
		}
		fmt.Println("-------------------")
	}
}

// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
// 查询评论数量最多的文章信息
func getMostCommentedPost() {
	var post Post

	// 第一步：查询评论数量最多的文章ID
	err := db.Debug().Model(&Post{}).
		Select("posts.id, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&post).Error
	if err != nil {
		fmt.Errorf("查询评论最多文章失败: %v", err)
	}
	fmt.Printf("评论最多的文章: %v\n-------------------------------", post.ID)

	// 第二步：完整加载文章及其关联信息
	err = db.Debug().Model(&Post{}).
		Preload("UserBlog").          // 加载文章作者
		Preload("Comments").          // 加载文章评论
		Preload("Comments.UserBlog"). // 加载评论作者
		First(&post, post.ID).        // 查询特定ID的文章
		Error
	if err != nil {
		fmt.Errorf("加载文章详情失败: %v\n", err)
	}

	if err != nil {
		log.Fatalf("查询失败: %v\n", err)
	}

	fmt.Printf("评论最多的文章:\n")
	fmt.Printf("ID: %d\n", post.ID)
	fmt.Printf("标题: %s\n", post.Title)
	fmt.Printf("作者: %s\n", post.UserBlog.Username)
	fmt.Printf("评论数量: %d\n", len(post.Comments))

	fmt.Println("\n评论列表:")
	for _, comment := range post.Comments {
		fmt.Printf("- %s (by %s)\n", comment.Content, comment.UserBlog.Username)
	}

	fmt.Println("-------------------")

}

// 题目3：钩子函数
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
// Post 模型的 AfterCreate 钩子
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新用户的文章计数
	fmt.Println("------AfterCreate,HOOK IS START-------------")

	result := tx.Debug().Model(&UserBlog{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + 1"))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Comment 模型的 AfterDelete 钩子
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("------AfterDelete,HOOK IS START-------------")

	// 检查文章的剩余评论数量
	var count int64
	result := tx.Debug().Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count)

	if result.Error != nil {
		return result.Error
	}

	// 如果没有评论了，更新文章状态
	if count == 0 {
		result = tx.Debug().Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "无评论")

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// 钩子测试
func HookTest() {
	// 清空数据
	// db.Exec("DELETE FROM comments")
	// db.Exec("DELETE FROM posts")
	// db.Exec("DELETE FROM user_blogs")

	// 创建测试用户
	user := UserBlog{
		Username: "hook_user",
		Email:    "hook@example.com",
		Password: "hookpass",
	}
	db.Create(&user)

	fmt.Printf("创建用户后文章数: %d\n", user.PostCount)

	// 测试 Post 钩子
	post := Post{
		Title:   "钩子测试文章",
		Content: "测试钩子函数",
		UserID:  user.ID,
	}
	db.Create(&post)

	// 重新加载用户数据查看计数
	db.First(&user, user.ID)
	fmt.Printf("创建文章后用户文章数: %d\n", user.PostCount)

	// 测试 Comment 钩子
	comment1 := Comment{
		Content: "测试评论1",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	comment2 := Comment{
		Content: "测试评论2",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	db.Create(&comment1)
	db.Create(&comment2)

	// 检查文章状态
	var currentPost Post
	db.First(&currentPost, post.ID)
	fmt.Printf("添加评论后文章状态: %s\n", currentPost.CommentStatus)

	// 删除评论
	db.Delete(&comment1)
	db.Delete(&comment2)

	// 再次检查文章状态
	db.First(&currentPost, post.ID)
	fmt.Printf("删除所有评论后文章状态: %s\n", currentPost.CommentStatus)
}
