package main

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

/*
	GetUsersAndFilter function queries all the people who are part of your organization and filters them into 4 as follows:
	1. People who do not have their name and email id updated and publically visible
	2. People who do not have their name updated
	3. People who do not have their email id updated
	4. People who use anything else that quicklook id as login id
*/
func GetUsersAndFilter() {
	invalidUsername := make([]string, 0)
	withoutName, withoutEmail, withoutNameAndEmail := 0, 0, 0
	ctx := context.Background()
	apiToken := GetSecret("<PATH TO YOUR GITHUB API TOKEN HERE>")
	tokensource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)
	// Create go-github client
	httpClient := oauth2.NewClient(ctx, tokensource)
	githubClient := github.NewClient(httpClient)
	log.Println("Gihub client created.")
	// Customize request option parameter to get maximum 100 users in 1 request
	requestOptions := github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalUserCount := 0

	/** Feature Flags **/
	// set up LaunchDarkly client
	ld_client, _ := ld.MakeClient(GetSecret("<PATH TO YOUR LAUNCHDARKLY API TOKEN HERE"), 5*time.Second)
	if ld_client.Initialized() {
		log.Println("LaunchDarkly SDK successfully initialized!")
	} else { //If LaunchDarkly client does not initilize, abort.
		log.Panic("LaunchDarkly SDK failed to initialize")
	}

	// Regular expression which checks for exact string match, string to start with 2 alplabets and end with 6 letters.
	var validID = regexp.MustCompile(`^[[:alpha:]]{2}[[:digit:]]{6}$`)
	log.Println("Starting processing users.")
	for {
		memberList, response, error := githubClient.Organizations.ListMembers(ctx, "<YOUR ORGANIZATION HERE>", &requestOptions)
		if error != nil {
			log.Panic(error)
		}

		for _, member := range memberList {
			curMember, _, err := githubClient.Users.Get(ctx, *member.Login)
			if err != nil {
				log.Panic(err)
			}

			/** Feature Flag User Setup **/
			// create a user in LaunchDarkly
			user := lduser.NewUserBuilder(*member.Login).
				Build()

			// if userID is marked as a specified key in LaunchDarkly, they will recieve an email.
			sendEmail, _ := ld_client.BoolVariation("<YOUR LAUNCHDARKLY FLAG HERE>", user, false)
			if sendEmail {
				log.Println(user, "feature flag: can be sent email")
			}

			if !validID.MatchString(*member.Login) {
				// Person does not use quicklook id as login
				invalidUsername = append(invalidUsername, *member.Login)

			} else {
				if curMember.Name == nil && curMember.Email == nil {
					// Person does not have name and email updated
					withoutNameAndEmail++
					if sendEmail {
						Compose_emails(*member.Login, "both", "Github")
					}

				} else if curMember.Name == nil {
					// Person does not have name updated
					withoutName++
					if sendEmail {
						Compose_emails(*member.Login, "name", "Github")
					}

				} else if curMember.Email == nil {
					// Person does not have email updated
					withoutEmail++
					if sendEmail {
						Compose_emails(*member.Login, "email", "Github")
					}
				}
			}
			totalUserCount++
		}

		log.Printf("%d users processed.\n", totalUserCount)
		if response.NextPage == 0 {
			break
		}
		requestOptions.Page = response.NextPage
	}
	// close LaunchDarkly after sending emails
	ld_client.Close()
	// Send slack alert for non standard logins
	SendSlackAlert(invalidUsername)

	log.Printf("Total %d users processed\n", totalUserCount)
	log.Printf("Users with non-standard profile: %d\n", withoutNameAndEmail+withoutName+withoutEmail+len(invalidUsername))
	log.Printf("Users without name on their profile: %d\n", withoutName)
	log.Printf("Users without email id on their profile: %d\n", withoutEmail)
	log.Printf("Users without name and email on their profile: %d", withoutNameAndEmail)
}
