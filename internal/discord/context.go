package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const maxHistoryDepth = 10

// buildConversationContext builds the conversation history for the AI
// recursively fetches the reply chain up to a maximum depth of 10 messages
func (c *DiscordClient) buildConversationContext(s *discordgo.Session, m *discordgo.MessageCreate) ([]Message, error) {
	var messages []Message

	// If this is a reply, recursively fetch the reply chain
	if m.ReferencedMessage != nil {
		history, err := c.fetchReplyChain(s, m.ReferencedMessage, maxHistoryDepth)
		if err != nil {
			log.Printf("Warning: failed to fetch full reply chain: %v", err)
			// Continue with partial history rather than failing completely
		}
		messages = append(messages, history...)
	}

	// Add the current message
	messages = append(messages, Message{
		Role:    "user",
		Content: c.cleanMessageContent(m.Message),
	})

	return messages, nil
}

func (c *DiscordClient) fetchReplyChain(s *discordgo.Session, firstMessage *discordgo.Message, depth int) ([]Message, error) {
	var messages []Message
	var currentMsg discordgo.Message

	// clone firstMessage into currentMsg
	currentMsg = *firstMessage

	for depth > 0 {
		if currentMsg.ReferencedMessage != nil {
			currentMsg = *currentMsg.ReferencedMessage
		} else if currentMsg.MessageReference != nil && currentMsg.MessageReference.MessageID != "" {
			parentMsg, err := s.ChannelMessage(currentMsg.MessageReference.ChannelID, currentMsg.MessageReference.MessageID)
			if err != nil {
				log.Printf("Warning: failed to fetch referenced message: %v", err)
				break
			}
			currentMsg = *parentMsg
		} else {
			break
		}

		var role string
		if currentMsg.Author.ID == c.botID {
			role = "assistant"
		} else {
			role = "user"
		}

		messages = append([]Message{{
			Role:    role,
			Content: c.cleanMessageContent(&currentMsg),
		}}, messages...)

		depth--
	}

	return messages, nil
}
