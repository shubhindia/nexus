package prompt

import (
	"strings"

	"github.com/shubhindia/nexus/internal/types"
)

type CodexCompiler struct{}

func NewCodexCompiler() *CodexCompiler {
	return &CodexCompiler{}
}

func (c *CodexCompiler) Messages(
	req *types.ResponsesRequest,
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

func (c *CodexCompiler) compileInstructions(
	instructions string,
) string {

	if instructions == "" {
		return ""
	}

	instructions = stripRuntimeSections(instructions)

	return strings.TrimSpace(instructions)
}

func (c *CodexCompiler) compileConversation(
	input []types.ResponseInput,
) []types.Message {

	messages := make([]types.Message, 0)

	for _, message := range input {

		text := inputText(message.Content)
		if text == "" {
			continue
		}

		switch message.Role {

		case types.RoleDeveloper.String():
			fallthrough

		case types.RoleSystem.String():

			system := c.compileInstructions(text)
			if system == "" {
				continue
			}

			messages = append(messages, types.Message{
				Role:    types.RoleSystem,
				Content: system,
			})

		case "user":

			messages = append(messages, types.Message{
				Role:    types.RoleUser,
				Content: text,
			})

		case "assistant":

			messages = append(messages, types.Message{
				Role:    types.RoleAssistant,
				Content: text,
			})
		}
	}

	return messages
}

func (c *CodexCompiler) Tools(
	req *types.ResponsesRequest,
) []types.Tool {

	return nil
}

func (c *CodexCompiler) ToolChoice(
	req *types.ResponsesRequest,
) *types.ToolChoice {

	return nil
}

func inputText(
	content []types.ResponseInputContent,
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
