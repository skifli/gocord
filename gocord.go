package gocord

import (
	"errors"

	"github.com/skifli/gocord/api"
)

// Client represents a Discord self-bot.
type Client struct {
	Gateway *api.Gateway // Gateway contains data relating to a Discord WebSocket connection.
	SelfBot *api.SelfBot // SelfBot contains data relating to the self-bot.
}

// Creates a SelfBot struct, used when initializing a Client struct.
func SelfBotToken(token string) *api.SelfBot {
	return &api.SelfBot{
		Token: token,
	}
}

func (client *Client) Init() error {
	if client.SelfBot == nil {
		return errors.New("self-bot data not supplied")
	} else if client.SelfBot.Token == "" {
		return errors.New("self-bot's token not supplied")
	}

	client.Gateway = api.CreateGateway(client.SelfBot)

	return nil
}

func (client *Client) Connect() error {
	return client.Gateway.Connect()
}

func (client *Client) AddEventHandler(event api.GatewayEventName, function any) error {
	return client.Gateway.Handlers.Add(event, function)
}
