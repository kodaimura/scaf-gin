package main

import (
	"io"
	"os"
	"github.com/gin-gonic/gin"

	"goscaf/pkg/logger"
	"goscaf/internal/router"
)

func main() {
    f, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
    gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	logger.SetOutput(f)
	logger.SetLevel("INFO")

	r := gin.Default()
	router.SetStatic(r)
	router.SetWeb(r)
	router.SetApi(r)
	r.Run(":3000")
}