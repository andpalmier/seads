package main

import (
	"fmt"
	"github.com/containrrr/shoutrrr"
	"time"
)

// Notifier interface used for notification channels
type Notifier interface {
	SendMessage(message string) error
}

// TelegramNotifier holds configurations for sending the message on Telegram
type TelegramNotifier struct {
	Token  string   `yaml:"token"`
	ChatId []string `yaml:"chatid"`
}

// SendMessage sends the specified message on Telegram
func (tn *TelegramNotifier) SendMessage(message string) error {
	chats := tn.ChatId[0]
	if len(tn.ChatId) > 1 {
		for _, mr := range tn.ChatId[1:] {
			chats += "," + mr
		}
	}

	url := fmt.Sprintf("telegram://%s@telegram?channels=%s", tn.Token, chats)
	return shoutrrr.Send(url, message)
}

// SlackNotifier holds configurations for sending the message on Slack
type SlackNotifier struct {
	Token    string   `yaml:"token"`
	Channels []string `yaml:"channels"`
}

// SendMessage sends the specified message on Slack
func (sn *SlackNotifier) SendMessage(message string) error {
	channels := sn.Channels[0]
	if len(sn.Channels) > 1 {
		for _, mr := range sn.Channels[1:] {
			channels += "," + mr
		}
	}

	url := fmt.Sprintf("slack://%s@%s", sn.Token, channels)
	return shoutrrr.Send(url, message)
}

// MailNotifier holds configurations for sending the message via email
type MailNotifier struct {
	Host       string   `yaml:"host"`
	Port       string   `yaml:"port"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Auth       string   `yaml:"auth"`
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// SendMessage sends the specified message via email
func (mn *MailNotifier) SendMessage(message string) error {
	mailrecipients := mn.Recipients[0]
	if len(mn.Recipients) > 1 {
		for _, mr := range mn.Recipients[1:] {
			mailrecipients += "," + mr
		}
	}
	url := fmt.Sprintf("smtp://%s:%s@%s:%s/?from=%s&to=%s&subject=seadscan notification",
		mn.Username, mn.Password, mn.Host, mn.Port, mn.From, mailrecipients)
	return shoutrrr.Send(url, message)
}

// notify creates the message to be sent and sends it using the specified notification services
func (config *Config) notify(toSend []ResultAd) {
	message := createMessage(toSend)

	notifiers := []Notifier{}
	if config.SlackNotifier != nil {
		notifiers = append(notifiers, config.SlackNotifier)
	}
	if config.TelegramNotifier != nil {
		notifiers = append(notifiers, config.TelegramNotifier)
	}
	if config.MailNotifier != nil {
		notifiers = append(notifiers, config.MailNotifier)
	}

	notificationSent := false

	for _, notifier := range notifiers {
		err := notifier.SendMessage(message)
		if err != nil {
			fmt.Printf("error sending message via notifier: %v\n", err)
			continue
		}
		notificationSent = true
	}
	if notificationSent {
		fmt.Println("notifications sent!")
	}
}

// createMessage assembles the message to be sent over the specified notification channels
func createMessage(toSend []ResultAd) string {
	message := "Here are the \"unexpected domains\" found during the last execution of seads:\n\n " +
		"Message creation date: " + time.Now().Format(time.DateTime) + "\n\n"
	for _, s := range toSend {
		m := formatNotification(s)
		message += m + "\n"
	}
	message += "\nThis message was automatically sent by seads (www.github.com/andpalmier/seads)"
	return message
}

// formatNotification formats the notification message
func formatNotification(resultAd ResultAd) string {
	return fmt.Sprintf("* Search engine: %s\n\tSearch term: %s\n\tDomain: %s\n\tFull link: %s\n",
		resultAd.Engine, resultAd.Query, DefangURL(resultAd.Domain), DefangURL(resultAd.Link))
}
