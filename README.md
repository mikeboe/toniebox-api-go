# Toniebox API - Go

> A Go port of [maximilianvoss/toniebox-api](https://github.com/maximilianvoss/toniebox-api)

A Go library for interacting with the Toniebox Cloud API. Control your Creative-Tonies, upload audio files, manage chapters, and more - all from your Go applications.

## What is Toniebox?

The [Toniebox](https://tonies.com) is a popular audio player designed for children. It uses figurines called "Tonies" that play stories, songs, and other audio content. Creative-Tonies are special figurines that allow you to upload your own audio content.

## Features

- 🔐 **Authentication** - Secure login with your Toniebox account credentials
- 👤 **User Management** - Retrieve your personal account information
- 🏠 **Household Access** - List and manage multiple households
- 🎭 **Creative-Tonie Control** - Access and control all your Creative-Tonies
- 📤 **File Upload** - Upload audio files to your Creative-Tonies
- 📝 **Chapter Management** - Add, remove, and organize chapters
- 🔄 **Real-time Updates** - Refresh to get the latest state from the cloud
- 🌐 **Proxy Support** - Optional proxy configuration for network environments

## Installation

```bash
go get github.com/mikeboe/toniebox-api-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    toniebox "github.com/mikeboe/toniebox-api-go"
)

func main() {
    // Create a new client
    client := toniebox.NewClient()
    
    // Login
    if err := client.Login("user@example.com", "password"); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect()
    
    // Get your households
    households, err := client.GetHouseholds()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get Creative-Tonies
    tonies, err := client.GetCreativeTonies(&households[0])
    if err != nil {
        log.Fatal(err)
    }
    
    // Work with a Tonie
    tonie := &tonies[0]
    fmt.Printf("Tonie: %s has %d chapters\n", tonie.Name, tonie.ChaptersPresent)
}
```

## Usage Examples

### Authentication

```go
// Create a client
client := toniebox.NewClient()

// Or with proxy support
client, err := toniebox.NewClientWithProxy("http://proxy.example.com:8080")

// Login
err := client.Login("user@example.com", "password")
```

### Get User Information

```go
me, err := client.GetMe()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User: %s %s\n", me.FirstName, me.LastName)
```

### List Households and Creative-Tonies

```go
// Get all households
households, err := client.GetHouseholds()
if err != nil {
    log.Fatal(err)
}

// Get Creative-Tonies from a household
tonies, err := client.GetCreativeTonies(&households[0])
if err != nil {
    log.Fatal(err)
}

for _, tonie := range tonies {
    fmt.Printf("Tonie: %s (Chapters: %d)\n", tonie.Name, tonie.ChaptersPresent)
}
```

### Upload an Audio File

```go
tonie := &tonies[0]

// Upload a file
err := tonie.UploadFile("My Story", "/path/to/audio.mp3")
if err != nil {
    log.Fatal(err)
}

// Commit changes
err = tonie.Commit()
if err != nil {
    log.Fatal(err)
}
```

### Manage Chapters

```go
// Find a chapter by title
chapter := tonie.FindChapterByTitle("Old Story")
if chapter != nil {
    // Delete the chapter
    tonie.DeleteChapter(chapter)
    
    // Commit changes
    tonie.Commit()
}
```

### Rename a Creative-Tonie

```go
tonie.Name = "New Name"
err := tonie.Commit()
if err != nil {
    log.Fatal(err)
}
```

### Refresh State

```go
// Refresh to get the latest state from the cloud
err := tonie.Refresh()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Current chapters: %d\n", tonie.ChaptersPresent)
```

## Running the Example

A complete example application is included in the `examples` directory:

```bash
# Set your credentials
export TONIEBOX_USERNAME="your-email@example.com"
export TONIEBOX_PASSWORD="your-password"

# Run the example
cd examples
go run main.go
```

## API Documentation

For detailed API documentation, see the [GoDoc](https://pkg.go.dev/github.com/mikeboe/toniebox-api-go).

### Main Types

- **Client** - Main API client for authentication and accessing resources
- **CreativeTonie** - Represents a Creative-Tonie figurine with methods for managing content
- **Household** - Represents a household/family group
- **Chapter** - Represents an audio chapter/track on a Creative-Tonie
- **Me** - User account information

### Main Methods

#### Client Methods
- `NewClient()` - Create a new API client
- `NewClientWithProxy(proxyURL)` - Create a client with proxy support
- `Login(username, password)` - Authenticate with your Toniebox account
- `GetMe()` - Get your user information
- `GetHouseholds()` - List all households you belong to
- `GetCreativeTonies(household)` - List Creative-Tonies in a household
- `Disconnect()` - Terminate the session

#### CreativeTonie Methods
- `UploadFile(title, filePath)` - Upload an audio file
- `Commit()` - Save changes to the cloud
- `Refresh()` - Reload the latest state
- `FindChapterByTitle(title)` - Find a chapter by its title
- `DeleteChapter(chapter)` - Remove a chapter

## Requirements

- Go 1.21 or higher
- Active Toniebox account
- Internet connection

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Credits

This is a Go port of the original Java implementation by [Maximilian Voss](https://github.com/maximilianvoss/toniebox-api).

## Disclaimer

This is an unofficial library and is not affiliated with, endorsed by, or connected to Boxine GmbH or the Toniebox brand. Use at your own risk.