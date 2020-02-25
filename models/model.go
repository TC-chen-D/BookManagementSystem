package models

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

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

// 结构体验证，required必须需要，gt大于
type Person struct {
	Name    string `form:"name" binding:"required"`
	Age     int    `form:"age" binding:"required,gt=10"`
	Address string `form:"address" binding:"required"`
}

//binding 绑定一些验证请求参数,自定义标签bookabledate表示可预约的时期
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

//定义bookabledate标签对应的验证方法
func BookableDate(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if date.Unix() < today.Unix() {
			return true
		}
	}
	return false
}

// 将该验证方法注册到validator验证器里面
func RegisterTagToValidate() *validator.Validate {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", BookableDate)
		return v
	}
	return nil
}
