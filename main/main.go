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
		"rwb\\***",
		"",
		&ews.Config{Dump: true, NTLM: false},
	)

	//var layout = time.RFC3339
	//startTime, _ := time.Parse(layout, "2025-04-22T09:00:00.000Z") // создается на 3 часа позже
	//
	//duration := 1 * time.Hour
	//err := ewsutil.CreateEvent(c, []string{"radkevich.karina@rwb.ru"}, []string{},
	//	"Release", "tassssk", "meet", startTime, duration)
	//
	//if err != nil {
	//	log.Fatal("err>: ", err.Error())
	//}
	//
	//fmt.Println("--- success ---")

	update(c)

}

func update(c ews.Client) {
	startTime, _ := time.Parse(time.RFC3339, "2025-04-22T16:10:00.000Z")
	endTime := startTime.Add(2 * time.Hour)

	err := ews.UpdateEvent(
		c,
		"AAMkADk0ZDcxYmQxLWEwODMtNDE5Ny1hOWQ5LTdiZjdmZDM0YWRlYgBGAAAAAABvD3r+tuTDR7aREjHyaftxBwBD5SlB15ajRa+MhxxWrQnNAAAAAAENAABD5SlB15ajRa+MhxxWrQnNAABGKHwDAAA=", // Item ID from created calendar item
		"DwAAABYAAABD5SlB15ajRa+MhxxWrQnNAABGKJgt", // Change Key from the same item
		"Updated taaaaask",                         // New body text
		startTime,
		endTime,
	)

	if err != nil {
		log.Fatal("update failed: ", err)
	}

	fmt.Println("--- update success ---")
}
