package main

import (
	"fmt"
	"hotel_app/config"
	"hotel_app/db"
	"log"
	"github.com/gin-gonic/gin"
	"hotel_app/handlers"
	"hotel_app/service"
	"hotel_app/repository"
)

func main() {
	cfg := config.Load()
	r := gin.Default()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate(database)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected and migrated!")
    hotelRepo := repository.NewHotelRepository(database)
	hotelService := service.NewHotelService(hotelRepo)
	hotelHandler := handlers.NewHotelHandler(hotelService)
	 api := r.Group("/api/v1")
    {
        hotels := api.Group("/hotels")
        {
            hotels.GET("/GetAllHotels", hotelHandler.GetHotels)
            hotels.GET("/GetHotelByID/:id", hotelHandler.GetHotelByID)
			hotels.POST("/CreateHotel", hotelHandler.CreateHotel)
			hotels.PUT("/UpdateHotel/:id", hotelHandler.UpdateHotel)	
			hotels.DELETE("/DeleteHotel/:id", hotelHandler.DeleteHotel)
        }
    }

    log.Printf("Server running on port %s", cfg.Port)
    r.Run(":" + cfg.Port)

}