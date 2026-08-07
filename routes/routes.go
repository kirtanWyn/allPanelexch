// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/kirtanwyn/subsidy-api-go/controller"
// )

// func SetupRouter() *gin.Engine {

// 	router := gin.Default()

// 	router.GET("/", controller.Home)
// 	router.GET("/health", controller.Health)
	

// 	return router
// }
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/controller"
)

func SetupRouter(router *gin.Engine) {

	router.GET("/", controller.Home)
	router.GET("/health", controller.Health)

	// Authentication Routes
	router.POST("/ajaxfiles/logincheck.php", controller.LoginCheck)
	router.POST("/api/login", controller.LoginCheck)

}