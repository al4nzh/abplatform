package config

type Config struct {
	HTTPPort    string
	PostgresDSN string
	KafkaBroker string
	KafkaTopic  string
}