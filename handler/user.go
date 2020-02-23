package handler

import (
	"BookManagementSystem/db"
	"BookManagementSystem/models"
	"BookManagementSystem/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user/login.html", nil)
}

func LoginHandler(c *gin.Context) {
	var u models.UserInfo
	err := c.ShouldBind(&u)
	if err != nil {
		//c.JSON(http.StatusOK,gin.H{
		//	"err": err.Error(),
		//})
		c.HTML(http.StatusOK,"user/login.html",gin.H{
			"err":"用户名或密码不能为空",
		})
		return
	}
	if u.UserName == "" || u.PassWord == "" {
		c.HTML(http.StatusBadRequest,"user/login.html",gin.H{
			"err":"无效的用户信息",
		})
		//c.String(http.StatusBadRequest, "无效的用户信息")
		return
	} else {
		ok := db.QueryUserLogin(u.UserName, u.PassWord)
		if !ok {
			c.HTML(http.StatusMovedPermanently,"user/login.html",gin.H{
				"err":"用户名或密码不对",
			})
			//c.Redirect(http.StatusMovedPermanently, "/user/login")
		}
		c.Redirect(http.StatusMovedPermanently, "/book/list")
	}

}

func ShowRegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user/register.html", nil)
}

func UserRegisterHandler(c *gin.Context) {
	userName := c.PostForm("username")
	passWord := c.PostForm("password")
	pwd := util.Md5(passWord)
	err := db.InsertBmsUser(userName, pwd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "user/login.html", nil)
}
