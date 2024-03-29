package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/berat703/mediago/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Server struct
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize Server
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Image{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run server
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
