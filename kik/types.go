package kik

import (
	"encoding/json"
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

Docs for Keyboards: https://dev.kik.com/#/docs/messaging#keyboards
*/

// SuggestedResponseKeyboard is the only keyboard type, if we add more we can utilize the Keyboard struct.
type SuggestedResponseKeyboard struct {
	Type string `json:"type"` // must be "suggested"

	To     string `json:"to,omitempty"`     // defaults to everyone in the conversation.
	Hidden bool   `json:"hidden,omitempty"` // defaults to false.

	// TODO: actually validate.
	Responses []interface{} `json:"responses,omitempty"`
}

// KeyboardTextResponse sets a text message in the keyboard tray.
type KeyboardTextResponse struct {
	Type string `json:"type"` // Type must be "text".
	Body string `json:"body"`

	Metadata string `json:"metadata,omitempty"` // Include an object to be returned back to your bot when the user responds using the picture suggested response. This may be a string or object, as needed.
}

// KeyboardPictureResponse sets a picture in the keyboard tray.
type KeyboardPictureResponse struct {
	Type   string `json:"type"` // Type must be "picture".
	PicUrl string `json:"picUrl"`

	Metadata string `json:"metadata,omitempty"` // Include an object to be returned back to your bot when the user responds using the picture suggested response. This may be a string or object, as needed.
}

// KeyboardFriendPickerResponse sends a friend picker response to the keyboard tray.
// It is a special message type that will allow a user to 'invite' their friends to use your bot.
// When you invoke the friend picker, the user receives a message to invite their friends.
// It must be set before KeyboardTextResponse.
type KeyboardFriendPickerResponse struct {
	Type string `json:"type"` // Must be "friend-picker".

	Body        string   `json:"body,omitempty"`        // The text to be shown to the user on the suggested response
	Min         int8     `json:"min,omitempty"`         // The minimum amount of friends the user can invite, must be between 1 - 100 and less than or equal to max.
	Max         int8     `json:"max,omitempty"`         // The maximum amount of friends the user can invite, must be between 1 - 100 and greater than or equal to min.
	Preselected []string `json:"preselected,omitempty"` // A predetermined list of users to be picked by the friend picker.
	Metadata    string   `json:"metadata,omitempty"`    // Include an object to be returned back to your bot when the user responds using the picture suggested response. This may be a string or object, as needed.
}

/*
Messaging Types

1. Send (from the bot).
2. Receive (from users to the bot).

*/

// Messages is a simple wrapper around a `Message` interface to satisfy the formatting for the Kik API payload.
type Messages struct {
	Messages []Message `json:"messages"`
}

// Message is a dummy interface so that all structs that embedd `Message` share a common interface.
type Message interface {
	message()
}

// Implement the dummy interface
func (t SendMessage) message() { return }

type SendMessage struct {
	To        string                      `json:"to"`                  // The user or group that will receive the message
	Type      string                      `json:"type"`                // The type of message. See Message Types for the values you can see in this field.
	Delay     int                         `json:"delay"`               // An interval (in milliseconds) to wait before sending the message.
	Keyboards []SuggestedResponseKeyboard `json:"keyboards,omitempty"` // SuggestedResponseKeyboard is currently the only valid keyboard type
	Id        string                      `json:"id,omitempty"`        // randomUUID() ID for this message.Use this to link messages to receipts.This will always be present for received messages.
	ChatId    string                      `json:"chatId,omitempty"`    // The identifier for the conversation your bot is involved in. This field is recommended for all responses in order for messages to be routed correctly (for example, if you're messaging a user in a group)
}

// ReceivedMessages is a simple wrapper around a `Receive` interface
// that knows how to Unmarshal the JSON returned from the Kik API into a valid struct Type.
type ReceivedMessages []Receive

// UnmarshalJSON knows how to parse Kik bot API responses into their correct types.
func (v *ReceivedMessages) UnmarshalJSON(data []byte) error {
	// This just splits up the JSON array into the raw JSON for each object
	var raw map[string][]json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	messages := raw["messages"]
	for _, r := range messages {
		// unamrshal into a map to check the "type" field
		var obj map[string]interface{}
		err := json.Unmarshal(r, &obj)
		if err != nil {
			return err
		}

		messageType := ""
		if t, ok := obj["type"].(string); ok {
			messageType = t
		}

		// unmarshal again into the correct type
		var actual Receive
		switch messageType {
		case "text":
			actual = &TextMessageReceive{}
		case "picture":
			actual = &PictureMessageReceive{}
		}

		err = json.Unmarshal(r, actual)
		if err != nil {
			return err
		}
		*v = append(*v, actual)
	}
	return nil
}

// Receive is a dummy interface so that all structs that embedd `Receive` share a common interface.
type Receive interface {
	receive()
}

// Implements dummy interface.
func (t ReceiveMessage) receive() { return }

type ReceiveMessage struct {
	ChatId               string   `json:"chatId"`       // The identifier for the conversation your bot is involved in. This field is recommended for all responses in order for messages to be routed correctly (for example, if you're messaging a user in a group)
	Id                   string   `json:"id"`           // randomUUID() ID for this message.Use this to link messages to receipts.This will always be present for received messages.
	From                 string   `json:"from"`         // The user who sent the message
	Type                 string   `json:"type"`         // The type of message. See Message Types for the values you can see in this field.
	Participants         []string `json:"participants"` // The type of conversation the message originated from.
	Timestamp            int      `json:"timestamp"`    // The time the message was sent from the Kik client
	ReadReceiptRequested bool     `json:"readReceiptRequested"`

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

	PicUrl      string       `json:"picUrl,omitempty"`    // A picture to be displayed in the message.
	Title       string       `json:"title,omitempty"`     // A title to be displayed at the top of the message.
	Text        string       `json:"text,omitempty"`      // Text to be displayed in the middle of the message.
	NoForward   bool         `json:"noForward,omitempty"` //	If true, the message will not be able to be forwarded to other recipients.
	KikJsData   string       `json:"kikJsData,omitempty"` //	A JSON payload that would be passed to a website using Kik.js.
	Attribution *Attribution `json:"attribution,omitempty"`
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
Configuration
*/

type Configuration struct {
	Webhook   string            `json:"webhook"` // A URL to a webhook to which calls will be made when users interact with your bot.
	*Features `json:"features"` // An object describing the features that are active or not active for your bot.

	StaticKeyboard *SuggestedResponseKeyboard `json:"staticKeyboard,omitempty"` // A keyboard object that shows when a user starts to mention your bot in a conversation.
}

type Features struct {
	ManuallySendReadReceipts bool `json:"manuallySendReadReceipts"` // If enabled, your bot will be responsible for sending its own read receipts to users when you receive messages.
	ReceiveReadReceipts      bool `json:"receiveReadReceipts"`      // If enabled, your bot will receive messages of type read-receipt messages from users.
	ReceiveDeliveryReceipts  bool `json:"receiveDeliveryReceipts"`  // If enabled, your bot will receive messages of type delivery-receipt messages from users.
	ReceiveIsTyping          bool `json:"receiveIsTyping"`          // If enabled, your bot will receive messages of type is-typing messages from users.
}

/*
Kik Codes

Kik codes are a type of 2D barcode similar to QR codes that users can scan.
To render an image of your Kik Code, go to https://api.kik.com/v1/code/<id>?c=<color-code>.
The request will return a 1024x1024 PNG-encoded image.

Docs for Kik Codes: https://dev.kik.com/#/docs/messaging#kik-codes-api
*/

type ScanData struct {
	Data string `json:"data"` // Will be embedded in the Kik Code that users can scan.
}

type Code struct {
	Id string `json:"id"` // The ID to reference a generated Kik code.
}

/*
Error Types
*/

var NotMessageTypeError = errors.New("not a valid message type")
var HttpError = errors.New("HTTP request did not return 200")
