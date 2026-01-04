package toniebox

import (
	"fmt"
)

// Client is the main interface for interacting with the Toniebox API.
// It provides methods for authentication and accessing Toniebox resources.
type Client struct {
	requestHandler *requestHandler
}

// NewClient creates a new Toniebox API client with default settings.
//
// Example:
//
//	client := toniebox.NewClient()
//	err := client.Login("user@example.com", "password")
func NewClient() *Client {
	return &Client{
		requestHandler: newRequestHandler(),
	}
}

// NewClientWithProxy creates a new Toniebox API client with a proxy.
// The proxyURL should be in the format "http://host:port" or "https://host:port".
//
// Example:
//
//	client, err := toniebox.NewClientWithProxy("http://proxy.example.com:8080")
func NewClientWithProxy(proxyURL string) (*Client, error) {
	handler, err := newRequestHandlerWithProxy(proxyURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		requestHandler: handler,
	}, nil
}

// Login authenticates the user with their Toniebox account credentials.
// This must be called before any other API methods.
//
// Parameters:
//   - username: The email address for your Toniebox account
//   - password: The password for your Toniebox account
//
// Returns an error if authentication fails.
//
// Example:
//
//	err := client.Login("user@example.com", "password")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Client) Login(username, password string) error {
	login := &Login{
		Email:    username,
		Password: password,
	}
	return c.requestHandler.login(login)
}

// GetMe retrieves personal information about the authenticated user.
//
// Returns the user's profile information or an error if the request fails.
//
// Example:
//
//	me, err := client.GetMe()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("User: %s %s\n", me.FirstName, me.LastName)
func (c *Client) GetMe() (*Me, error) {
	return c.requestHandler.getMe()
}

// GetHouseholds retrieves all households that the user belongs to.
// A household represents a family or group that shares Tonieboxes.
//
// Returns a slice of households or an error if the request fails.
//
// Example:
//
//	households, err := client.GetHouseholds()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, household := range households {
//	    fmt.Printf("Household: %s (ID: %s)\n", household.Name, household.ID)
//	}
func (c *Client) GetHouseholds() ([]Household, error) {
	return c.requestHandler.getHouseholds()
}

// GetCreativeTonies retrieves all Creative-Tonies in a specific household.
// Creative-Tonies are the figurines that can have custom audio content uploaded to them.
//
// Parameters:
//   - household: The household to retrieve Creative-Tonies from
//
// Returns a slice of Creative-Tonies or an error if the request fails.
//
// Example:
//
//	households, _ := client.GetHouseholds()
//	tonies, err := client.GetCreativeTonies(&households[0])
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tonie := range tonies {
//	    fmt.Printf("Tonie: %s (Chapters: %d)\n", tonie.Name, tonie.ChaptersPresent)
//	}
func (c *Client) GetCreativeTonies(household *Household) ([]CreativeTonie, error) {
	return c.requestHandler.getCreativeTonies(household)
}

// FindChapterByTitle searches for a chapter with the given title on this Creative-Tonie.
//
// Parameters:
//   - title: The title to search for
//
// Returns the chapter if found, or nil if not found.
//
// Example:
//
//	chapter := tonie.FindChapterByTitle("My Audio Track")
//	if chapter != nil {
//	    fmt.Printf("Found chapter: %s\n", chapter.Title)
//	}
func (ct *CreativeTonie) FindChapterByTitle(title string) *Chapter {
	for i := range ct.Chapters {
		if ct.Chapters[i].Title == title {
			return &ct.Chapters[i]
		}
	}
	return nil
}

// DeleteChapter removes a chapter from this Creative-Tonie.
// Note: You must call Commit() after this to persist the changes.
//
// Parameters:
//   - chapter: The chapter to delete
//
// Example:
//
//	chapter := tonie.FindChapterByTitle("Old Track")
//	if chapter != nil {
//	    tonie.DeleteChapter(chapter)
//	    tonie.Commit()
//	}
func (ct *CreativeTonie) DeleteChapter(chapter *Chapter) {
	var newChapters []Chapter
	for i := range ct.Chapters {
		if ct.Chapters[i].ID != chapter.ID {
			newChapters = append(newChapters, ct.Chapters[i])
		}
	}
	ct.Chapters = newChapters
}

// UploadFile uploads an audio file to this Creative-Tonie.
// The file will be added as a new chapter with the specified title.
// Note: You must call Commit() after this to persist the changes.
//
// Parameters:
//   - title: The title for the new chapter
//   - filePath: The path to the audio file to upload
//
// Returns an error if the upload fails.
//
// Example:
//
//	err := tonie.UploadFile("My Story", "/path/to/audio.mp3")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = tonie.Commit()
func (ct *CreativeTonie) UploadFile(title, filePath string) error {
	if ct.requestHandler == nil {
		return fmt.Errorf("tonie not properly initialized")
	}
	return ct.requestHandler.uploadFile(ct, filePath, title)
}

// Commit saves all changes made to this Creative-Tonie to the Toniebox cloud.
// This must be called after making changes like renaming, uploading, or deleting chapters.
//
// Returns an error if the commit fails.
//
// Example:
//
//	tonie.Name = "New Name"
//	err := tonie.Commit()
//	if err != nil {
//	    log.Fatal(err)
//	}
func (ct *CreativeTonie) Commit() error {
	if ct.requestHandler == nil {
		return fmt.Errorf("tonie not properly initialized")
	}
	return ct.requestHandler.commitTonie(ct)
}

// Refresh reloads the current state of this Creative-Tonie from the Toniebox cloud.
// This is useful to see the latest changes, such as transcoding status.
//
// Returns an error if the refresh fails.
//
// Example:
//
//	err := tonie.Refresh()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Chapters present: %d\n", tonie.ChaptersPresent)
func (ct *CreativeTonie) Refresh() error {
	if ct.requestHandler == nil {
		return fmt.Errorf("tonie not properly initialized")
	}

	refreshed, err := ct.requestHandler.refreshTonie(ct)
	if err != nil {
		return err
	}

	// Update fields
	ct.ID = refreshed.ID
	ct.Name = refreshed.Name
	ct.Live = refreshed.Live
	ct.Private = refreshed.Private
	ct.ImageURL = refreshed.ImageURL
	ct.TranscodingErrors = refreshed.TranscodingErrors
	ct.Transcoding = refreshed.Transcoding
	ct.SecondsPresent = refreshed.SecondsPresent
	ct.SecondsRemaining = refreshed.SecondsRemaining
	ct.ChaptersPresent = refreshed.ChaptersPresent
	ct.ChaptersRemaining = refreshed.ChaptersRemaining
	ct.Chapters = refreshed.Chapters
	ct.HouseholdID = refreshed.HouseholdID

	return nil
}
