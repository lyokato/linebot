package config

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Web *WebConfig `toml:"web"`
	Bot *BotConfig `toml:"bot"`
}

type BotConfig struct {
	ChannelId                int    `toml:"channel_id"`
	ChannelSecret            string `toml:"channel_secret"`
	MID                      string `toml:"channel_mid"`
	ClientWorkerQueueSize    int    `toml:"client_worker_queue_size"`
	EventDispatcherQueueSize int    `toml:"event_dispatcher_queue_size"`
}

type WebConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func (wc *WebConfig) Address() string {
	return fmt.Sprintf("%s:%d", wc.Host, wc.Port)
}

func LoadFromFile(filePath string) *Config {
	var c Config
	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}
	return &c
}
