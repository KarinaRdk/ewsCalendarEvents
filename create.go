package ews

import (
	"encoding/xml"
	"time"

	"github.com/mhewedy/ews/response"
)

const (
	busyStatus  = "Busy"
	typeText    = "Text"
	sendAndSave = "SendToAllAndSaveCopy" //  necessary to see event in the calendar
)

type CreateItem struct {
	XMLName                struct{}          `xml:"m:CreateItem"`
	MessageDisposition     string            `xml:"MessageDisposition,attr"`
	SendMeetingInvitations string            `xml:"SendMeetingInvitations,attr"`
	SavedItemFolderId      SavedItemFolderId `xml:"m:SavedItemFolderId"`
	Items                  Items             `xml:"m:Items"`
}

type Items struct {
	Message      []Message      `xml:"t:Message"`
	CalendarItem []CalendarItem `xml:"t:CalendarItem"`
}

type SavedItemFolderId struct {
	DistinguishedFolderId DistinguishedFolderId `xml:"t:DistinguishedFolderId"`
}

type DistinguishedFolderId struct {
	Id string `xml:"Id,attr"`
}

type CalendarItem struct {
	Subject              string      `xml:"t:Subject"`
	Body                 Body        `xml:"t:Body"`
	Start                time.Time   `xml:"t:Start"`
	End                  time.Time   `xml:"t:End"`
	IsAllDayEvent        bool        `xml:"t:IsAllDayEvent"`
	LegacyFreeBusyStatus string      `xml:"t:LegacyFreeBusyStatus"`
	Location             string      `xml:"t:Location"`
	RequiredAttendees    []Attendees `xml:"t:RequiredAttendees"`
}

type Mailbox struct {
	EmailAddress string `xml:"t:EmailAddress"`
}

type Attendee struct {
	Mailbox Mailbox `xml:"t:Mailbox"`
}

type Attendees struct {
	Attendee []Attendee `xml:"t:Attendee"`
}

type Create struct {
	User     string
	Subject  string
	Body     string
	Location string
	Start    time.Time
	End      time.Time
}

// CreateEvent outer layer wrapper for sending create event request
func CreateEvent(c Client, model Create) (string, string, error) {

	requiredAttendees := []Attendee{
		{Mailbox: Mailbox{EmailAddress: model.User}},
	}

	m := CalendarItem{
		Subject: model.Subject,
		Body: Body{
			BodyType: typeText,
			Body:     model.Body,
		},
		Start:                model.Start,
		End:                  model.End,
		IsAllDayEvent:        false,
		LegacyFreeBusyStatus: busyStatus,
		Location:             model.Location,
		RequiredAttendees:    []Attendees{{Attendee: requiredAttendees}},
	}

	return CreateCalendarItem(c, m)
}

func CreateCalendarItem(c Client, ci CalendarItem) (string, string, error) {
	item := &CreateItem{
		SendMeetingInvitations: sendAndSave,
		SavedItemFolderId:      SavedItemFolderId{DistinguishedFolderId{Id: "calendar"}},
	}
	item.Items.CalendarItem = append(item.Items.CalendarItem, ci)

	xmlBytes, err := xml.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", "", err
	}

	bb, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return "", "", err
	}

	itemID, changeKey, err := response.ParseCreateResponse(bb)
	if err != nil {
		return "", "", err
	}

	return itemID, changeKey, nil
}
