package main

import (
	"fmt"
	"hotel_app/config"
	"hotel_app/db"
	"hotel_app/handlers"
	"hotel_app/repository"
	"hotel_app/service"
	"log"

	"github.com/gin-gonic/gin"
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
	//Hotels
	hotelRepo := repository.NewHotelRepository(database)
	hotelService := service.NewHotelService(hotelRepo)
	hotelHandler := handlers.NewHotelHandler(hotelService)
	//Rooms
	roomRepo := repository.NewRoomRepository(database)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handlers.NewRoomHandler(roomService)
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
		rooms := api.Group("/rooms")
		{
			rooms.GET("/GetRoomsByHotelID/:id", roomHandler.GetByHotelId)
			rooms.GET("/GetRoomsByID/:id", roomHandler.GetRoomById)
			rooms.POST("/CreateRoom", roomHandler.CreateRoom)
			rooms.PUT("/UpdateRoom", roomHandler.UpdateRoom)
			rooms.DELETE("/DeleteRoom/:id", roomHandler.DeleteRoom)
		}
	}

	log.Printf("Server running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)

}
