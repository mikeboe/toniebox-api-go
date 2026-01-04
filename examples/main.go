package main

import (
	"fmt"
	"log"
	"os"

	toniebox "github.com/mikeboe/toniebox-api-go"
)

func main() {
	// Get credentials from environment variables
	username := os.Getenv("TONIEBOX_USERNAME")
	password := os.Getenv("TONIEBOX_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Please set TONIEBOX_USERNAME and TONIEBOX_PASSWORD environment variables")
	}

	// Create a new client
	client := toniebox.NewClient()

	// Login to the Toniebox API
	fmt.Println("Logging in...")
	if err := client.Login(username, password); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("✓ Login successful")

	// Make sure to disconnect when done
	defer func() {
		fmt.Println("\nDisconnecting...")
		if err := client.Disconnect(); err != nil {
			log.Printf("Warning: Disconnect failed: %v", err)
		} else {
			fmt.Println("✓ Disconnected")
		}
	}()

	// Get personal information
	fmt.Println("\nFetching user information...")
	me, err := client.GetMe()
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("✓ Hello, %s %s! (Email: %s)\n", me.FirstName, me.LastName, me.Email)

	// Get all households
	fmt.Println("\nFetching households...")
	households, err := client.GetHouseholds()
	if err != nil {
		log.Fatalf("Failed to get households: %v", err)
	}
	fmt.Printf("✓ Found %d household(s)\n", len(households))

	if len(households) == 0 {
		fmt.Println("No households found.")
		return
	}

	// Use the first household
	household := &households[0]
	fmt.Printf("\nUsing household: %s (ID: %s)\n", household.Name, household.ID)

	// Get all Creative-Tonies in the household
	fmt.Println("\nFetching Creative-Tonies...")
	tonies, err := client.GetCreativeTonies(household)
	if err != nil {
		log.Fatalf("Failed to get Creative-Tonies: %v", err)
	}
	fmt.Printf("✓ Found %d Creative-Tonie(s)\n", len(tonies))

	// Display information about each Tonie
	for i, tonie := range tonies {
		fmt.Printf("\n--- Creative-Tonie #%d ---\n", i+1)
		fmt.Printf("  Name: %s\n", tonie.Name)
		fmt.Printf("  ID: %s\n", tonie.ID)
		fmt.Printf("  Chapters Present: %d\n", tonie.ChaptersPresent)
		fmt.Printf("  Chapters Remaining: %d\n", tonie.ChaptersRemaining)
		fmt.Printf("  Seconds Present: %.2f\n", tonie.SecondsPresent)
		fmt.Printf("  Seconds Remaining: %.2f\n", tonie.SecondsRemaining)
		fmt.Printf("  Live: %t\n", tonie.Live)
		fmt.Printf("  Private: %t\n", tonie.Private)
		fmt.Printf("  Transcoding: %t\n", tonie.Transcoding)

		if len(tonie.Chapters) > 0 {
			fmt.Printf("\n  Chapters:\n")
			for j, chapter := range tonie.Chapters {
				fmt.Printf("    %d. %s (%.2f seconds)\n", j+1, chapter.Title, chapter.Seconds)
			}
		}
	}

	// Example: Working with a specific Creative-Tonie
	if len(tonies) > 0 {
		fmt.Println("\n\n=== Example Operations ===")
		tonie := &tonies[0]
		fmt.Printf("Working with: %s\n", tonie.Name)

		// Example: Refresh to get the latest state
		fmt.Println("\nRefreshing Tonie state...")
		if err := tonie.Refresh(); err != nil {
			log.Printf("Failed to refresh: %v", err)
		} else {
			fmt.Printf("✓ Refreshed. Current chapters: %d\n", tonie.ChaptersPresent)
		}

		// Example: Find a chapter by title
		if len(tonie.Chapters) > 0 {
			firstChapter := tonie.Chapters[0]
			fmt.Printf("\nSearching for chapter: %s\n", firstChapter.Title)
			found := tonie.FindChapterByTitle(firstChapter.Title)
			if found != nil {
				fmt.Printf("✓ Found chapter: %s (ID: %s)\n", found.Title, found.ID)
			}
		}

		// Example: Rename the Tonie (commented out to avoid unwanted changes)
		/*
			fmt.Println("\nExample: Renaming Tonie (commented out)")
			originalName := tonie.Name
			tonie.Name = "Test Name"
			if err := tonie.Commit(); err != nil {
				log.Printf("Failed to commit: %v", err)
			} else {
				fmt.Println("✓ Tonie renamed")
				// Restore original name
				tonie.Name = originalName
				tonie.Commit()
			}
		*/

		// Example: Upload a file (commented out to avoid unwanted changes)
		/*
			fmt.Println("\nExample: Uploading a file (commented out)")
			if err := tonie.UploadFile("My Audio", "/path/to/audio.mp3"); err != nil {
				log.Printf("Failed to upload: %v", err)
			} else {
				if err := tonie.Commit(); err != nil {
					log.Printf("Failed to commit: %v", err)
				} else {
					fmt.Println("✓ File uploaded")
				}
			}
		*/
	}

	fmt.Println("\n✓ Example completed successfully!")
}
