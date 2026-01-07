package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

	log.Printf("Message received from %s in server %s", m.Author.Username, m.GuildID)

	// Show typing indicator
	err := s.ChannelTyping(m.ChannelID)
	if err != nil {
		log.Printf("Error sending typing indicator: %v", err)
	}

	// Get conversation context
	messages, err := c.buildConversationContext(s, m)
	if err != nil {
		log.Printf("Error building conversation context: %v", err)
		return
	}

	// Get response from Mistral AI
	response, err := c.messageGeneratorHandler(messages)
	if err != nil {
		log.Printf("Error getting response from AI: %v", err)
		return
	}

	// Send the response
	_, err = s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
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

func (c *DiscordClient) cleanMessageContent(msg *discordgo.Message) string {
	// Replace bot mentions
	cleaned := strings.ReplaceAll(msg.Content, fmt.Sprintf("<@%s>", c.botID), "@Assistant")

	// Replace user mentions by @username
	for _, user := range msg.Mentions {
		cleaned = strings.ReplaceAll(cleaned, fmt.Sprintf("<@%s>", user.ID), "@"+user.DisplayName())
	}

	// Replace role mentions by @role
	for _, roleID := range msg.MentionRoles {
		role, err := c.session.State.Role(msg.GuildID, roleID)
		if err == nil {
			cleaned = strings.ReplaceAll(cleaned, fmt.Sprintf("<@&%s>", roleID), "@"+role.Name)
		}
	}

	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
