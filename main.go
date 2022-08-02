package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

var (
	token = os.Getenv("token")

	commands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "say",
			Description: "says what you say",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "ephemeral",
					Description: "If the response should only be visible to you",
					Required:    true,
				},
			},
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	gatewayLogger := log.New(log.LstdFlags | log.Lshortfile)
	gatewayLogger.SetLevel(log.LevelDebug)

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(gateway.WithLogger(gatewayLogger)),
		bot.WithEventListenerFunc(commandListener),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest().SetGlobalCommands(client.ApplicationID(), commands); err != nil {
		log.Fatal("error while registering commands: ", err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		)
		if err != nil {
			event.Client().Logger().Error("error on sending response: ", err)
		}
	}
}
