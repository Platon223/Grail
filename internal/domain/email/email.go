package email

import (
	"fmt"

	"encoding/base64"
	"time"

	gmailapi "google.golang.org/api/gmail/v1"
)


type Email struct {
	ID            string
	ThreadID      string
	Subject       string
	From          string
	To            string
	Date          time.Time
	Snippet       string
	Body          string
	LabelIDs      []string
	Unread        bool
	HasAttachment bool
}


func SendReply(svc *gmailapi.Service, to, subject, body, threadID string) error {

    raw := fmt.Sprintf(
        "To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
        to, subject, body,
    )

    encoded := base64.URLEncoding.EncodeToString([]byte(raw))

    msg := &gmailapi.Message{
        Raw:      encoded,
        ThreadId: threadID,
    }

    _, err := svc.Users.Messages.Send("me", msg).Do()
    return err
}

func FromMessage(msg *gmailapi.Message) (*Email, error) {
	e := &Email{
		ID:       msg.Id,
		ThreadID: msg.ThreadId,
		LabelIDs: msg.LabelIds,
		Snippet:  msg.Snippet,
	}

	for _, l := range msg.LabelIds {
		if l == "UNREAD" {
			e.Unread = true
		}
	}

	if msg.Payload != nil {
		for _, h := range msg.Payload.Headers {
			switch h.Name {
			case "Subject":
				e.Subject = h.Value
			case "From":
				e.From = h.Value
			case "To":
				e.To = h.Value
			case "Date":
				e.Date, _ = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", h.Value)
			}
		}
		e.Body = extractBody(msg.Payload)
		e.HasAttachment = hasAttachments(msg.Payload)
	}

	return e, nil
}

func extractBody(payload *gmailapi.MessagePart) string {
	if payload == nil {
		return ""
	}
	if payload.MimeType == "text/plain" && payload.Body != nil {
		return decode(payload.Body.Data)
	}
	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" && part.Body != nil {
			return decode(part.Body.Data)
		}
	}
	for _, part := range payload.Parts {
		if part.MimeType == "text/html" && part.Body != nil {
			return decode(part.Body.Data)
		}
	}
	for _, part := range payload.Parts {
		if body := extractBody(part); body != "" {
			return body
		}
	}
	return ""
}

func hasAttachments(payload *gmailapi.MessagePart) bool {
	if payload == nil {
		return false
	}
	for _, part := range payload.Parts {
		if part.Body != nil && part.Body.AttachmentId != "" {
			return true
		}
	}
	return false
}

func decode(data string) string {
	b, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		b, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			return ""
		}
	}
	return string(b)
}
