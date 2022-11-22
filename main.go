package main

import (
	"cfo/config"
	"cfo/model"
	"cfo/repository/db/mysql"
	"cfo/service"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/h4ckitt/goTelegram"
)

var (
	manager *service.Manager
	pending model.PendingReplyCache
)

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

	pending = model.NewPendingReplyCache()
	manager = service.NewManager(&b, repo, pending)

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
		processCallback(update)
	}
}

func processText(update goTelegram.Update) {
	switch update.Command {
	case "/start":
		manager.SendGenericMessage(fmt.Sprintf("Hello %s\nSend /help To See How To Use This Bot.", update.Message.From.Firstname), update.Message.Chat.ID)
	case "/add":
		message := update.Message
		text := strings.TrimLeft(message.Text, "/add ")
		err := manager.SaveEntry(message.Chat.ID, message.MessageID, text)

		if err != nil {
			manager.SendGenericMessage("An error occurred while saving the entry", message.Chat.ID)
			log.Println(err)
			return
		}

	case "/show":
		message := update.Message
		dates := strings.Fields(message.Text)
		var result []model.Spending
		var err error
		if len(dates) == 1 { // only /show was passed
			result, err = manager.RetrieveSpendingByDateRanges(message.Chat.ID)
		} else {
			switch dates[1] {
			case "week":
				result, err = manager.RetrieveThisWeekSpending(message.Chat.ID)
			case "yesterday":
				result, err = manager.RetrieveYesterdaySpending(message.Chat.ID)
			case "month":
				result, err = manager.RetrieveThisMonthSpending(message.Chat.ID)
			case "category":
				var (
					result []model.CategorySpending
					err    error
				)
				// if only /show category was sent
				if len(dates) == 2 {
					result, err = manager.RetrieveCategoriesSpendingByDateRanges(message.Chat.ID)
				} else {
					result, err = manager.RetrieveCategoriesSpendingByDateRanges(message.Chat.ID, dates[2:]...)
				}

				if err != nil {
					manager.SendGenericMessage("An error occurred while retrieving your entries", message.Chat.ID)
					log.Println(err)
					return
				}

				err = manager.SendCategorySpending(message.Chat.ID, result...)

				if err != nil {
					log.Println(err)
				}
				return
			default:
				result, err = manager.RetrieveSpendingByDateRanges(message.Chat.ID, dates[1:]...)
			}
		}
		if err != nil {
			manager.SendGenericMessage("An error occurred while retrieving your entries", message.Chat.ID)
			log.Println(err)
			return
		}

		err = manager.SendSpendingData(message.Chat.ID, result...)

		if err != nil {
			log.Println(err)
		}
		return
	case "/visualize":
		fallthrough
	case "/notionize":
		message := update.Message
		manager.SendGenericMessage("This Feature Is Coming Soon", message.Chat.ID)
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
	default:
		if pending.HasCachedSpending(update.Message.Chat.ID) {
			err := manager.CompleteSaveEntry(update.Message.Chat.ID, strings.ToLower(update.Message.Text), true)

			if err != nil {
				log.Println(err)
			}
			manager.DeleteMessage(update.Message.MessageID-1, update.Message.Chat.ID)
			manager.SendGenericMessage("Entry Saved Successfully", update.Message.Chat.ID)
		}
	}
}

func processCallback(update goTelegram.Update) {
	defer func() {
		err := manager.AnswerCallBackQuery(update.CallbackQuery.ID)

		if err != nil {
			log.Println(err)
		}
	}()
	callBacks := strings.SplitN(update.CallbackQuery.Data, "-", 2)

	switch callBacks[0] {
	case "category":
		err := manager.CompleteSaveEntry(update.CallbackQuery.Message.Chat.ID, callBacks[1], false)

		if err != nil {
			log.Println(err)
		}
		manager.EditText("Entry Saved Successfully", update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Chat.ID)

	}
}
