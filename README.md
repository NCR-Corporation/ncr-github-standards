# Github Standards

Github Standards is a set of automation to help ensure Github standards are met for users and repos accross the NCR organization. 
Currently, this program checks that users within the specific Github organizations have their full name and email visible.

## Installation

Clone this repository:
```bash
git clone https://github.com/NCR-Corporation/ncr-github-standards.git
```

Install Go using the instructions [here](https://golang.org/doc/install)

Install GCloud Secret Manager Client Libraries:

```bash
go get -u cloud.google.com/go/secretmanager/apiv1
go get -u google.golang.org/genproto/googleapis/cloud/secretmanager/v1
```

Set up Authentication:

```bash
gcloud auth application-default login
```

Install LaunchDarkly Feature flags:
```bash
go get gopkg.in/launchdarkly/go-server-sdk.v5
go get gopkg.in/launchdarkly/go-sdk-common.v2/lduser
```

To run the application directly, run:

```bash
go build
go run *.go
```
from the src folder

## Feature Flags

This application uses LaunchDarkly Feature flags to hide, enable, or disable features before the application is ready for release.
In the case for Github Standards, we wanted to test our api calls and email components before sending out mass emails.
Using feature flags, we could filter our userbase to a select few test users to send mock emails to, before expanding to the entire organization.

**Please extensively test the code using the feature flag to avoid accidently sending mass emails!**

# Running the Program on your own

## API Keys

In order to sucessfully run this project, new users will have to create the following API Keys:
- Github API token
- Slackbot API token
- LaunchDarkly API token
- SendGrid API token

Once created and stored, they can be accessed using the getSecret method.

## Github Actions

At the current time, our Github Actions workflow is set to only run once a year. If you choose to run the program more frequently, update the cron schedule. 
Below is a sample for running the code once a day. 

    - cron: '1 0 * * *'

As tokens for this project were stored in Google Cloud Platform's Secret Manager, in step 3 of the workflow, new users will need to store their GCP Authentication keys in their own Github repository.

*To note: the final run step of the workflow is commented out. When running the program on your own, uncomment line 49 in github_standards_actions.yml.*

## Expanding to other platforms

Currently, Github-Standards has templates to instruct users on how to update their name, email, or both name and email for Github. Because each platform has
different ways to have a user change/update their name or email, to instruct users on how to update their Outlook/Slack/Teams account, there are two options:  
1. Change the contents of existing HTML templates to fit your platform:  
    Navigate to src/templates/ and pick a file to change, and adjust the contents to fit your needs.  
2. Create your own HTML file:  
    1. Using the template.html found src/templates/, create a new HTML file in the same directory with instructions on how to help the user fix their issue.    
    2. In src/compose_emails.go, starting on line 47, replace the following code with the path to your new file  
       - ex. filepath.Abs("./templates/nameOnly.html") -> filepath.Abs("./templates/{YOUR FILE NAME}].html") 

