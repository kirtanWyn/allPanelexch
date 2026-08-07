package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/kirtanwyn/allPanelexch/config"
	"github.com/kirtanwyn/allPanelexch/routes"
)

func main() {

	// err := godotenv.Load()
	err := godotenv.Load("./.env")

	if err != nil {
		fmt.Println("error in load dotenv")
		panic(err)
	}

	// if err := config.LoadConfig(); err != nil {
	// 	log.Fatal(err)
	// }

	if err := config.InitDatabase(); err != nil {
		log.Fatal(err)
	}

	if err := config.InitRedis(); err != nil {
		log.Fatal(err)
	}

	defer config.CloseDatabase()

	log.Println("Server Started...")

	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

routes.SetupRouter(router)
// routes.SetupAccessControlRouter(router)

	log.Println("🚀 Server Started on Port 9000")

	// router.Run(":9000")
	router.Run(":" + os.Getenv("PORT"))
}
