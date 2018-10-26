package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Answer struct {
	ID       string // server で割り振る
	Language string `json:"language"`
	Script   string `json:"script"`
}

func (a *Answer) filename() string { // TODO: 固定する
	filename := strings.Join([]string{a.ID, "py"}, ".")
	return filename
}

func (a *Answer) dirname() string { // TODO: 固定する
	filename := a.filename()
	dir := strings.Join([]string{"code", filename}, "/")
	return dir
}

func (a *Answer) save() error {
	filename := a.dirname()
	return ioutil.WriteFile(filename, []byte(a.Script), 0600)
}

func (a *Answer) run() string {
	// golang:1.9にもともと入っているpythonで実行(version: 2.7.13)
	filename := a.dirname()
	out, err := exec.Command("python", filename).Output()
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	return string(out)
}

func (a *Answer) remove() {
	filename := a.dirname()
	if err := os.Remove(filename); err != nil {
		fmt.Println(err)
	}
}

type ResponceFromRunner struct {
	Result string
}

func (a *Answer) runOnDocker() string {
	// 実行するサーバにfile名をリクエスト
	// TODO:Data Volumeから読み取り、python-runnerで実行
	filename := a.filename()

	// リクエストの作成
	p := `{"filename":""}`
	index := 13
	q := p[:index] + filename + p[index:]

	var jsonStr = []byte(q)
	req, err := http.NewRequest("POST", "http://python-runner:9090/", bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := ResponceFromRunner{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	return string(result.Result)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func pythonHandler(c *gin.Context) {
	// ランダムなidの生成
	id_ := RandStringRunes(16)

	a := &Answer{ID: id_}
	if err := c.Bind(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		return
	}

	err := a.save()
	if err != nil {
		panic(err)
	}
	defer a.remove()

	//result := a.run()
	result := a.runOnDocker()
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/python", pythonHandler)

	router.Run(":7070")
}
