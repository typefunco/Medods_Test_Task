package utils

import (
	"fmt"
	"time"
)

var mockEmails = map[string]string{
	"123e4567-e89b-12d3-a456-426614174000": "user1@example.com",
	"223e4567-e89b-12d3-a456-426614174001": "user2@example.com",
}

func SendWarningEmail(userID string) {
	email, ok := mockEmails[userID]
	if !ok {
		fmt.Println("Email not found for user:", userID)
		return
	}

	subject := "Security Warning: IP Address Changed"
	body := fmt.Sprintf("Dear user,\n\nWe've detected a change in your IP address during your latest request. If this was not you, please take appropriate security measures.\n\nTime: %s\n", time.Now().Format(time.RFC1123))

	fmt.Printf("Sending email to: %s\nSubject: %s\n\n%s\n", email, subject, body)

}
