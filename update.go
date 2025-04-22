package ews

import (
	"encoding/xml"
	"time"

	"github.com/mhewedy/ews/response"
)

type UpdateItem struct {
	XMLName                               xml.Name     `xml:"m:UpdateItem"`
	Xmlns                                 string       `xml:"xmlns,attr"`
	MessageDisposition                    string       `xml:"MessageDisposition,attr"`
	ConflictResolution                    string       `xml:"ConflictResolution,attr"`
	SendMeetingInvitationsOrCancellations string       `xml:"SendMeetingInvitationsOrCancellations,attr"`
	ItemChanges                           []ItemChange `xml:"m:ItemChanges>t:ItemChange"`
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type ItemChange struct {
	ItemId  ItemId         `xml:"t:ItemId"`
	Updates []SetItemField `xml:"t:Updates>t:SetItemField"`
}

type ItemUpdate struct {
	SetItemField      *SetItemField      `xml:"t:SetItemField,omitempty"`
	AppendToItemField *AppendToItemField `xml:"t:AppendToItemField,omitempty"`
}

type SetItemField struct {
	FieldURI     FieldURI             `xml:"t:FieldURI"`
	CalendarItem *UpdatedCalendarItem `xml:"t:CalendarItem,omitempty"`
	Message      *Message             `xml:"t:Message,omitempty"`
	Subject      *Subject             `xml:"t:Subject,omitempty"`
}

type AppendToItemField struct {
	FieldURI FieldURI `xml:"t:FieldURI"`
	Message  *Message `xml:"t:Message"`
}

type FieldURI struct {
	FieldURI string `xml:"FieldURI,attr"`
}

type UpdatedCalendarItem struct {
	XMLName xml.Name   `xml:"t:CalendarItem"`
	Start   *time.Time `xml:"t:Start,omitempty"`
	End     *time.Time `xml:"t:End,omitempty"`
	Subject string     `xml:"t:Subject,omitempty"`
}

type Message struct {
	Body Body `xml:"t:Body"`
}

type Body struct {
	BodyType string `xml:"BodyType,attr"`
	Body     string `xml:",chardata"`
}

type Subject struct {
	XMLName xml.Name `xml:"t:Subject"`
	Body    string   `xml:",chardata"`
}

type Update struct {
	ItemID    string
	ChangeKey string
	Subject   string
	Body      string
	Start     time.Time
	End       time.Time
}

func UpdateEvent(c Client, model Update) (string, string, error) {
	update := UpdateItem{
		Xmlns:                                 "http://schemas.microsoft.com/exchange/services/2006/messages",
		MessageDisposition:                    "SaveOnly",
		ConflictResolution:                    "AutoResolve",
		SendMeetingInvitationsOrCancellations: "SendToNone",
		ItemChanges: []ItemChange{
			{
				ItemId: ItemId{
					Id:        model.ItemID,
					ChangeKey: model.ChangeKey,
				},
				Updates: []SetItemField{
					{
						FieldURI: FieldURI{FieldURI: "calendar:Start"},
						CalendarItem: &UpdatedCalendarItem{
							Start: &model.Start,
						},
					},
					{
						FieldURI: FieldURI{FieldURI: "calendar:End"},
						CalendarItem: &UpdatedCalendarItem{
							End: &model.End,
						},
					},
					{
						FieldURI: FieldURI{FieldURI: "item:Subject"},
						CalendarItem: &UpdatedCalendarItem{
							Subject: model.Subject,
						},
					},
					{
						FieldURI: FieldURI{FieldURI: "item:Body"},
						Message: &Message{
							Body: Body{
								BodyType: "Text",
								Body:     model.Body,
							},
						},
					},
				},
			},
		},
	}
	return UpdateCalendarItem(c, update)
}

func UpdateCalendarItem(c Client, i UpdateItem) (string, string, error) {
	xmlBytes, err := xml.MarshalIndent(&i, "", "  ")
	if err != nil {
		return "", "", err
	}

	resp, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return "", "", err
	}

	return response.ParseUpdateResponse(resp)
}
