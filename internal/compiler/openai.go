package compiler

import (
	"encoding/json"
	"strings"

	"github.com/shubhindia/nexus/internal/schema/openai"
	"github.com/shubhindia/nexus/internal/toolstate"
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
		Model:       req.Model,
		Messages:    c.messages(req),
		Tools:       c.tools(req),
		ToolChoice:  c.toolChoice(req),
		Temperature: nil,
		TopP:        nil,
		MaxTokens:   nil,
		Stream:      req.Stream,
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
	callNames := map[string]string{}
	skippedCalls := map[string]string{}

	for _, message := range input {

		switch message.Type {
		case "function_call":
			call := compileFunctionCall(message)
			if call == nil {
				if message.CallID != "" && message.Name != "" {
					skippedCalls[message.CallID] = message.Name
				}
				continue
			}

			messages = append(messages, types.Message{
				Role:      types.RoleAssistant,
				ToolCalls: []types.ToolCall{*call},
			})

			if call.ID != "" && call.Name != "" {
				callNames[call.ID] = call.Name
			}

			continue

		case "function_call_output":
			toolName := message.Name
			if toolName == "" {
				toolName = callNames[message.CallID]
			}

			if _, skipped := skippedCalls[message.CallID]; skipped || message.CallID == "" {
				content := strings.TrimSpace(message.Output)
				if content == "" {
					continue
				}

				if toolName != "" {
					content = "Tool " + toolName + " output:\n" + content
				}

				messages = append(messages, types.Message{
					Role:    types.RoleUser,
					Content: content,
				})

				continue
			}

			messages = append(messages, types.Message{
				Role:       types.RoleTool,
				Name:       toolName,
				ToolCallID: message.CallID,
				Content:    strings.TrimSpace(message.Output),
			})

			continue
		}

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

func compileFunctionCall(
	input openai.ResponseInput,
) *types.ToolCall {
	if input.Name == "" {
		return nil
	}

	arguments := json.RawMessage(input.Arguments)
	if len(arguments) == 0 {
		arguments = json.RawMessage(`{}`)
	}

	if !json.Valid(arguments) {
		encoded, err := json.Marshal(input.Arguments)
		if err != nil {
			return nil
		}

		arguments = encoded
	}

	thoughtSignature := input.ThoughtSignature
	if thoughtSignature == "" {
		thoughtSignature = toolstate.LookupThoughtSignature(input.CallID)
	}

	if thoughtSignature == "" {
		return nil
	}

	return &types.ToolCall{
		ID:               input.CallID,
		Name:             input.Name,
		Arguments:        arguments,
		ThoughtSignature: thoughtSignature,
	}
}

func (c *OpenAICompiler) tools(
	req *openai.ResponsesRequest,
) []types.Tool {

	if len(req.Tools) == 0 {
		return nil
	}

	tools := make([]types.Tool, 0, len(req.Tools))

	for _, tool := range req.Tools {
		if tool.Type != "function" {
			continue
		}

		tools = append(tools, types.Tool{
			Type:        types.ToolTypeFunction,
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
		})
	}

	return tools
}
func (c *OpenAICompiler) toolChoice(
	req *openai.ResponsesRequest,
) *types.ToolChoice {

	if req.ToolChoice == nil {
		return nil
	}

	return &types.ToolChoice{
		Mode: toToolChoiceMode(req.ToolChoice.Mode),
		Name: req.ToolChoice.Name,
	}
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

func toToolChoiceMode(mode string) types.ToolChoiceMode {
	switch mode {

	case openai.ToolChoiceAuto:
		return types.ToolChoiceAuto

	case openai.ToolChoiceNone:
		return types.ToolChoiceNone

	case openai.ToolChoiceRequired:
		return types.ToolChoiceTool

	default:
		return types.ToolChoiceAuto
	}
}
