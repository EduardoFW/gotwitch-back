package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
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
}

var ServerSettings = &Server{}

func getParam(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Setup() {
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

	println("Loading settings...")
	ServerSettings.Addr = getParam("ADDR", ":8080")
	ServerSettings.ReadTimeout, _ = strconv.Atoi(getParam("READ_TIMEOUT", "10"))
	ServerSettings.WriteTimeout, _ = strconv.Atoi(getParam("WRITE_TIMEOUT", "10"))
	ServerSettings.IdleTimeout, _ = strconv.Atoi(getParam("IDLE_TIMEOUT", "60"))

	ServerSettings.TwitchClientID = getParam("TWITCH_CLIENT_ID", "")
	ServerSettings.TwitchClientSecret = getParam("TWITCH_CLIENT_SECRET", "")

	println("Settings loaded!")
}
