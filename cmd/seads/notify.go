package main

import (
	"fmt"
	"github.com/containrrr/shoutrrr"
	"time"
)

// Notifier interface used for notification channels
type Notifier interface {
	SendNotificationMessage(message string) error
}

// TelegramNotifier holds configurations for sending the message via Telegram
type TelegramNotifier struct {
	Token   string   `yaml:"token"`
	ChatIDs []string `yaml:"chatids"`
}

// SendNotificationMessage sends the specified message via Telegram
func (tn *TelegramNotifier) SendNotificationMessage(message string) error {
	chats := tn.ChatIDs[0]
	if len(tn.ChatIDs) > 1 {
		for _, mr := range tn.ChatIDs[1:] {
			chats += "," + mr
		}
	}

	url := fmt.Sprintf("telegram://%s@telegram?channels=%s", tn.Token, chats)
	return shoutrrr.Send(url, message)
}

// SlackNotifier holds configurations for sending the message via Slack
type SlackNotifier struct {
	Token    string   `yaml:"token"`
	Channels []string `yaml:"channels"`
}

// SendNotificationMessage sends the specified message via Slack
func (sn *SlackNotifier) SendNotificationMessage(message string) error {
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
	From       string   `yaml:"from"`
	Recipients []string `yaml:"recipients"`
}

// SendNotificationMessage sends the specified message via email
func (mn *MailNotifier) SendNotificationMessage(message string) error {
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

// sendNotifications creates the message to be sent and sends it using the specified notification services
func (config *Config) sendNotifications(adsToNotify []AdResult) {
	message := createNotificationMessage(adsToNotify)

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

	notificationsSent := false

	for _, notifier := range notifiers {
		err := notifier.SendNotificationMessage(message)
		if err != nil {
			fmt.Printf("error sending message via notifier: %v\n", err)
			continue
		}
		notificationsSent = true
	}
	if notificationsSent {
		fmt.Println("notifications sent!")
	}
}

// createNotificationMessage assembles the message to be sent over the specified notification channel
func createNotificationMessage(toSend []AdResult) string {
	message := "Here are the \"unexpected domains\" found during the last execution of seads:\n\n" +
		"Message creation date: " + time.Now().Format(time.DateTime) + "\n\n"
	for _, resultAd := range toSend {
		formattedMessage := formatNotificationMessage(resultAd)
		message += formattedMessage + "\n"
	}
	message += "\nThis message was automatically sent by seads (github.com/andpalmier/seads)"
	return message
}

// formatNotificationMessage formats the notification message to be sent
func formatNotificationMessage(resultAd AdResult) string {
	return fmt.Sprintf("* Search engine: %s\n\tSearch term: %s\n\tDomain: %s\n\tFull link: %s\n",
		resultAd.Engine, resultAd.Query, defangAdURL(resultAd.FinalDomainURL), defangAdURL(resultAd.FinalRedirectURL))
}
