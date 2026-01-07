package app

import (
	"github.com/iwa/Clode/internal/ai"
	"github.com/iwa/Clode/internal/config"
	"github.com/iwa/Clode/internal/discord"
)

type App struct {
	DiscordClient *discord.DiscordClient
	AIClient      *ai.AIClient
}

func NewApp(config *config.Config) (*App, error) {
	aiClient, err := ai.NewAIClient(config.AIToken, string(config.AIMode), config.AIModel)
	if err != nil {
		return nil, err
	}

	discordClient, err := discord.NewDiscordClient(config.DiscordToken, aiClient.Chat)
	if err != nil {
		return nil, err
	}

	return &App{
		DiscordClient: discordClient,
		AIClient:      aiClient,
	}, err
}
