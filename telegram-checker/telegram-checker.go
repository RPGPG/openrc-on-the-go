package main

import (
	"context"
	"fmt"
	"log"
	"openrc-on-the-go/config"
	"openrc-on-the-go/wrapper"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	services map[string]string
	pass     string
)

func getStatuses(services map[string]string) string {
	outs := []string{"Alias :: Name :: Status"}
	for alias, service := range services {
		started := wrapper.IsServiceStarted(service)
		if started {
			outs = append(outs, fmt.Sprintf("%s :: %s :: 🟢",
				alias, service))
		} else {
			outs = append(outs, fmt.Sprintf("%s :: %s :: 🔴",
				alias, service))
		}
	}
	return strings.Join(outs, "\n")
}

func main() {
	cmdArgs := os.Args[1:]
	configPath, err := config.GetConfigPath(cmdArgs)
	if err != nil {
		log.Fatal(err)
	}
	apiKey := ""
	services, pass, apiKey, err = config.LoadConfigTelegram(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(apiKey, opts...)
	if err != nil {
		log.Fatal(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == fmt.Sprintf("/openrc-check %s", pass) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   getStatuses(services),
		})
		log.Println("Status checked")
	}
}
