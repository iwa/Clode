package app

import (
	"github.com/iwa/Clode/internal/ai"
	"github.com/iwa/Clode/internal/discord"
)

type App struct {
	DiscordClient discord.DiscordClient
	AIClient      ai.AIClient
}

func NewApp(discordToken, aiKey, agentID string) (*App, error) {
	aiClient, err := ai.NewAIClient(aiKey, agentID, "")
	if err != nil {
		return nil, err
	}

	discordClient, err := discord.NewDiscordClient(discordToken, aiClient.Chat)
	if err != nil {
		return nil, err
	}

	return &App{
		DiscordClient: *discordClient,
		AIClient:      *aiClient,
	}, err
}
