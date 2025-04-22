package ews

import (
	"encoding/xml"
	"errors"
	"time"
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

type Message struct {
	ItemClass    string     `xml:"t:ItemClass"`
	Subject      string     `xml:"t:Subject"`
	Body         Body       `xml:"t:Body"`
	Sender       OneMailbox `xml:"t:Sender"`
	ToRecipients XMailbox   `xml:"t:ToRecipients"`
}

type CalendarItem struct {
	Subject                    string      `xml:"t:Subject"`
	Body                       Body        `xml:"t:Body"`
	ReminderIsSet              bool        `xml:"t:ReminderIsSet"`
	ReminderMinutesBeforeStart int         `xml:"t:ReminderMinutesBeforeStart"`
	Start                      time.Time   `xml:"t:Start"`
	End                        time.Time   `xml:"t:End"`
	IsAllDayEvent              bool        `xml:"t:IsAllDayEvent"`
	LegacyFreeBusyStatus       string      `xml:"t:LegacyFreeBusyStatus"`
	Location                   string      `xml:"t:Location"`
	RequiredAttendees          []Attendees `xml:"t:RequiredAttendees"`
	OptionalAttendees          []Attendees `xml:"t:OptionalAttendees"`
	Resources                  []Attendees `xml:"t:Resources"`
}

type Body struct {
	BodyType string `xml:"BodyType,attr"`
	Body     []byte `xml:",chardata"`
}

type OneMailbox struct {
	Mailbox Mailbox `xml:"t:Mailbox"`
}

type XMailbox struct {
	Mailbox []Mailbox `xml:"t:Mailbox"`
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
type createItemResponseBodyEnvelop struct {
	XMLName xml.Name                      `xml:"Envelope"`
	Body    createItemResponseBodyContent `xml:"Body"`
}

type createItemResponseBodyContent struct {
	CreateItemResponse createItemResponse `xml:"CreateItemResponse"`
}

type createItemResponse struct {
	ResponseMessages createItemResponseMessages `xml:"ResponseMessages"`
}

type createItemResponseMessages struct {
	CreateItemResponseMessage createItemResponseMessage `xml:"CreateItemResponseMessage"`
}

type createItemResponseMessage struct {
	ResponseClass string `xml:"ResponseClass,attr"`
	MessageText   string `xml:"MessageText"`
	Items         items  `xml:"Items"`
}

type items struct {
	CalendarItem []calendarItemResponse `xml:"CalendarItem"`
}

type calendarItemResponse struct {
	ItemId ItemId `xml:"ItemId"`
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

// CreateMessageItem
// https://docs.microsoft.com/en-us/exchange/client-developer/web-service-reference/createitem-operation-email-message
func CreateMessageItem(c Client, m ...Message) error {

	item := &CreateItem{
		MessageDisposition: "SendAndSaveCopy",
		SavedItemFolderId:  SavedItemFolderId{DistinguishedFolderId{Id: "sentitems"}},
	}
	item.Items.Message = append(item.Items.Message, m...)

	xmlBytes, err := xml.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}

	bb, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return err
	}

	if err := checkCreateItemResponseForErrors(bb); err != nil {
		return err
	}

	return nil
}

// CreateCalendarItem
// https://docs.microsoft.com/en-us/exchange/client-developer/web-service-reference/createitem-operation-calendar-item
func CreateCalendarItem(c Client, ci CalendarItem) (string, string, error) {
	item := &CreateItem{
		SendMeetingInvitations: "SendToAllAndSaveCopy",
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

	itemID, changeKey, err := parseCreateItemResponse(bb)
	if err != nil {
		return "", "", err
	}

	return itemID, changeKey, nil
}

func checkCreateItemResponseForErrors(bb []byte) error {
	var soapResp createItemResponseBodyEnvelop
	if err := xml.Unmarshal(bb, &soapResp); err != nil {
		return err
	}

	resp := soapResp.Body.CreateItemResponse.ResponseMessages.CreateItemResponseMessage
	if resp.ResponseClass == "Error" {
		return errors.New(resp.MessageText)
	}
	return nil
}

func parseCreateItemResponse(bb []byte) (string, string, error) {
	var soapResp createItemResponseBodyEnvelop
	if err := xml.Unmarshal(bb, &soapResp); err != nil {
		return "", "", err
	}

	resp := soapResp.Body.CreateItemResponse.ResponseMessages.CreateItemResponseMessage
	if resp.ResponseClass == "Error" {
		return "", "", errors.New(resp.MessageText)
	}

	if len(resp.Items.CalendarItem) == 0 {
		return "", "", errors.New("no CalendarItem returned")
	}

	item := resp.Items.CalendarItem[0]
	return item.ItemId.Id, item.ItemId.ChangeKey, nil
}
