package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/iwa/Clode/internal/ai"
)

type MessageGeneratorHandler func(messages []ai.AIMessage) (string, error)

type DiscordClient struct {
	session                 *discordgo.Session
	botID                   string
	messageGeneratorHandler MessageGeneratorHandler
}

func NewDiscordClient(discordToken string, handler MessageGeneratorHandler) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, err
	}

	client := &DiscordClient{
		session:                 session,
		messageGeneratorHandler: handler,
	}

	client.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	client.session.AddHandler(client.messageCreate)

	return client, nil
}

func (c *DiscordClient) Start() error {
	if err := c.session.Open(); err != nil {
		return err
	}

	user, err := c.session.User("@me")
	if err != nil {
		return err
	}
	c.botID = user.ID

	log.Printf("Bot is running as %s (ID: %s)", user.Username, user.ID)
	return nil
}

func (c *DiscordClient) Stop() error {
	return c.session.Close()
}
