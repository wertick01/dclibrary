package models

import "time"

type Booking struct {
	Id             int           `json: "id"`
	BookId         int           `json: "book_id"`
	UserId         int           `json: "user_id"`
	DateOfIssue    time.Duration `json: "date_of_issue"`
	DateOfDelivery time.Duration `json: "date_of_delivery"`
}
