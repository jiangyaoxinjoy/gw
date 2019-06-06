package routers

import (
	"gw2/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRouters(r *gin.Engine, tc *controllers.BaseController) {
	r.POST("/api", tc.Get)
}
