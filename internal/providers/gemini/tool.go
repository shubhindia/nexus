package gemini

import (
	"encoding/json"

	"github.com/shubhindia/nexus/internal/types"
)

var allowedTools = map[string]struct{}{
	"exec_command":                {},
	"write_stdin":                 {},
	"list_mcp_resources":          {},
	"list_mcp_resource_templates": {},
	"read_mcp_resource":           {},
	"update_plan":                 {},
}

func convertTools(tools []types.Tool) []geminiTool {
	if len(tools) == 0 {
		return nil
	}

	declarations := make([]geminiFunctionDeclaration, 0, len(tools))

	for _, tool := range tools {
		if _, ok := allowedTools[tool.Name]; !ok {
			continue
		}

		declarations = append(declarations, geminiFunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  sanitizeToolParameters(tool.Parameters),
		})
	}

	if len(declarations) == 0 {
		return nil
	}

	return []geminiTool{{
		FunctionDeclarations: declarations,
	}}
}

func convertToolChoice(choice *types.ToolChoice) *geminiToolConfig {
	if choice == nil {
		return nil
	}

	switch choice.Mode {
	case types.ToolChoiceAuto:
		return &geminiToolConfig{
			FunctionCallingConfig: &geminiFunctionCallingConfig{
				Mode: "AUTO",
			},
		}
	case types.ToolChoiceNone:
		return &geminiToolConfig{
			FunctionCallingConfig: &geminiFunctionCallingConfig{
				Mode: "NONE",
			},
		}
	case types.ToolChoiceTool:
		cfg := &geminiToolConfig{
			FunctionCallingConfig: &geminiFunctionCallingConfig{
				Mode: "ANY",
			},
		}

		if choice.Name != "" {
			cfg.FunctionCallingConfig.AllowedFunctionNames = []string{choice.Name}
		}

		return cfg
	}

	return nil
}

func sanitizeToolParameters(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}

	var schema any
	if err := json.Unmarshal(raw, &schema); err != nil {
		return raw
	}

	sanitizeSchema(schema)

	cleaned, err := json.Marshal(schema)
	if err != nil {
		return raw
	}

	return cleaned
}

func sanitizeSchema(node any) {
	switch v := node.(type) {
	case map[string]any:
		delete(v, "additionalProperties")

		for _, child := range v {
			sanitizeSchema(child)
		}

	case []any:
		for _, child := range v {
			sanitizeSchema(child)
		}
	}
}
