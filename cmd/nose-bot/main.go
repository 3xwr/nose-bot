package main

import (
	"fmt"
	"log"
	"nose-bot/internal/pkg/client"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func main() {
	setupConfig()
	dg, err := discordgo.New("Bot " + viper.GetString("discord-token"))
	if err != nil {
		log.Fatalf("Cannot start the session: %v", err)
	}
	defer dg.Close()
	discordClient := client.NewDiscordClient(dg)

	_ = discordClient
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}

func setupConfig() {
	viper.SetConfigName("local")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
