package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ignacio-J-Maylin/arithmetic-calculator/config"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/handlers/authHandlers"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/handlers/userHandlers"
	"github.com/Ignacio-J-Maylin/arithmetic-calculator/middlewares"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
		}
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error conectándose a la base de datos: %v", err)
	}
	defer db.Close()

	log.Println("Conexión exitosa a la base de datos")

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	log.Println("Migraciones ejecutadas correctamente.")

	mux := http.NewServeMux()

	mux.Handle("/api/v1/users/credits", middlewares.AuthMiddleware(http.HandlerFunc(userHandlers.HandleCredits(db))))
	mux.Handle("/api/v1/users/operation", middlewares.AuthMiddleware(http.HandlerFunc(userHandlers.PerformOperation(db))))
	mux.Handle("/api/v1/records/history", middlewares.AuthMiddleware(http.HandlerFunc(userHandlers.GetRecordsHistory(db))))
	mux.Handle("/api/v1/records/delete", middlewares.AuthMiddleware(http.HandlerFunc(userHandlers.DeleteRecordHandler(db))))

	mux.HandleFunc("/api/v1/logout", http.HandlerFunc(authHandlers.Logout()))
	mux.HandleFunc("/api/v1/login", authHandlers.Login(db))
	mux.HandleFunc("/api/v1/refresh", authHandlers.RefreshToken(db))
	mux.HandleFunc("/api/v1/signup", authHandlers.SignUp(db))
	mux.HandleFunc("/api/v1/operations", userHandlers.GetOperations(db))

	corsMux := middlewares.CorsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMux))
}
