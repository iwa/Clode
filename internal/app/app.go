package app

import (
	"github.com/iwa/Clode/internal/ai"
	"github.com/iwa/Clode/internal/discord"
)

type App struct {
	discordClient discord.DiscordClient
	aiClient      ai.AIClient
}

func NewApp() *App {
	return &App{}
}
