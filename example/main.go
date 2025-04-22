package main

import (
	"fmt"
	"time"

	"github.com/mhewedy/ews"
)

const layout = time.RFC3339

func main() {

	//  This is how a client might be initialised:
	c := ews.NewClient(
		"https://mail.corp.rwb.ru/EWS/Exchange.asmx",
		"rwb\\login",
		"password",
		&ews.Config{Dump: true, NTLM: false},
	)
	var itemID, changeID string

	itemID, changeID, err := createEvent(c)
	//itemID, changeID, err := updateEvent(c)
	//err := deleteEvent(c)
	if err != nil {
		fmt.Println("err: ", err.Error())
	} else {
		fmt.Println("--- success ---\nitemID:", itemID, "\nchangeID:", changeID)
	}
}

// createEvent shows how to call ews.CreateEvent. On success, it returns both an item ID and a change ID.
// Both are required for updating or deleting the event.
func createEvent(c ews.Client) (string, string, error) {
	start, _ := time.Parse(layout, "2026-04-25T08:15:00.000Z") // создается на 3 часа позже
	finish, _ := time.Parse(layout, "2026-04-25T08:30:00.000Z")

	model := ews.Create{
		User:     "radkevich.karina@rwb.ru",
		Subject:  "Release service",
		Body:     "Test event",
		Location: "https://meet.wb.ru/WBPay_Releases",
		Start:    AdjustToExchangeTimezone(start, 3),
		End:      AdjustToExchangeTimezone(finish, 3),
	}
	return ews.CreateEvent(c, model)
}

// updateEvent demonstrates how to update an event. On success, a new change ID is returned.
// Any subsequent update or deletion of the event requires the latest change ID.
// UpdateEvent allows changing the subject, body, and time of the event.
func updateEvent(client ews.Client) (string, string, error) {
	start, _ := time.Parse(layout, "2025-04-26T12:01:00.000Z") // создается на 3 часа позже
	finish, _ := time.Parse(layout, "2025-04-26T12:10:00.000Z")

	model := ews.Update{
		ItemID:    "AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwWAAA=",
		ChangeKey: "DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKKtQ",
		Subject:   "Release keeper",
		Body:      "Updated event",
		Start:     AdjustToExchangeTimezone(start, 3),
		End:       AdjustToExchangeTimezone(finish, 3),
	}
	itemID, changeID, err := ews.UpdateEvent(client, model)

	return itemID, changeID, err
}

// deleteEvent is an example of how to delete a specific calendar event.
func deleteEvent(c ews.Client) error {
	return ews.DeleteEvent(c,
		"AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwWAAA=",
		"DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKKtw")
}

func AdjustToExchangeTimezone(t time.Time, hoursOffset int) time.Time {
	return t.Add(time.Duration(-hoursOffset) * time.Hour)
}
