package toniebox

const (
	// API endpoints
	openIDConnect    = "https://login.tonies.com/auth/realms/tonies/protocol/openid-connect/token"
	creativeTonies   = "https://api.tonie.cloud/v2/households/%s/creativetonies"
	creativeTonie    = "https://api.tonie.cloud/v2/households/%s/creativetonies/%s"
	session          = "https://api.tonie.cloud/v2/sessions"
	me               = "https://api.tonie.cloud/v2/me"
	households       = "https://api.tonie.cloud/v2/households"
	fileUpload       = "https://api.tonie.cloud/v2/file"
	fileUploadAmazon = "https://bxn-toniecloud-prod-upload.s3.amazonaws.com/"

	// HTTP headers
	contentTypeJSON = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"

	// OAuth parameters
	grantTypePassword = "password"
	clientID          = "my-tonies"
	scopeOpenID       = "openid"
)
