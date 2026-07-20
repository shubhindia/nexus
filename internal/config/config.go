package config

import "os"

const (
	ProviderLlamaCPP = "llama.cpp"

	defaultNexusPort = ":9000"
	defaultProvider  = ProviderLlamaCPP
	defaultLlamaURL  = "http://192.168.1.11:8080"
)

type Config struct {
	NexusPort string
	Provider  string
	LlamaURL  string
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
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
