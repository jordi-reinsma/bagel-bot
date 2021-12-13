package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jordi-reinsma/bagel/core"
	"github.com/jordi-reinsma/bagel/db"
	"github.com/jordi-reinsma/bagel/slack"
)

var (
	RESET = os.Getenv("RESET") == "1"
	MOCK  = os.Getenv("MOCK") == "1"
)

func main() {
	fmt.Println("Starting bagel...")
	rand.Seed(time.Now().UnixNano())

	dB := db.MustConnect(RESET)
	defer dB.Close()

	skip, err := core.ShouldSkipExecution(dB)
	if err != nil {
		panic(err)
	}
	if skip {
		fmt.Println("Skipping execution")
		return
	}

	slack := slack.New(MOCK)

	channelIDs := slack.GetChannelUUIDs()

	for _, channelID := range channelIDs {
		fmt.Println("Generating matches for channel", channelID)

		userUUIDs, err := slack.GetUserUUIDs(channelID)
		if err != nil {
			panic(err)
		}

		users, err := dB.AddAndGetUsers(userUUIDs)
		if err != nil {
			panic(err)
		}

		fmt.Println("Users:", len(users))

		matches, err := core.GenerateMatches(dB, users)
		if err != nil {
			panic(err)
		}

		for _, match := range matches {
			err = slack.SendMessage(match, channelID)

			if err != nil {
				fmt.Println(match, err)
				continue
			}
			err = dB.UpdateMatch(match)
			if err != nil {
				fmt.Println(match, err)
				continue
			}
		}

		err = slack.SendChannelMessage(channelID)
		if err != nil {
			panic(err)
		}

		fmt.Println("Matches:", len(matches))
	}

	err = dB.SaveExecution()
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}
