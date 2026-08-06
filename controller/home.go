package controller

import ( 
	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/utils"
	"github.com/kirtanwyn/allPanelexch/config"
)
func Home(c *gin.Context) {

	// c.JSON(200, gin.H{
	// 	"status":  true,
	// 	"message": "Subsidy API is Running",
	// })
	utils.Success(c, "Subsidy API Running", nil)
} 
    

func Health(c *gin.Context) {

	dbStatus := "connected"
	redisStatus := "connected"

	// Check Database
	if err := config.DB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	// Check Redis
	if _, err := config.Redis.Ping(c.Request.Context()).Result(); err != nil {
		redisStatus = "disconnected"
	}

	// data : {"database": dbStatus,
	// 	"redis":    redisStatus,}

	// Fixed: Initialize as a map using the := operator
    // data := map[string]string{
    //     "database": dbStatus,
    //     "redis":    redisStatus,
    // }
    
    // Alternatively, if utils.Success accepts an interface{}, you can use gin.H:
    data := gin.H{
        "database": dbStatus,
        "redis":    redisStatus,
    }
	// utils.Success(c, "Application is healthy", data)
	// c.JSON(http.StatusOK, gin.H{
	// 	"success":  true,
	// 	"message":  "Application is healthy",
	// 	"database": dbStatus,
	// 	"redis":    redisStatus,
	// })
	c.JSON(200, gin.H{
		"success":  true,
		"message":  "Application is healthy",
		"data":data,
	})
}