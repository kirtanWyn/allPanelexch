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

	// router := routes.SetupRouter(),routes.SetupAcessControllRouter()

router := gin.New()

router.Use(gin.Recovery())

routes.SetupRouter(router)
// routes.SetupAccessControlRouter(router)

	log.Println("🚀 Server Started on Port 9000")

	// router.Run(":9000")
	router.Run(":" + os.Getenv("PORT"))
}
