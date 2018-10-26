package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type Script struct {
	Filename string `json:"filename"`
}

func runHandler(c *gin.Context) {
	var s Script

	if err := c.Bind(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
		return
	}

	dir := strings.Join([]string{"code", s.Filename}, "/")
	// TODO: change
	out, err := exec.Command("python", dir).Output()
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	result := string(out)
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

	router.POST("/", runHandler)

	// TODO: change
	router.Run(":9090")
}
