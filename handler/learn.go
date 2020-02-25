package handler

import (
	"BookManagementSystem/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 参数作为url
func LearnExample1Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": c.Param("name"),
		"id":   c.Param("id"),
	})
}

// 获取参数、带默认值
func LearnExample2Handler(c *gin.Context) {
	firstName := c.Query("first_name")
	lastName := c.DefaultQuery("last_name", "last_default_name")
	c.String(http.StatusOK, "%s,%s", firstName, lastName)
}

// 验证请求参数，结构体验证
func LearnExample3Handler(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBind(&person); err != nil {
		c.String(http.StatusInternalServerError, "%v", err)
		c.Abort()
		return
	}
	c.String(http.StatusOK, "%v", person)
}


func LearnExample4Handler(c *gin.Context) {
	validate := models.RegisterTagToValidate()
	var book models.Booking
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	if err := validate.Struct(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"booking": book,
	})

}
