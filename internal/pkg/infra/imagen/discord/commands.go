package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name: "imagen",
		Type: discordgo.MessageApplicationCommand,
		Description: "指定した情報で画像を生成する",
	},
}
