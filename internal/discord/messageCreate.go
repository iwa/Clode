package discord

import "github.com/bwmarrin/discordgo"

// handler for incoming messages
func (c *DiscordClient) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore itself
	if m.Author.ID == c.botID {
		return
	}

	// check if the bot is mentioned and should respond
	shouldRespond := c.shouldRespondToMessage(m)
	if !shouldRespond {
		return
	}

	// TODO: add all the logic to read messages and generate response...
}

// shouldRespondToMessage checks if the bot should respond to the given message
// respond if bot is mentioned or if someone replies to a bot's message
func (c *DiscordClient) shouldRespondToMessage(m *discordgo.MessageCreate) bool {
	// Check if the bot is mentioned
	for _, mention := range m.Mentions {
		if mention.ID == c.botID {
			return true
		}
	}

	// Check if the message is a reply to the bot
	if m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID == c.botID {
		return true
	}

	return false
}
