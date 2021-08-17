package main

import (
	"context"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func GetSecret(path string) string {
	secret := accessSecretVersion(path)
	return secret
}

// Calls secret manager API to grab token and signing secret
func accessSecretVersion(name string) string {
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating secret manager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalf("Error calling the secret manager API: %v", err)
	}

	return string(result.Payload.Data)
}
