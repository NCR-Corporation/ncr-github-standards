package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

/*
The SendSlackAlert function takes a list of users who do not use quicklook IDs for login and sends an alert in slack channel specified in the secret webhook
	invalidUsername: A list of login ids which are not in quicklook ID format
*/

func SendSlackAlert(invalidUsername []string) {
	// This is the webhook of Slack channel to which the message will be sent.
	url := GetSecret("<YOUR SLACK BOT API TOKEN HERE>")

	// Building the message
	str := fmt.Sprintf(":memo: The following *%d* login ids do not comply with Github Standards:", len(invalidUsername))
	headerText := slack.NewTextBlockObject("mrkdwn", str, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	divider := slack.NewDividerBlock();

	list := ""
	for _, user := range invalidUsername {
		list += "<https://github.com/"+user + "|" + user + ">" + "\n"
	}
	listText := slack.NewTextBlockObject("mrkdwn", list, false, false)
	listSection := slack.NewSectionBlock(listText, nil, nil)

	msg := slack.NewBlockMessage(
		headerSection,
		divider,
		listSection,
	)

	//Build a JSON response
	b, err := json.Marshal(&msg)
	if err != nil {
		log.Println(err)
		return
	}


	// Create a request and send
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Printf("Slack message sent. %d users do not use quicklook ID as login", len(invalidUsername))
}
