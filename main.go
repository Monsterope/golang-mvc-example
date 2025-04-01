package main

import (
	"fmt"
	"monsterloveshop/app"
	"monsterloveshop/config"
	"monsterloveshop/controllers"
	"monsterloveshop/databases"
	"monsterloveshop/middleware"
	"monsterloveshop/routes"
	"monsterloveshop/store"
)

func init() {
	config.Load()
	// if os.Getenv("APP_ENV") != "production" {
	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal("Error loading .env file")
	// 	}
	// }
}

func main() {
	config.Load()

	dbConfig := &databases.DatabaseConfig{}
	dbConfig.Connect()
	// dbConfig.AutoMigrate()

	authStoreInstance := store.NewRedisAuthStore(config.GetEnv("redis.dns"))
	if authStoreInstance == nil {
		fmt.Println("Failed for connect Redis Auth")
	}
	controller := controllers.NewController(dbConfig, authStoreInstance)
	middleware := middleware.NewMiddlewareAuthRedis(authStoreInstance)

	app := app.NewApp()
	routes := routes.NewRoute(app.App, controller, middleware)
	routes.RouteApi()
	routes.RouteBo()

	app.Start(":8080")
}
