package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

// Config struct
type Server struct {
	// Port
	Addr string
	// ReadTimeout
	ReadTimeout int
	// WriteTimeout
	WriteTimeout int
	// IdleTimeout
	IdleTimeout int
	// Twitch Client ID
	TwitchClientID string
	// Twitch Client Secret
	TwitchClientSecret string
	// Postgres DSN
	PostgresDSN string
}

var ServerSettings = &Server{}
var DB *gorm.DB

func getParam(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Setup() {
	loadEnv()

	println("Loading settings...")
	ServerSettings.Addr = getParam("ADDR", ":8080")
	ServerSettings.ReadTimeout, _ = strconv.Atoi(getParam("READ_TIMEOUT", "10"))
	ServerSettings.WriteTimeout, _ = strconv.Atoi(getParam("WRITE_TIMEOUT", "10"))
	ServerSettings.IdleTimeout, _ = strconv.Atoi(getParam("IDLE_TIMEOUT", "60"))

	ServerSettings.TwitchClientID = getParam("TWITCH_CLIENT_ID", "")
	ServerSettings.TwitchClientSecret = getParam("TWITCH_CLIENT_SECRET", "")

	ServerSettings.PostgresDSN = "host=" + getParam("POSTGRES_HOST", "localhost") + " port=" + getParam("POSTGRES_PORT", "5432") + " user=" + getParam("POSTGRES_USER", "postgres") + " password=" + getParam("POSTGRES_PASSWORD", "") + " dbname=" + getParam("POSTGRES_DB", "twitch_go_backend")

	println("Settings loaded!")
}

func loadEnv() {
	println("Loading .env file...")
	env := os.Getenv("TWITCH_GO_BACKEND_ENV")
	println("Detected environment: " + env)
	switch env {
	case "devl":
		godotenv.Load(".env.devl")
		println("Loading .env.devl")
		break
	case "prod":
		godotenv.Load(".env.prod")
		println("Loading .env.prod")
		break
	default:
		godotenv.Load()
		println("Loading .env")
	}
}
