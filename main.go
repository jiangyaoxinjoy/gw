package main

import (
	"gw2/controllers"
	"gw2/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	baseController := &controllers.BaseController{}
	r := gin.Default()
	routers.CreateRouters(r, baseController)
	r.Run("127.0.0.1:4699")
}
