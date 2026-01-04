package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	toniebox "github.com/mikeboe/toniebox-api-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// It's okay if .env doesn't exist, we might be using env vars directly
	}

	// In a real scenario, you would get these from your database or storage
	// For this example, we'll first login to get a token, then demonstrate using it
	username := os.Getenv("TONIEBOX_USERNAME")
	password := os.Getenv("TONIEBOX_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set TONIEBOX_USERNAME and TONIEBOX_PASSWORD environment variables")
	}

	// 1. Initial Login to get a token
	fmt.Println("1. Performing initial login to get a token...")
	initialClient := toniebox.NewClient()
	token, err := initialClient.Login(username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("✓ Got token (Access: %s..., Refresh: %s...)\n",
		token.AccessToken[:10], token.RefreshToken[:10])

	// 2. Create a NEW client and use the token directly
	fmt.Println("\n2. Creating a new client using the existing token...")
	newClient := toniebox.NewClient()

	// Set the token directly
	newClient.SetToken(token)

	// 3. Verify it works by fetching user info
	fmt.Println("3. Verifying authentication with new client...")
	me, err := newClient.GetMe()
	if err != nil {
		log.Fatalf("Failed to get user info with token: %v", err)
	}
	fmt.Printf("✓ Success! Hello, %s %s!\n", me.FirstName, me.LastName)

	// 4. Verify it works by fetching households
	households, err := newClient.GetHouseholds()
	if err != nil {
		log.Fatalf("Failed to get households with token: %v", err)
	}
	fmt.Printf("✓ Successfully retrieved %d household(s)\n", len(households))
}
