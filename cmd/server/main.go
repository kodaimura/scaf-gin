package main

import (
	"github.com/gin-gonic/gin"

	"goscaf/internal/router"
)

func main() {
	r := gin.Default()
	router.Set(r)
	r.Run(":3000")
}