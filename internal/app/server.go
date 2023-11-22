package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DefaultDatabaseURL = os.Getenv("DEFAULT_DATABASE_URL")

func StartServer() {
	g := gin.Default()

	db, err := gorm.Open(postgres.Open(DefaultDatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&PromoModel{}); err != nil {
		log.Fatal(err)
	}

	promoController := PromoController{db}

	g.GET("/promos", promoController.GetPromos)
	g.GET("/promos/:id", promoController.GetPromo)
	g.POST("/promos", promoController.CreatePromo)
	g.PUT("/promos/:id", promoController.UpdatePromo)

	g.Run()
}
