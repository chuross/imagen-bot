package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	CommandImagen,
}

var DeprecatedCommands = []*discordgo.ApplicationCommand{}

var CommandImagen = &discordgo.ApplicationCommand{
	Name: "imagen",
	Type: discordgo.MessageApplicationCommand,
}

var CommandWorkspace = &discordgo.ApplicationCommand{
	Name: "workspace",
	Type: discordgo.MessageApplicationCommand,
}
