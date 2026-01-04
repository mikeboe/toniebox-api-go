package toniebox

// JWTToken represents the authentication token returned by the API
type JWTToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// Login represents the credentials for logging into the Toniebox API
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Me represents personal information about the authenticated user
type Me struct {
	Email                        string `json:"email"`
	UUID                         string `json:"uuid"`
	FirstName                    string `json:"firstName"`
	LastName                     string `json:"lastName"`
	Sex                          string `json:"sex"`
	AcceptedTermsOfUse           bool   `json:"acceptedTermsOfUse"`
	Tracking                     bool   `json:"tracking"`
	AuthCode                     string `json:"authCode"`
	ProfileImage                 string `json:"profileImage"`
	Verified                     bool   `json:"isVerified"`
	EduUser                      bool   `json:"isEduUser"`
	NotificationCount            int    `json:"notificationCount"`
	RequiresVerificationToUpload bool   `json:"requiresVerificationToUpload"`
}

// Household represents a Toniebox household
type Household struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	Image                       string `json:"image"`
	ForeignCreativeTonieContent bool   `json:"foreignCreativeTonieContent"`
	Access                      string `json:"access"`
	CanLeave                    bool   `json:"canLeave"`
	OwnerName                   string `json:"ownerName"`
}

// Chapter represents a chapter/track on a Creative-Tonie
type Chapter struct {
	ID          string  `json:"id"`
	File        string  `json:"file"`
	Title       string  `json:"title"`
	Seconds     float64 `json:"seconds"`
	Transcoding bool    `json:"transcoding"`
}

// CreativeTonie represents a Creative-Tonie figurine
type CreativeTonie struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Live              bool      `json:"live"`
	Private           bool      `json:"private"`
	ImageURL          string    `json:"imageUrl"`
	TranscodingErrors []string  `json:"transcodingErrors"`
	Transcoding       bool      `json:"transcoding"`
	SecondsPresent    float64   `json:"secondsPresent"`
	SecondsRemaining  float64   `json:"secondsRemaining"`
	ChaptersPresent   int       `json:"chaptersPresent"`
	ChaptersRemaining int       `json:"chaptersRemaining"`
	Chapters          []Chapter `json:"chapters"`
	HouseholdID       string    `json:"householdId"`

	// Internal fields not serialized to JSON
	household      *Household      `json:"-"`
	requestHandler *requestHandler `json:"-"`
}

// AmazonBean represents the Amazon S3 upload response
type AmazonBean struct {
	FileID  string      `json:"fileId"`
	Request RequestBean `json:"request"`
}

// RequestBean represents Amazon S3 upload request details
type RequestBean struct {
	URL    string     `json:"url"`
	Fields FieldsBean `json:"fields"`
}

// FieldsBean represents Amazon S3 upload form fields
type FieldsBean struct {
	Key               string `json:"key"`
	Policy            string `json:"policy"`
	XAmzAlgorithm     string `json:"x-amz-algorithm"`
	XAmzCredential    string `json:"x-amz-credential"`
	XAmzDate          string `json:"x-amz-date"`
	XAmzSignature     string `json:"x-amz-signature"`
	XAmzSecurityToken string `json:"x-amz-security-token"`
}
