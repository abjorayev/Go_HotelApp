package main

import (
	"fmt"
	"hotel_app/config"
	"hotel_app/db"
	"hotel_app/handlers"
	"hotel_app/middleware"
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

	err = db.Seed(database)
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
	//Bookings
	bookingRepo := repository.NewBookingRepository(database)
	bookingService := service.NewBookingService(bookingRepo, roomRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	//Auth
	authRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)
	auth := middleware.AuthMiddleware()
	admin := middleware.RequireRole("admin")
	api := r.Group("/api/v1")
	{
		hotels := api.Group("/hotels")
		{
			hotels.GET("/GetAllHotels", hotelHandler.GetHotels)
			hotels.GET("/GetHotelByID/:id", hotelHandler.GetHotelByID)
			hotels.POST("/CreateHotel", auth, admin, hotelHandler.CreateHotel)
			hotels.PUT("/UpdateHotel/:id", auth, admin, hotelHandler.UpdateHotel)
			hotels.DELETE("/DeleteHotel/:id", auth, admin, hotelHandler.DeleteHotel)
		}
		rooms := api.Group("/rooms")
		{
			rooms.GET("/GetRoomsByHotelID/:id", roomHandler.GetByHotelId)
			rooms.GET("/GetRoomsByID/:id", roomHandler.GetRoomById)
			rooms.POST("/CreateRoom", auth, admin, roomHandler.CreateRoom)
			rooms.PUT("/UpdateRoom", auth, admin, roomHandler.UpdateRoom)
			rooms.DELETE("/DeleteRoom/:id", auth, admin, roomHandler.DeleteRoom)
		}
		bookings := api.Group("/bookings")
		{
			bookings.GET("/GetBookingById/:id", bookingHandler.GetByBookingId)
			bookings.GET("/GetBookingsByRoomId/:id", bookingHandler.GetByRoomid)
			bookings.POST("/CreateBooking", auth, bookingHandler.CreateBooking)
			bookings.DELETE("/DeleteBooking", auth, bookingHandler.DeleteBooking)
		}
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	log.Printf("Server running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)

}
