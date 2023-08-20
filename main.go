package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	Viper()
	r := gin.Default()

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/public/index.html")
	})

	r.POST("/login", func(c *gin.Context) {
		if c.PostForm("PASSWORD") != viper.GetString("PASSWORD") {
			Error(c, "密码错误")
			return
		}
		Success(c, "登陆成功")
	})

	r.POST("/chat", func(c *gin.Context) {
		if c.PostForm("PASSWORD") != viper.GetString("PASSWORD") {
			Error(c, "密码错误")
			return
		}

		ask := c.PostForm("ask")
		user := c.PostForm("user")
		client := &http.Client{
			Timeout: time.Second * 60,
		}
		reqData := `{"temperature": 0, "user": "` + base64.StdEncoding.EncodeToString([]byte(user)) + `", "messages":` + ask + `}`
		fmt.Println(reqData)
		req, err := http.NewRequest("POST", viper.GetString("ENDPOINT"), strings.NewReader(reqData))
		if err != nil {
			Error(c, "生成请求失败："+err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("api-key", viper.GetString("AZURE_OPENAI_KEY"))
		res, err := client.Do(req)

		if res != nil {
			defer res.Body.Close()
		}
		if err != nil {
			Error(c, "获取响应失败："+err.Error())
			return
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			Error(c, "读取响应失败："+err.Error())
			return
		}

		fmt.Println(string(data))

		type Message struct {
			Content string `json:"content"`
		}

		type Messages struct {
			Message Message `json:"message"`
		}

		type Response struct {
			Choices []Messages      `json:"choices"`
			Error   json.RawMessage `json:"error"`
		}

		var resp Response
		if err = json.Unmarshal(data, &resp); err != nil {
			Error(c, "Unmarshal failed："+err.Error())
			return
		}

		if len(resp.Choices) == 0 {
			if len(resp.Error) != 0 {
				Error(c, string(resp.Error))
				return
			}
			Error(c, "resp.Choices length is 0")
			return
		}

		message := resp.Choices[0].Message.Content
		if message == "" {
			Error(c, "no result")
			return
		}

		Success(c, message)
	})
	r.Run(":8989") // listen and serve on 0.0.0.0:8080
}

func Error(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"type":    "error",
		"message": message,
	})
}

func Success(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"type":    "success",
		"message": message,
	})
}

func Viper() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".") // optionally look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		errMessage := "未找到配置文件 .env ：" + err.Error()
		log.Fatalln(errMessage)
	}
	viper.WatchConfig()
}
