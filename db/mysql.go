package db

import (
	"BookManagementSystem/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var MysqlClient *sqlx.DB

func InitDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/chendong"
	MysqlClient, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	MysqlClient.SetMaxOpenConns(100)
	MysqlClient.SetMaxIdleConns(16)
	return
}


// 插入用户数据
func InsertBmsUser(username , password string) (err error) {
	sqlStr := "insert into bms_user(username, password) value (?,?)"
	_, err = MysqlClient.Exec(sqlStr, username, password)
	if err != nil {
		logrus.Error("插入用户信息失败")
		return
	}
	return
}

// 查询用户登录信息
func QueryUserLogin(username , password string) bool{
	sqlStr := "select password from bms_user where username =?;"
	rows,err := MysqlClient.Query(sqlStr,username)
	if err != nil {
		logrus.Error("query failed")
		fmt.Printf("query failed, err:%v\n", err)
		return false
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var user models.UserInfo
		err := rows.Scan(&user.PassWord)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return false
		}
		if user.PassWord != password{
			return false
		}
	}
	return true
}

// 查数据库
func QueryAllBook() (bookList []*models.Book, err error) {
	sqlStr := "select id, title, price from book;"
	err = MysqlClient.Select(&bookList, sqlStr)
	if err != nil {
		fmt.Println("查询所有书籍信息失败")
		return
	}
	return
}

// 查询单个书籍
func QueryBookByID(id int64) (book models.Book, err error) {
	sqlStr := "select id, title, price from book where id =?;"
	err = MysqlClient.Get(&book, sqlStr, id)
	if err != nil {
		fmt.Println("查询书籍信息失败")
		return
	}
	return
}

// 插入数据
func InsertBook(title string, price float64) (err error) {
	sqlStr := "insert into book(title, price) value (?,?)"
	_, err = MysqlClient.Exec(sqlStr, title, price)
	if err != nil {
		fmt.Println("插入书籍信息失败")
		return
	}
	return
}

// 删除数据
func DeleteBook(id int64) (err error) {
	sqlStr := "delete from book where id = ?"
	_, err = MysqlClient.Exec(sqlStr, id)
	if err != nil {
		fmt.Println("删除书籍信息失败")
		return
	}
	return
}

func EditBook(title string, price float64, id int64) (err error) {
	sqlStr := "update book set title=?,price=? where id =?"
	_, err = MysqlClient.Exec(sqlStr, title, price, id)
	if err != nil {
		fmt.Println("编辑书籍信息失败")
		return
	}
	return
}

