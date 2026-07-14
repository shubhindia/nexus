package compiler

import (
	"strings"

	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/types"
)

type OpenAICompiler struct{}

func NewOpenAICompiler() *OpenAICompiler {
	return &OpenAICompiler{}
}

func (c *OpenAICompiler) Compile(
	req *openai.ResponsesRequest,
) *types.ChatRequest {

	return &types.ChatRequest{
		Model: req.Model,

		Messages: c.messages(req),

		Tools: c.tools(req),

		ToolChoice: c.toolChoice(req),

		Temperature: nil,
		TopP:        nil,
		MaxTokens:   nil,

		Stream: req.Stream,
	}
}

func (c *OpenAICompiler) messages(
	req *openai.ResponsesRequest,
) []types.Message {

	messages := make([]types.Message, 0)

	if system := c.compileInstructions(req.Instructions); system != "" {
		messages = append(messages, types.Message{
			Role:    types.RoleSystem,
			Content: system,
		})
	}

	messages = append(messages, c.compileConversation(req.Input)...)

	return messages
}

func (c *OpenAICompiler) compileInstructions(
	instructions string,
) string {

	if instructions == "" {
		return ""
	}

	instructions = stripRuntimeSections(instructions)

	return strings.TrimSpace(instructions)
}

func (c *OpenAICompiler) compileConversation(
	input []openai.ResponseInput,
) []types.Message {

	messages := make([]types.Message, 0)

	for _, message := range input {

		text := inputText(message.Content)
		if text == "" {
			continue
		}

		switch types.ParseRole(message.Role) {

		case types.RoleDeveloper:
			fallthrough

		case types.RoleSystem:

			system := c.compileInstructions(text)
			if system == "" {
				continue
			}

			messages = append(messages, types.Message{
				Role:    types.RoleSystem,
				Content: system,
			})

		case types.RoleUser:

			messages = append(messages, types.Message{
				Role:    types.RoleUser,
				Content: text,
			})

		case types.RoleAssistant:

			messages = append(messages, types.Message{
				Role:    types.RoleAssistant,
				Content: text,
			})
		}
	}

	return messages
}

func (c *OpenAICompiler) tools(
	req *openai.ResponsesRequest,
) []types.Tool {

	return nil
}

func (c *OpenAICompiler) toolChoice(
	req *openai.ResponsesRequest,
) *types.ToolChoice {

	return nil
}

func inputText(
	content []openai.ResponseInputContent,
) string {

	var text strings.Builder

	for _, item := range content {

		if item.Type != "input_text" {
			continue
		}

		if text.Len() > 0 {
			text.WriteByte('\n')
		}

		text.WriteString(item.Text)
	}

	return text.String()
}

func stripRuntimeSections(
	instructions string,
) string {

	sections := []string{
		"permissions instructions",
		"environment_context",
		"skills_instructions",
		"plugins_instructions",
	}

	for _, section := range sections {
		instructions = stripSection(
			instructions,
			section,
		)
	}

	return instructions
}

func stripSection(
	text string,
	section string,
) string {

	startTag := "<" + section + ">"
	endTag := "</" + section + ">"

	for {

		start := strings.Index(
			text,
			startTag,
		)
		if start == -1 {
			break
		}

		end := strings.Index(
			text[start:],
			endTag,
		)
		if end == -1 {
			break
		}

		end += start + len(endTag)

		text = text[:start] + text[end:]
	}

	return text
}
