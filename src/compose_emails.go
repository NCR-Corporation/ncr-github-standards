package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

/**
The Compose Emails function is used to send an email to Github users who need their profile updated
	user: the QLID of the user who needs an email sent to them.
	option: the type of email that needs to be sent
		- name: only the users name needs updating
		- email: only the users email settings need updating
		- both: both name and email need updating
	platform: the type of platform needing updating (ex. Github, Slack, etc)
		"Platform" will be placed in the sender email and email subject (ex. NCR [Platform] Standards)
NOTE: Because each platform has different steps to have users update their profiles, the current html files
(both.html, emailOnly.html, and nameOnly.html) are specific to updating Github. Options are to: change the content of
the html files to fit your needs, or add new ones and replace the string to the correct filepath. Further instructions
found in the README.md
 **/
func Compose_emails(user string, option string, platform string) {
	log.Println("Compose Emails: email being sent to:", user)

	// create new *SGMailV3
	m := mail.NewV3Mail()
	personalization := mail.NewPersonalization()

	// sender info
	from := mail.NewEmail("NCR "+platform+" Standards", "standards@ncr.com")
	m.SetFrom(from)

	// recipient info
	email := user + "@ncr.com"
	to := mail.NewEmail(user, email)
	personalization.AddTos(to)

	// email subject
	personalization.Subject = "Please update your " + platform + " account"

	// selects the type of email to send
	var file string
	if option == "name" { // only name
		file, _ = filepath.Abs("./templates/nameOnly.html")
	} else if option == "email" { // only email
		file, _ = filepath.Abs("./templates/emailOnly.html")
	} else if option == "both" { // need to update both name and email
		file, _ = filepath.Abs("./templates/both.html")
	} else { //Logs the case where an incorrect option is specified for a particular user and returns.
		log.Printf("Incorrect option specified for user %s. Please choose between: name, email, both", user)
		return
	}

	
	// replace QLID with the user's Quicklook ID. If no "QLID" mentioned in html file, will not substitute
	personalization.SetSubstitution("QLID", user)

	msg, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// creating the email body, and add it to the email
	content := mail.NewContent("text/html", string(msg))
	m.AddContent(content)

	// add the recipients and subject
	m.AddPersonalizations(personalization)

	// send the email
	client := sendgrid.NewSendClient(GetSecret("<PATH TO YOUR SENDGRID TOKEN HERE>"))
	response, err := client.Send(m)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent. Response satus code: ", response.StatusCode)
	}
}
