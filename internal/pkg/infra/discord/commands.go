package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	CommandImagen,
}

var DeprecatedCommands = []*discordgo.ApplicationCommand{
	CommandImagenTxt,
}

var CommandImagen = &discordgo.ApplicationCommand{
	Name: "imagen",
	Type: discordgo.MessageApplicationCommand,
}

var CommandImagenTxt = &discordgo.ApplicationCommand{
	Name: "imagen-txt",
	Type: discordgo.MessageApplicationCommand,
}
