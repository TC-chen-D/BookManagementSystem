package handler

import (
	"BookManagementSystem/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func BookListHandler(c *gin.Context) {
	bookList, err := db.QueryAllBook()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	c.HTML(http.StatusOK, "book/book_list.tmpl", gin.H{
		"code": 0,
		"data": bookList,
	})
}

func NewBookHandler(c *gin.Context) {
	// 给用户返回一个添加书籍页面的处理函数
	c.HTML(http.StatusOK, "book/new_book.html", nil)
}

func CreateBookHandler(c *gin.Context) {
	// 创建书籍的处理函数
	// 从form表单取数据
	titleVal := c.PostForm("title")
	priceVal := c.PostForm("price")
	price, err := strconv.ParseFloat(priceVal, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的价格参数",
		})
		return
	}
	err = db.InsertBook(titleVal, price)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "插入数据失败，请重试！",
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/book/list")
}

func DeleteBookHandler(c *gin.Context) {
	idStr := c.Query("id")
	idVal, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "插入数据失败，请重试！",
		})
		return
	}
	err = db.DeleteBook(idVal)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	// 删除成功跳转到书籍列表页
	c.Redirect(http.StatusMovedPermanently, "/book/list")
}

func EditBookHandler(c *gin.Context) {
	idStr := c.Query("id")
	if len(idStr) == 0 {
		c.String(http.StatusBadRequest, "无效的请求")
		return
	}
	bookID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的请求")
		return
	}
	if c.Request.Method == "POST" {
		// 获取用户提交的数据，去数据库更新对应的书籍数据，跳转回、book/list页面查看是否修改成功
		titleVal := c.PostForm("title")
		priceStr := c.PostForm("price")
		priceVal, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "无效的价格信息")
			return
		}
		err = db.EditBook(titleVal, priceVal, bookID)
		if err != nil {
			c.String(http.StatusInternalServerError, "更新数据失败")
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/book/list")
	} else {
		// 需要给模板渲染上原来的旧数据
		bookObj, err := db.QueryBookByID(bookID)
		if err != nil {
			c.String(http.StatusBadRequest, "无效的书籍id")
			return
		}
		c.HTML(http.StatusOK, "book/book_edit.html", bookObj)
	}
}

func ShowUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func UploadHandler(c *gin.Context) {
	// 提取用户上传的文件
	fileObj, err := c.FormFile("filename")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  err,
		})
		return
	}
	// fileObj: 上传的文件对象
	// fileObj.filename 拿到上传文件的文件名
	filePath := fmt.Sprintf("./%s", fileObj.Filename)
	// 保存文件到本地的路径
	c.SaveUploadedFile(fileObj, filePath)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "OK",
	})
}
