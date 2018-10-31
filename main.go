package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"path/filepath"
	"path"
	"github.com/spf13/viper"
)

var (
	UploadDir       string
	URLPrefix       string
	EnableBasicAuth bool
	UserName        string
	Password        string
)

func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/file_transfer")
	viper.SetConfigName("file_transfer")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("can't read config file")
		panic(err)
	}
	UploadDir = viper.GetString("upload_dir")
	URLPrefix = viper.GetString("url_prefix")
	EnableBasicAuth = viper.GetBool("basic_auth.enabled")
	UserName = viper.GetString("basic_auth.username")
	Password = viper.GetString("basic_auth.password")

	log.Println("UploadDir", UploadDir)
	log.Println("URLPrefix", URLPrefix)
	log.Println("EnableBasicAuth", EnableBasicAuth)
	log.Println("starting...")
}

func DealUpload(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	// Upload the file to specific dst.
	dst := filepath.Join(UploadDir, file.Filename)
	c.SaveUploadedFile(file, dst)
	c.String(http.StatusOK, "DownloadPath: %s\n",
		path.Join(URLPrefix, file.Filename))
}

func main() {
	router := gin.Default()
	if EnableBasicAuth {
		router.Use(gin.BasicAuth(gin.Accounts{
			UserName: Password, //用户名：密码
		}))
	}
	router.StaticFS("/download", http.Dir(UploadDir))
	router.POST("/upload", DealUpload)
	router.Run(":8080")
}
