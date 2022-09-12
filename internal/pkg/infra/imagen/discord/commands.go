package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name: "imagen",
		Type: discordgo.MessageApplicationCommand,
	},
}
