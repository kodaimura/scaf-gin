package main

import (
	"github.com/gin-gonic/gin"

	"goscaf/internal/router"
)

func main() {
	r := gin.Default()
	router.SetStatic(r)
	router.SetWeb(r)
	router.SetApi(r)
	r.Run(":3000")
}