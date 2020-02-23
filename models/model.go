package models

// 专门用来定义与数据库对应的结构体

type Book struct {
	ID    int64   `db:"id"`
	Title string  `db:"title"`
	Price float64 `db:"price"`
}

type UserInfo struct {
	UserName string `form:"username" json:"username"`
	PassWord string `form:"password" json:"password"`
}
