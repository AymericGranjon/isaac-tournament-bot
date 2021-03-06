package main

import (
	"database/sql"

	"github.com/Zamiell/isaac-tournament-bot/src/models"
	"github.com/bwmarrin/discordgo"
)

func commandTimeDelete(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database)
	var race models.Race
	if v, err := raceGet(m.ChannelID); err == sql.ErrNoRows {
		discordSend(m.ChannelID, "You can only use that command in a race channel.")
		return
	} else if err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can reschedule this match.")
		return
	}

	// Check to see if this race has already been scheduled
	if race.State == "initial" {
		discordSend(m.ChannelID, "There is no need to rescheulde until both racers have already agreed to a time.")
		return
	}

	// Set the scheduled time to null
	if err := db.Races.UnsetDatetimeScheduled(m.ChannelID); err != nil {
		msg := "Failed to unset the scheduled time: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the state
	race.State = "initial"
	if err := db.Races.SetState(m.ChannelID, race.State); err != nil {
		msg := "Failed to set the state: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	discordSend(m.ChannelID, "The currently scheduled time has been deleted. Please suggest a new time with the `!time` command.")
	log.Info("Race \"" + race.Name() + "\" rescheduled; state set to 0.")
}
