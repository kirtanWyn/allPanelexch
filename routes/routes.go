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
	router.POST("/signup", controller.Signup)
	router.POST("/verify_otp", controller.VerifySignupOTP)
	router.POST("/login", controller.LoginCheck)
	router.POST("/change_password", controller.ChangePassword)
	router.GET("/logout", controller.Logout)
	// router.POST("/logout", controller.Logout)
	router.GET("/get_session", controller.GetSession)

	// Account Statements
	router.POST("/ajaxfiles/account_statement.php", controller.GetAccountStatement)
	router.POST("/ajaxfiles/account_bet_statement.php", controller.GetAccountBetStatement)
	router.POST("/account_statement", controller.GetAccountStatement)
	router.POST("/account_bet_statement", controller.GetAccountBetStatement)
	router.POST("/ajaxfiles/current_bet.php", controller.GetCurrentBets) 
	router.POST("/current_bet", controller.GetCurrentBets)

	// Button Value Change
	router.POST("/ajaxfiles/button_value_change.php", controller.UpdateButtonValue)
	router.POST("/button_value_change", controller.UpdateButtonValue)

	// Activity Logs
	router.POST("/ajaxfiles/activity_log.php", controller.GetActivityLogs)
	router.POST("/activity_log", controller.GetActivityLogs)

	// Casino Results
	router.POST("/ajaxfiles/casino_result.php", controller.GetCasinoResult)
	router.POST("/casino_result", controller.GetCasinoResult)
	router.POST("/ajaxfiles/get_result_cards.php", controller.GetResultCards)
	router.POST("/get_result_cards", controller.GetResultCards)

	// Balance
	router.POST("/ajaxfiles/refresh_balance.php", controller.RefreshBalance)
	router.POST("/refresh_balance", controller.RefreshBalance)
}