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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/kirtanwyn/allPanelexch/controller"
	"log"
)

func SetupRouter(router *gin.Engine) {
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", "", []byte("secret_session_key"))
	if err != nil {
		log.Fatalf("Failed to connect to Redis session store: %v", err)
	}
	
	// Apply session middleware using PHPSESSID to match PHP behavior
	router.Use(sessions.Sessions("PHPSESSID", store))

	router.GET("/", controller.Home)
	router.GET("/health", controller.Health)

	// Authentication Routes
	// router.POST("/ajaxfiles/logincheck.php", controller.LoginCheck)
	router.POST("/api/login", controller.LoginCheck)
	router.GET("/logout", controller.Logout)
	// router.POST("/logout", controller.Logout)
	router.GET("/get_session", controller.GetSession)

}