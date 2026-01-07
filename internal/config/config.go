package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string

	AIToken string
	AIMode  AIMode
	AIModel string
}

// AIMode enum
type AIMode string

const (
	MistralAPI   AIMode = "mistral-api"
	MistralAgent AIMode = "mistral-agent"
)

func NewConfig(discordToken, aiToken, aiMode, aiModel string) *Config {
	return &Config{
		DiscordToken: discordToken,
		AIToken:      aiToken,
		AIMode:       AIMode(aiMode),
		AIModel:      aiModel,
	}
}

func GenerateConfigFromEnv() *Config {
	// Parse .env file into runtime env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Local function
	getRequiredEnv := func(key string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		log.Fatalf("%s environment variable is required", key)
		return ""
	}

	// Import required env vars
	discordToken := getRequiredEnv("DISCORD_TOKEN")
	aiToken := getRequiredEnv("AI_TOKEN")
	aiModel := getRequiredEnv("AI_MODEL")

	aiMode := os.Getenv("AI_MODE")
	if aiMode == "" || (aiMode != string(MistralAPI) && aiMode != string(MistralAgent)) {
		log.Fatal("AI_MODE not set or invalid")
	}

	return NewConfig(discordToken, aiToken, aiMode, aiModel)
}
