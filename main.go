package main

import (
	"cfo/config"
	"cfo/repository/db/mysql"
	"cfo/service"
	"fmt"
	"github.com/h4ckitt/goTelegram"
	"log"
	"net/http"
	"strings"
)

var manager service.Service

func main() {

	err := config.ReadConfig(".env")

	if err != nil {
		log.Fatalf("an error occurred while reading config; %v\n", err)
	}

	b, err := goTelegram.NewBot(config.GetConfig().TBotAPIKey)

	if err != nil {
		log.Fatalf("an error occurred while creating the bot: %v\n", err)
	}

	b.SetHandler(handler)

	repo, err := mysql.NewMySQLHandler()

	if err != nil {
		log.Fatalf("an error occurred while creating the repo: %v\n", err)
	}

	manager = service.NewManager(&b, repo)

	port := config.GetConfig().PORT

	log.Println("Starting Server On Port " + port + "....")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(b.UpdateHandler)); err != nil {
		log.Fatalln(err)
	}

}

func handler(update goTelegram.Update) {
	switch update.Type {
	case "text":
		processText(update)
	case "callback":
	}
}

func processText(update goTelegram.Update) {
	switch update.Command {
	case "/start":
		manager.SendGenericMessage(fmt.Sprintf("Hello %s\nSend /help To See How To Use This Bot.", update.Message.From.Firstname), update.Message.Chat.ID)
	case "/add":
		message := update.Message
		err := manager.SaveEntry(message.Chat.ID, message.Text)

		if err != nil {
			manager.SendGenericMessage("An error occurred while saving the entry", message.Chat.ID)
			log.Println(err)
			return
		}

		manager.SendGenericMessage("Entry Saved Successfully", message.Chat.ID)
	case "/show":
		message := update.Message
		dates := strings.Split(message.Text, " ")
		result, err := manager.RetrieveSpendingByDateRanges(message.Chat.ID, dates...)

		if err != nil {
			manager.SendGenericMessage("An error occurred while retrieving your entries", message.Chat.ID)
			log.Println(err)
			return
		}

		err = manager.SendSpendingData(message.Chat.ID, result...)

		if err != nil {
			log.Println(err)
		}
	case "/visualize":
	case "/notionize":
	case "/help":
		helpText := `Hello {username}

Here Are The List Of Commands You Can Use For Interacting With This Bot.

You Can Send /command help to get more verbose info on how to use each command ex. /add help
- /add             - Add A New Entry
- /show           - Show your spending
- /visualize     - Send An Image Visualizing Your Spending By Categories
- /notionize    - Create A Notion Page With A Tabulated View Of Your Spending`

		helpText = strings.NewReplacer("{username}", update.Message.From.Firstname).Replace(helpText)

		manager.SendGenericMessage(helpText, update.Message.Chat.ID)
	}
}
