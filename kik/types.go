package kik

import (
	"errors"
)

// User is the response body of a User profile from the Kik bot API.
type User struct {
	FirstName              string
	LastName               string
	ProfilePicLastModified int64
	ProfilePicUrl          string
}

/*
Keyboard Types
*/

// Keyboard would be a useful abstraction if this struct is larger or was used independently.
//type Keyboard struct {
//	To     string `json:"to,omitempty"`     // defaults to everyone in the conversation.
//	Hidden bool   `json:"hidden,omitempty"` // defaults to false.
//	Type   string `json:"type"`
//}

// SuggestedResponseKeyboard is the only keyboard type, if we add more we can utilize the Keyboard struct.
type SuggestedResponseKeyboard struct {
	To     string `json:"to,omitempty"`     // defaults to everyone in the conversation.
	Hidden bool   `json:"hidden,omitempty"` // defaults to false.
	Type   string `json:"type"`             // must be "suggested"
	// Must be one of *Response structs. Validated at runtime due to lack of generics.
	Responses []interface{} `json:"responses,omitempty"`
}

// Response would be a useful abstraction if this struct is larger or was used independently.
//type Response struct {
//	Type     string `json:"type"`
//	Metadata string `json:"metadata,omitempty"`
//}

type KeyboardTextResponse struct {
	Type string `json:"type"`
	Body string `json:"body"`

	Metadata string `json:"metadata,omitempty"`
}

type KeyboardPictureResponse struct {
	Type   string `json:"type"`
	PicUrl string `json:"picUrl"`

	Metadata string `json:"metadata,omitempty"`
}

type KeyboardFriendPickerResponse struct {
	Type string `json:"type"`

	Body        string   `json:"body,omitempty"`
	Min         int8     `json:"min,omitempty"`
	Max         int8     `json:"max,omitempty"`
	Preselected []string `json:"preselected,omitempty"`
	Metadata    string   `json:"metadata,omitempty"`
}

/*
Messaging Types

1. Send (from the bot).
2. Receive (from users to the bot).

*/

type SendMessage struct {
	To        string                      `json:"to"`
	Type      string                      `json:"type"`                // The type of message. See Message Types for the values you can see in this field.
	Delay     int                         `json:"delay"`               // An interval (in milliseconds) to wait before sending the message.
	Keyboards []SuggestedResponseKeyboard `json:"keyboards,omitempty"` // SuggestedResponseKeyboard is currently the only valid keyboard type
	Id        string                      `json:"id,omitempty"`        // randomUUID() ID for this message.Use this to link messages to receipts.This will always be present for received messages.
	ChatId    string                      `json:"chatId,omitempty"`    // The identifier for the conversation your bot is involved in. This field is recommended for all responses in order for messages to be routed correctly (for example, if you're messaging a user in a group)
}

type ReceiveMessage struct {
	ChatId               string `json:"chatId"`       // The identifier for the conversation your bot is involved in. This field is recommended for all responses in order for messages to be routed correctly (for example, if you're messaging a user in a group)
	Id                   string `json:"id"`           // randomUUID() ID for this message.Use this to link messages to receipts.This will always be present for received messages.
	From                 string `json:"from"`         // The user who sent the message
	Type                 string `json:"type"`         // The type of message. See Message Types for the values you can see in this field.
	Participants         string `json:"participants"` // The type of conversation the message originated from.
	Timestamp            int    `json:"timestamp"`
	ReadReceiptRequested bool   `json:"readReceiptRequested"`

	ChatType string `json:"chatType,omitempty"` // The type of conversation the message originated from.
	Mention  string `json:"mention,omitempty"`  // The username of the bot mentioned in the message.
	Metadata string `json:"metadata,omitempty"` // Metadata that was provided by your bot when sending the user a suggested response.
}

type Attribution struct {
	Name    string `json:"name"`              // The name that will appear in the attribution bar.
	IconUrl string `json:"iconUrl,omitempty"` // BROKEN in KIK API: The URL specifying an icon that will appear in the attribution bar.
}

// TextMessage for sending from the bot.
type TextMessage struct {
	SendMessage
	Body     string `json:"body"`               // The text of the message.
	TypeTime int    `json:"typeTime,omitempty"` // An interval (in milliseconds) to appear to be typing to the recipient before the message is sent. This occurs after delay.
}

// TextMessageReceive is the data structure returned from the Kik API when a user sends the bot a text message.
type TextMessageReceive struct {
	ReceiveMessage
	Body string `json:"body"` // The text of the message.
}

type PictureMessage struct {
	SendMessage
	PicUrl      string       `json:"picUrl"`
	Attribution *Attribution `json:"attribution,omitempty"`
}
type PictureMessageReceive struct {
	ReceiveMessage
	PicUrl      string       `json:"picUrl"`
	Attribution *Attribution `json:"attribution,omitempty"`
}

type LinkMessage struct {
	SendMessage
	Url string `json:"url"`

	PicUrl    string `json:"picUrl,omitempty"`    // A picture to be displayed in the message.
	Title     string `json:"title,omitempty"`     // A title to be displayed at the top of the message.
	Text      string `json:"text,omitempty"`      // Text to be displayed in the middle of the message.
	NoForward bool   `json:"noForward,omitempty"` //	If true, the message will not be able to be forwarded to other recipients.
	KikJsData string `json:"kikJsData,omitempty"` //	A JSON payload that would be passed to a website using Kik.js.
	//Attribution *Attribution `json:"attribution,omitempty"`
}
type LinkMessageReceive struct {
	ReceiveMessage
	Url string `json:"url"`

	PicUrl      string       `json:"picUrl,omitempty"`    // A picture to be displayed in the message.
	NoForward   bool         `json:"noForward,omitempty"` //	If true, the message will not be able to be forwarded to other recipients.
	KikJsData   string       `json:"kikJsData,omitempty"` //	A JSON payload that would be passed to a website using Kik.js.
	Attribution *Attribution `json:"attribution,omitempty"`
}
type VideoMessage struct {
	SendMessage
	VideoUrl string `json:"videoUrl"` // The URL of the video or GIF you wish to send.

	Loop        bool         `json:"loop,omitempty"`     // Whether or not the video should loop when played.
	Muted       bool         `json:"muted,omitempty"`    // Whether or not the video should be played without audio.
	Autoplay    bool         `json:"autoplay,omitempty"` // Whether or not the video should be played inline.These messages will only be played inline if they are below 1 MB in size.
	NoSave      bool         `json:"noSave,omitempty"`   // If true, the user will not be allowed to save the video to their device.
	Attribution *Attribution `json:"attribution,omitempty"`
}

type VideoMessageReceive struct {
	ReceiveMessage
	VideoUrl    string       `json:"videoUrl"` // The URL of the video or GIF you wish to send.
	Attribution *Attribution `json:"attribution,omitempty"`
}

/*
Error Types
*/

var NotMessageTypeError = errors.New("not a valid message type")
var HttpError = errors.New("HTTP request did not return 200")
