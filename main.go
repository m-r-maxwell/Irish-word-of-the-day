package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"gopkg.in/gomail.v2"
)

func main() {
	// on my machine this is a cronjob, to see all cron jobs for a user enter crontab -e
	url := "https://www.teanglann.ie/en"
	c := colly.NewCollector()

	var wordOfTheDay, definition string

	c.OnHTML(".wod .entry", func(e *colly.HTMLElement) {
		wordOfTheDay = strings.TrimSpace(e.ChildText("a.headword"))
		definitions := e.ChildTexts(".sense .def")
		definition = strings.Join(definitions, ", ")
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		message := fmt.Sprintf("Word of the Day: %s\nDefinition: %s\nFrom: %s\n", wordOfTheDay, definition, url)

		sender := "your-email-here"
		password := "your-password-here"
		recipient := "your-email-here"

		m := gomail.NewMessage()
		m.SetHeader("From", sender)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", "Irish Word of the Day")
		m.SetBody("text/plain", message)

		d := gomail.NewDialer("smtp-goes-here", 1000, sender, password)

		if err := d.DialAndSend(m); err != nil {
			log.Fatal(err)
		}
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
