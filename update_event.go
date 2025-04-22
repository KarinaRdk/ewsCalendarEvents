package ews

import (
	"encoding/xml"
	"errors"
	"time"
)

type UpdateItem struct {
	XMLName                               xml.Name     `xml:"m:UpdateItem"`
	Xmlns                                 string       `xml:"xmlns,attr"`
	MessageDisposition                    string       `xml:"MessageDisposition,attr"`
	ConflictResolution                    string       `xml:"ConflictResolution,attr"`
	SendMeetingInvitationsOrCancellations string       `xml:"SendMeetingInvitationsOrCancellations,attr"`
	ItemChanges                           []ItemChange `xml:"m:ItemChanges>t:ItemChange"`
}

type responseMessage struct {
	ResponseClass string `xml:"ResponseClass,attr"`
	MessageText   string `xml:"MessageText"`
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
	FieldURI     FieldURI            `xml:"t:FieldURI"`
	CalendarItem *UpdateCalendarItem `xml:"t:CalendarItem,omitempty"`
	Message      *Message            `xml:"t:Message,omitempty"`
}

type AppendToItemField struct {
	FieldURI FieldURI `xml:"t:FieldURI"`
	Message  *Message `xml:"t:Message"`
}

type FieldURI struct {
	FieldURI string `xml:"FieldURI,attr"`
}

type UpdateCalendarItem struct {
	XMLName xml.Name   `xml:"t:CalendarItem"`
	Start   *time.Time `xml:"t:Start,omitempty"`
	End     *time.Time `xml:"t:End,omitempty"`
}

// Update работает когда есть только body
type Message struct {
	Body Body `xml:"t:Body"`
}

type Body struct {
	BodyType string `xml:"BodyType,attr"`
	Body     string `xml:",chardata"`
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
	update := &UpdateItem{
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
						CalendarItem: &UpdateCalendarItem{
							Start: &model.Start,
						},
					},
					{
						FieldURI: FieldURI{FieldURI: "calendar:End"},
						CalendarItem: &UpdateCalendarItem{
							End: &model.End,
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

	xmlBytes, err := xml.MarshalIndent(update, "", "  ")
	if err != nil {
		return "", "", err
	}

	resp, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return "", "", err
	}

	return parseUpdateItemResponse(resp)
}

func parseUpdateItemResponse(resp []byte) (string, string, error) {
	var soapResp struct {
		Body struct {
			UpdateItemResponse struct {
				ResponseMessages struct {
					ResponseMessage struct {
						ResponseClass string `xml:"ResponseClass,attr"`
						MessageText   string `xml:"MessageText"`
						Items         struct {
							CalendarItem []calendarItemResponse `xml:"CalendarItem"`
						} `xml:"Items"`
					} `xml:"UpdateItemResponseMessage"`
				} `xml:"ResponseMessages"`
			} `xml:"UpdateItemResponse"`
		} `xml:"Body"`
	}

	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return "", "", err
	}

	r := soapResp.Body.UpdateItemResponse.ResponseMessages.ResponseMessage
	if r.ResponseClass == "Error" {
		return "", "", errors.New(r.MessageText)
	}

	if len(r.Items.CalendarItem) == 0 {
		return "", "", errors.New("no CalendarItem returned")
	}

	item := r.Items.CalendarItem[0]
	return item.ItemId.Id, item.ItemId.ChangeKey, nil
}
