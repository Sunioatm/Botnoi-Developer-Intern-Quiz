package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

var userContexts sync.Map

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					handleTextMessage(event, message, bot)
				case *linebot.StickerMessage:
					handleStickerMessage(event, message, bot)
				default:
					log.Printf("Unsupported message content: %T\n", message)
				}
			} else {
				log.Printf("Unsupported event type: %T\n", event)
			}
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleTextMessage(event *linebot.Event, message *linebot.TextMessage, bot *linebot.Client) {
	replyToken := event.ReplyToken
	userMessage := strings.ToLower(message.Text)
	userID := event.Source.UserID

	if context, hasContext := userContexts.Load(userID); hasContext {
		handleQuickReplyResponse(replyToken, userMessage, context.(string), bot)
		userContexts.Delete(userID)
		return
	}

	switch userMessage {
	case "sticker", "stickers":
		replyStickerMessage(replyToken, bot)
	// Button reply
	case "button", "buttons":
		replyButtonTemplate(replyToken, bot)
	// Carousel reply
	case "carousel":
		replyCarouselTemplate(replyToken, bot)
	// Quick reply
	case "botnoi":
		replyBotnoiMessage(replyToken, bot, userID)
	// Text reply
	default:
		replyTextMessage(replyToken, message.Text, bot)
	}
}

func replyTextMessage(replyToken, text string, bot *linebot.Client) {
	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent text reply.")
	}
}

func replyStickerMessage(replyToken string, bot *linebot.Client) {
	stickerMessage := linebot.NewStickerMessage("1", "1")
	if _, err := bot.ReplyMessage(replyToken, stickerMessage).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent sticker reply.")
	}
}

func replyButtonTemplate(replyToken string, bot *linebot.Client) {
	imageURL := "https://scontent.fbkk12-5.fna.fbcdn.net/v/t39.30808-6/380592600_707567551394609_124810049998601882_n.jpg?_nc_cat=107&ccb=1-7&_nc_sid=5f2048&_nc_eui2=AeG4zttAXiMN_Gi1WiE7B48g4HvNO4_1m23ge807j_WbbTBz6F977N7oANiUEVKfxQpWXTx2VRrl1i0WoI_26VCw&_nc_ohc=PozFWHK9OmwQ7kNvgGEV_j6&_nc_ht=scontent.fbkk12-5.fna&oh=00_AYB6wuKmGGsc-fKAMPRrcx-u_LsXOZ5RfniyYljsG-NHlA&oe=664A4913"
	buttonTemplate := linebot.NewButtonsTemplate(
		imageURL, "Botnoi", "Please select",
		linebot.NewURIAction("View detail", "https://botnoigroup.com/th/"),
		linebot.NewPostbackAction("Buy", "action=buy&itemid=123", "", ""),
		linebot.NewPostbackAction("Add to cart", "action=add&itemid=123", "", ""),
	)
	templateMessage := linebot.NewTemplateMessage("This is a buttons template", buttonTemplate)

	if _, err := bot.ReplyMessage(replyToken, templateMessage).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent buttons template reply.")
	}
}

func replyCarouselTemplate(replyToken string, bot *linebot.Client) {
	carouselTemplate := linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			"https://pbs.twimg.com/media/GLw0iODaoAAIvC5?format=jpg&name=4096x4096",
			"This is Aespa",
			"Description 1",
			linebot.NewURIAction("Supernova MV", "https://www.youtube.com/watch?v=phuiiNCxRMg"),
			linebot.NewPostbackAction("No action", "action=buy&itemid=222", "", ""),
		),
		linebot.NewCarouselColumn(
			"https://kprofiles.com/wp-content/uploads/2022/02/Nmixx-roller-coaster.jpeg",
			"This is NMIXX",
			"Description 2",
			linebot.NewURIAction("Roller Coaster MV", "https://www.youtube.com/watch?v=fqBAzCH4-9g"),
			linebot.NewPostbackAction("No action", "action=buy&itemid=222", "", ""),
		),
		linebot.NewCarouselColumn(
			"https://kprofiles.com/wp-content/uploads/2024/03/ILLIT-800x800.jpg",
			"This is ILLIT",
			"Description 3",
			linebot.NewURIAction("Magnetic MV", "https://www.youtube.com/watch?v=Vk5-c_v4gMU"),
			linebot.NewPostbackAction("No action", "action=buy&itemid=111", "", ""),
		),
	)

	templateMessage := linebot.NewTemplateMessage("This is a carousel template", carouselTemplate)

	if _, err := bot.ReplyMessage(replyToken, templateMessage).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent carousel template reply.")
	}
}

func replyBotnoiMessage(replyToken string, bot *linebot.Client, userID string) {
	quickReplyItems := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("Yes", "Yes")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("No", "No")),
	)

	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage("Do you like botnoi?").WithQuickReplies(quickReplyItems),
	).Do(); err != nil {
		log.Print(err)
	} else {
		userContexts.Store(userID, "botnoi")
		log.Println("Sent Botnoi question with quick reply options.")
	}
}

func handleQuickReplyResponse(replyToken, userResponse, context string, bot *linebot.Client) {
	var responseMessage string
	if context == "botnoi" {
		if userResponse == "yes" {
			responseMessage = "I’m glad you like Botnoi, because I like too."
		} else if userResponse == "no" {
			responseMessage = "I will pretend that you clicked that by mistake."
		} else {
			responseMessage = "What did you mean by that."
		}
	}

	if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(responseMessage)).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent quick reply response.")
	}
}

func handleStickerMessage(event *linebot.Event, message *linebot.StickerMessage, bot *linebot.Client) {
	replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
		log.Print(err)
	} else {
		log.Println("Sent sticker reply.")
	}
}
