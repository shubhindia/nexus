package config

import "os"

const (
	ProviderLlamaCPP = "llama.cpp"
	ProviderGemini   = "gemini"

	defaultNexusPort = ":9000"
	defaultProvider  = ProviderLlamaCPP
	defaultLlamaURL  = "http://192.168.1.11:8080"
	defaultGeminiURL = "https://generativelanguage.googleapis.com/v1beta"
)

type Config struct {
	NexusPort string
	Provider  string
	LlamaURL  string
	GeminiURL string
	GeminiKey string
}

func Load() Config {
	return Config{
		NexusPort: envOrDefault(
			"NEXUS_PORT",
			defaultNexusPort,
		),
		Provider: envOrDefault(
			"NEXUS_PROVIDER",
			defaultProvider,
		),
		LlamaURL: envOrDefault(
			"NEXUS_LLAMACPP_URL",
			defaultLlamaURL,
		),
		GeminiURL: envOrDefault(
			"NEXUS_GEMINI_URL",
			defaultGeminiURL,
		),
		GeminiKey: os.Getenv("NEXUS_GEMINI_API_KEY"),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
