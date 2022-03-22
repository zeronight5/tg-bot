package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type AppConfig struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

func main() {
	var args = os.Args
	var configFile = args[1]
	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config AppConfig
	_ = json.Unmarshal(configJson, &config)

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help", "start":
				msg.Text = "type /torrent ID."
			case "torrent":
				arguments := update.Message.CommandArguments()
				if arguments == "" {
					msg.Text = random(config.Host)
				} else {
					magnet := getMagnet(search(config.Host, arguments))
					msg.Text = magnet
				}
			default:
				msg.Text = "I don't know that command"
			}
			_, err := bot.Send(msg)
			if err != nil {
				continue
			}
		}
	}
}

func random(host string) string {
	html := parseHtml("https://" + host + "/search")
	if html == nil {
		return "nothing."
	}
	find := html.Find(".tags-box a")
	ids := make([]string, find.Size())
	find.Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		ids = append(ids, link)
	})
	link := ids[rand.Intn(len(ids))]
	id := link[strings.LastIndex(link, "/")+1:]
	return getMagnet(search(host, id))
}

func parseHtml(url string) *goquery.Document {
	if strings.HasPrefix(url, "//") {
		url = "https:" + url
	}
	if url == "nothing." {
		return nil
	}
	fmt.Printf("visit link: %s\n", url)
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return nil
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
		return nil
	}
	return doc
}

func search(host string, id string) string {
	doc := parseHtml("https://" + host + "/search/" + id)

	if doc == nil {
		return "nothing."
	}

	var url string

	// Find the review items
	doc.Find(".data-list a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		link, _ := s.Attr("href")
		sizeDate := s.Find(".size-date").Text()
		size := strings.Split(sizeDate, " / ")[0]
		size = strings.Split(size, ":")[1]
		// fmt.Printf("Review %d: %s - %s\n", i, link, size)
		if strings.Contains(size, "GB") {
			sizeGB := strings.ReplaceAll(size, "GB", "")
			sizeNum, _ := strconv.ParseFloat(sizeGB, 2)
			if sizeNum > 4 {
				url = link
				return
			}
		}
	})
	return url
}

func getMagnet(url string) string {
	if url == "" {
		return "nothing."
	}
	doc := parseHtml(url)

	if doc == nil {
		return "nothing."
	}

	var magnet string

	// Find the review items
	doc.Find(".magnet-link").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		magnet = s.Text()
		if magnet != "" {
			return
		}
	})
	return magnet
}
