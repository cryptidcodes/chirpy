package main
import _ "github.com/lib/pq"
import (
	"github.com/cryptidcodes/chirpy/internal/database"
	"sync/atomic"
	"fmt"
	"net/http"
	"os"
	"database/sql"
	"log"
	"github.com/joho/godotenv"
)

type apiConfig struct {
		fileserverHits atomic.Int32
		dbQueries *database.Queries
		platform string
	}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// load environment variables from .env file
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("Error connecting to the database: ", err)
	}

	dbQueries := database.New(db)

	// init config struct
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: dbQueries,
		platform: platform,
	}

	// create a new http.ServeMux to handle requests
	mux := http.NewServeMux()

	// handle /app/ with middleware to count hits
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// additional endpoint handlers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)

	// admin endpoint handlers
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	
	// create a new http.Server struct
	server := &http.Server{
		Addr:	":" + port,
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	// start the server
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}