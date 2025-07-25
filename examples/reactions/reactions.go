package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	debug := flag.Bool("debug", false, "Show JSON output")
	flag.Parse()

	// Get token from environment variable
	apiToken := os.Getenv("SLACK_BOT_TOKEN")
	if apiToken == "" {
		fmt.Println("SLACK_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	api := slack.New(apiToken, slack.OptionDebug(*debug))

	var (
		postAsUserName  string
		postAsUserID    string
		postToUserName  string
		postToUserID    string
		postToChannelID string
	)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error getting channels: %s\n", err)
		return
	}

	// Post as the authenticated user.
	postAsUserName = authTest.User
	postAsUserID = authTest.UserID

	// Posting to DM with self causes a conversation with slackbot.
	postToUserName = authTest.User
	postToUserID = authTest.UserID

	// Find the channel.
	channel, _, _, err := api.OpenConversation(&slack.OpenConversationParameters{ChannelID: postToUserID})
	if err != nil {
		fmt.Printf("Error opening IM: %s\n", err)
		return
	}
	postToChannelID = channel.ID

	fmt.Printf("Posting as %s (%s) in DM with %s (%s), channel %s\n", postAsUserName, postAsUserID, postToUserName, postToUserID, postToChannelID)

	// Post a message.
	channelID, timestamp, err := api.PostMessage(postToChannelID, slack.MsgOptionText("Is this any good?", false))
	if err != nil {
		fmt.Printf("Error posting message: %s\n", err)
		return
	}

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// React with :+1:
	if err = api.AddReaction("+1", msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}

	// React with :-1:
	if err = api.AddReaction("cry", msgRef); err != nil {
		fmt.Printf("Error adding reaction: %s\n", err)
		return
	}

	// Get all reactions on the message.
	msgReactions, err := api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil {
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message...\n", len(msgReactions))
	for _, r := range msgReactions {
		fmt.Printf("  %d users say %s\n", r.Count, r.Name)
	}

	// List all of the users reactions.
	listReactions, _, err := api.ListReactions(slack.NewListReactionsParameters())
	if err != nil {
		fmt.Printf("Error listing reactions: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("All reactions by %s...\n", authTest.User)
	for _, item := range listReactions {
		fmt.Printf("%d on a %s...\n", len(item.Reactions), item.Type)
		for _, r := range item.Reactions {
			fmt.Printf("  %s (along with %d others)\n", r.Name, r.Count-1)
		}
	}

	// Remove the :cry: reaction.
	err = api.RemoveReaction("cry", msgRef)
	if err != nil {
		fmt.Printf("Error remove reaction: %s\n", err)
		return
	}

	// Get all reactions on the message.
	msgReactions, err = api.GetReactions(msgRef, slack.NewGetReactionsParameters())
	if err != nil {
		fmt.Printf("Error getting reactions: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("%d reactions to message after removing cry...\n", len(msgReactions))
	for _, r := range msgReactions {
		fmt.Printf("  %d users say %s\n", r.Count, r.Name)
	}
}
