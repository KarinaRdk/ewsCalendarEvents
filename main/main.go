package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mhewedy/ews"
)

func main() {
	c := ews.NewClient(
		"https://mail.corp.rwb.ru/EWS/Exchange.asmx",
		"*",
		"*",
		&ews.Config{Dump: true, NTLM: false},
	)
	//
	//var layout = time.RFC3339
	//startTime, _ := time.Parse(layout, "2025-04-23T09:00:00.000Z") // создается на 3 часа позже
	//AdjustToExchangeTimezone(startTime, 3)
	//duration := 1 * time.Hour
	//err := ewsutil.CreateEvent(c, []string{"radkevich.karina@rwb.ru"}, []string{},
	//	"Release", "created task", "meet", AdjustToExchangeTimezone(startTime, 3), duration)
	//
	//if err != nil {
	//	log.Fatal("err>: ", err.Error())
	//}

	//fmt.Println("--- success ---")

	//update(c)
	Delete(c)

}

func update(c ews.Client) {
	originalStartTime, _ := time.Parse(time.RFC3339, "2025-04-23T15:30:00.000Z")
	adjustedStart := AdjustToExchangeTimezone(originalStartTime, 3)
	originalEndTime, _ := time.Parse(time.RFC3339, "2025-04-23T16:30:00.000Z")
	adjustedEndTime := AdjustToExchangeTimezone(originalEndTime, 3)
	err := ews.UpdateEvent(
		c,
		"AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwGAAA=", // Item ID from created calendar item
		"DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKJ6v", // Change Key from the same item changes with each update
		"Set new again", // New body text
		adjustedStart,
		adjustedEndTime,
	)

	if err != nil {
		log.Fatal("update failed: ", err)
	}

	fmt.Println("--- update success ---")
}
func AdjustToExchangeTimezone(t time.Time, hoursOffset int) time.Time {
	return t.Add(time.Duration(-hoursOffset) * time.Hour)
}

func Delete(c ews.Client) {
	err := ews.DeleteEvent(c,
		"AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwGAAA=", // Item ID from created calendar item
		"DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKJ6v")

	if err != nil {
		log.Fatal("update failed: ", err)
	}

	fmt.Println("--- update success ---")
}
