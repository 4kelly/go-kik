package kik

// User is the response body of a User profile from the Kik bot API.
type User struct {
	FirstName              string
	LastName               string
	ProfilePicLastModified int64
	ProfilePicUrl          string
}
