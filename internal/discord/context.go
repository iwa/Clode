package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/iwa/Clode/internal/ai"
)

const (
	maxHistoryDepth = 10
)

// buildConversationContext builds the conversation history for the AI
// recursively fetches the reply chain up to a maximum depth of 10 messages
func (c *DiscordClient) buildConversationContext(s *discordgo.Session, m *discordgo.MessageCreate) ([]ai.AIMessage, error) {
	var messages []ai.AIMessage

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
	messages = append(messages, ai.AIMessage{
		Role:    "user",
		Content: fmt.Sprintf("%s: %s", m.Author.DisplayName(), c.cleanMessageContent(m.Message)),
	})

	return messages, nil
}

func (c *DiscordClient) fetchReplyChain(s *discordgo.Session, firstMessage *discordgo.Message, depth int) ([]ai.AIMessage, error) {
	var messages []ai.AIMessage
	var currentMsg discordgo.Message

	// clone firstMessage into currentMsg
	currentMsg = *firstMessage

	for depth > 0 {
		// Add message to context
		var role, content string
		if currentMsg.Author.ID == c.botID {
			role = "assistant"
			content = c.cleanMessageContent(&currentMsg)
		} else {
			role = "user"
			content = fmt.Sprintf("%s: %s", currentMsg.Author.DisplayName(), c.cleanMessageContent(&currentMsg))
		}

		messages = append([]ai.AIMessage{{
			Role:    role,
			Content: content,
		}}, messages...)

		// Determine the parent message, if any
		// Break the loop if none found
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

		// Decrease depth
		depth--
	}

	return messages, nil
}
