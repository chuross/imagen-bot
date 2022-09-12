package discord

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	CommandImagen,
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
