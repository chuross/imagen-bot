package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	CommandImagen,
}

var DeprecatedCommands = []*discordgo.ApplicationCommand{
	CommandImagenUpscaling,
}

var CommandImagen = &discordgo.ApplicationCommand{
	Name: "imagen",
	Type: discordgo.MessageApplicationCommand,
}

var CommandImagenUpscaling = &discordgo.ApplicationCommand{
	Name: "imagen-upscaling",
	Type: discordgo.MessageApplicationCommand,
}
