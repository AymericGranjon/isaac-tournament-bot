package main

import (
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	discord          *discordgo.Session
	discordBotID     string
	discordGuildName string
	discordGuildID   string
	commandMutex     = new(sync.Mutex)
)

func discordInit() {
	// Read the OAuth secret from the environment variable
	discordToken := os.Getenv("DISCORD_TOKEN")
	if len(discordToken) == 0 {
		log.Fatal("The \"DISCORD_TOKEN\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	discordGuildName = os.Getenv("DISCORD_SERVER_NAME")
	if len(discordGuildName) == 0 {
		log.Fatal("The \"DISCORD_SERVER_NAME\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	// Bot accounts must be prefixed with "Bot"
	if d, err := discordgo.New("Bot " + discordToken); err != nil {
		log.Fatal("Failed to create a Discord session:", err)
		return
	} else {
		discord = d
	}

	// Register function handlers for various events
	discord.AddHandler(discordReady)
	discord.AddHandler(discordMessageCreate)

	// Register function handlers for all of the commands
	commandInit()

	// Open the websocket and begin listening
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening Discord session: ", err)
		return
	}
}

/*
	Event handlers
*/

func discordReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Info("Discord bot connected with username: " + event.User.Username)
	discordBotID = event.User.ID

	// Get the guild ID
	var guilds []*discordgo.UserGuild
	if v, err := s.UserGuilds(1, "", ""); err != nil {
		log.Fatal("Failed to get the Discord guilds:", err)
		return
	} else {
		guilds = v
	}

	foundGuild := false
	for _, guild := range guilds {
		if guild.Name == discordGuildName {
			foundGuild = true
			discordGuildID = guild.ID
			break
		}
	}
	if !foundGuild {
		log.Fatal("Failed to find the ID of the \"" + discordGuildName + "\" Discord server.")
	}
}

func discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == discordBotID {
		return
	}

	// Log the message
	var channelName string
	if v, err := discord.Channel(m.ChannelID); err != nil {
		log.Error("Failed to get the channel name for the channel ID of \""+m.ChannelID+"\":", err)
	} else {
		channelName = v.Name
	}
	log.Info("[#" + channelName + "] <" + m.Author.Username + "#" + m.Author.Discriminator + "> " + m.Content)

	// Commands for this bot will start with a "!", so we can ignore everything else
	message := strings.ToLower(m.Content)
	args := strings.SplitN(message, " ", 2) // We use SplitN because there
	command := args[0]
	args = args[1:] // This will be an empty slice if there is nothing after the command
	if !strings.HasPrefix(command, "!") {
		return
	}
	command = strings.TrimPrefix(command, "!")

	// Check to see if there is a command handler for this command
	if _, ok := commandHandlerMap[command]; !ok {
		discordSend(m.ChannelID, "That is not a valid command.")
		return
	}

	commandMutex.Lock()
	commandHandlerMap[command](m, args)
	commandMutex.Unlock()
}

/*
	Miscellaneous functions
*/

func discordSend(channelID string, message string) {
	if _, err := discord.ChannelMessageSend(channelID, message); err != nil {
		log.Error("Failed to send \"" + message + "\" to \"" + channelID + "\": " + err.Error())
		return
	}
}