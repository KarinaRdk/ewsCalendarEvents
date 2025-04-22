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
		"rwb\\NN_0ki",
		"90019767711917Lf*",
		&ews.Config{Dump: true, NTLM: false},
	)
	var itemID, changeID string

	//itemID, changeID, err := createEvent(c)
	itemID, changeID, err := updateEvent(c)
	if err != nil {
		fmt.Println("err: ", err.Error())
	} else {
		fmt.Println("--- success ---\nitemID:", itemID, "\nchangeID:", changeID)
	}
}

// createEvent is an example of how to call ews.CreateEvent. In case of success an item ID and a change ID are returned
// they are both required for update or deletion of the event
func createEvent(c ews.Client) (string, string, error) {
	start, _ := time.Parse(layout, "2025-04-25T08:15:00.000Z") // создается на 3 часа позже
	finish, _ := time.Parse(layout, "2025-04-25T08:30:00.000Z")

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

// updateEvent is an example of how update of event might be evoked. In case of success a new change ID is returned
// a subsequent update or deletion of the event require the most recent change ID. UpdateEvent allows to change
// subject, body, and time of the event
func updateEvent(client ews.Client) (string, string, error) {
	start, _ := time.Parse(layout, "2025-04-25T10:01:00.000Z") // создается на 3 часа позже
	finish, _ := time.Parse(layout, "2025-04-25T10:10:00.000Z")

	model := ews.Update{
		ItemID:    "AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwVAAA=",
		ChangeKey: "DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKKmz",
		Subject:   "release KEEPER",
		Body:      "Updated again body",
		Start:     AdjustToExchangeTimezone(start, 3),
		End:       AdjustToExchangeTimezone(finish, 3),
	}
	itemID, changeID, err := ews.UpdateEvent(client, model)

	return itemID, changeID, err
}

// deleteEvent is an example of how to call deletion of a specific calendar event:
func deleteEvent(c ews.Client) error {
	return ews.DeleteEvent(c,
		"AQMkAGNiYzYzODU3LThjNGEtNDM1NQAtYjEyMC01YzlhMTVjYwIyOQBGAAADcZyuoZlofUKjE+TTcoFoDAcATh2JFwMGt0+69GVJ8SYjsQAAAgENAAAATh2JFwMGt0+69GVJ8SYjsQAAAhHXAAAA",
		"DwAAABYAAABOHYkXAwa3T7r0ZUnxJiOxAAAAABMO")
}

func AdjustToExchangeTimezone(t time.Time, hoursOffset int) time.Time {
	return t.Add(time.Duration(-hoursOffset) * time.Hour)
}
