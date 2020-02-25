package main

import (
	//"gopkg.in/go-playground/validator.v8"
	"os"
	"time"

	"BookManagementSystem/db"
	"BookManagementSystem/handler"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// 日志输出初始化
func InitLogs() {
	log.SetFormatter(&log.TextFormatter{ // 设置日志格式为text格式
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式
	})
	log.SetOutput(os.Stdout)     // 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	log.SetLevel(log.DebugLevel) // 设置日志级别为debug以上
	log.SetReportCaller(true)

	hook := newLfsHook("debug") // 创建hook
	log.AddHook(hook)
}

// 本地日志文件分割，按天存储
func newLfsHook(logLevel string) log.Hook {
	writer, err := rotatelogs.New(
		"./log/log.%Y%m%d%H%M", // 文件名

		// rotatelogs.WithLinkName("./log/log"), // WithLinkName为最新的日志建立软连接

		rotatelogs.WithRotationTime(24*time.Hour), // WithRotationTime设置日志分割的时间

		// WithMaxAge和WithRotationCount二者只能设置一个
		rotatelogs.WithRotationCount(7), // WithRotationCount设置文件清理前最多保存的个数
		// rotatelogs.WithMaxAge(24*time.Hour),        // WithMaxAge设置文件清理前的最长保存时间
	)
	if err != nil {
		log.Fatalf("配置存储日志文件出错: %v", err)
	}

	level, err := log.ParseLevel(logLevel)
	if err == nil {
		log.SetLevel(level)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{DisableColors: true})

	return lfsHook
}



func main() {
	InitLogs()
	err := db.InitDB()
	if err != nil {
		log.Error("connect mysql fail", err.Error())
	}
	r := gin.Default()
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	v.RegisterValidation("bookabledate", BookableDate)
	//}
	r.LoadHTMLGlob("templates/**/*")
	userGroup := r.Group("/user")
	{
		userGroup.GET("/register", handler.ShowRegisterHandler)
		userGroup.POST("/register", handler.UserRegisterHandler)
		userGroup.GET("/login", handler.ShowLoginHandler)
		userGroup.POST("/login", handler.LoginHandler)
	}

	bookGroup := r.Group("/book")
	{
		// 查看所有书籍数据
		bookGroup.GET("/list", handler.BookListHandler)
		// 返回一个页面给用户填写新增的书籍信息
		bookGroup.GET("/new", handler.NewBookHandler)
		bookGroup.POST("/new", handler.CreateBookHandler)
		bookGroup.GET("/delete", handler.DeleteBookHandler)
		bookGroup.Any("/edit", handler.EditBookHandler)
		bookGroup.GET("/upload", handler.ShowUpload)
		bookGroup.POST("/upload", handler.UploadHandler)
	}

	learnGroup := r.Group("/test1")
	{
		learnGroup.GET("/:name/:id", handler.LearnExample1Handler)
	}
	learn2Group := r.Group("test2")
	{
		learn2Group.GET("/example1", handler.LearnExample2Handler)
		learn2Group.GET("/example2", handler.LearnExample3Handler)
		learn2Group.GET("/bookable", handler.LearnExample4Handler)
	}

	r.Run(":9090")
}
