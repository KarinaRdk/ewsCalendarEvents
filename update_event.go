package ews

import (
	"encoding/xml"
	"errors"
)

//const (
//	ResponseClassSuccess = "Success"
//	ResponseClassWarning = "Warning"
//	ResponseClassError   = "Error"
//)

type UpdateItem struct {
	XMLName                               xml.Name     `xml:"m:UpdateItem"`
	Xmlns                                 string       `xml:"xmlns,attr"`
	MessageDisposition                    string       `xml:"MessageDisposition,attr"`
	ConflictResolution                    string       `xml:"ConflictResolution,attr"`
	SendMeetingInvitationsOrCancellations string       `xml:"SendMeetingInvitationsOrCancellations,attr"`
	ItemChanges                           []ItemChange `xml:"m:ItemChanges>t:ItemChange"`
}

type DeleteItem struct {
	XMLName                  xml.Name `xml:"m:DeleteItem"`
	Xmlns                    string   `xml:"xmlns:m,attr"`
	DeleteType               string   `xml:"DeleteType,attr"`
	SendMeetingCancellations string   `xml:"SendMeetingCancellations,attr"`
	AffectedTaskOccurrences  string   `xml:"AffectedTaskOccurrences,attr"`
	PerformReminderAction    bool     `xml:"PerformReminderAction,attr"`
	ItemIds                  []ItemId `xml:"m:ItemIds>t:ItemId"`
}
type deleteItemResponseBodyEnvelop struct {
	XMLName xml.Name                   `xml:"Envelope"`
	Body    deleteItemResponseBodyBody `xml:"Body"`
}

type deleteItemResponseBodyBody struct {
	DeleteItemResponse deleteItemResponse `xml:"DeleteItemResponse"`
}

type deleteItemResponse struct {
	ResponseMessages deleteItemResponseMessages `xml:"ResponseMessages"`
}

type deleteItemResponseMessages struct {
	DeleteItemResponseMessage responseMessage `xml:"DeleteItemResponseMessage"`
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
	FieldURI     FieldURI      `xml:"t:FieldURI"`
	CalendarItem *CalendarItem `xml:"t:CalendarItem,omitempty"`
	Message      *Message      `xml:"t:Message,omitempty"`
}

type AppendToItemField struct {
	FieldURI FieldURI `xml:"t:FieldURI"`
	Message  *Message `xml:"t:Message"`
}

type FieldURI struct {
	FieldURI string `xml:"FieldURI,attr"`
}

//type CalendarItem struct {
//	XMLName xml.Name   `xml:"t:CalendarItem"`
//	Start   *time.Time `xml:"t:Start,omitempty"`
//	End     *time.Time `xml:"t:End,omitempty"`
//}

//type ItemId struct {
//	Id        string `xml:"Id,attr"`
//	ChangeKey string `xml:"ChangeKey,attr"`
//}

//type Message struct {
//	Body Body `xml:"t:Body"`
//}
//
//type Body struct {
//	BodyType string `xml:"BodyType,attr"`
//	Value    string `xml:",chardata"`
//}

//	type Mailbox struct {
//		EmailAddress string `xml:"t:EmailAddress"`
//	}

func checkUpdateItemResponseForErrors(resp []byte) error {
	var soapResp struct {
		Body struct {
			UpdateItemResponse struct {
				ResponseMessages struct {
					ResponseMessage struct {
						ResponseClass string `xml:"ResponseClass,attr"`
						MessageText   string `xml:"MessageText"`
					} `xml:"UpdateItemResponseMessage"`
				} `xml:"ResponseMessages"`
			} `xml:"UpdateItemResponse"`
		} `xml:"Body"`
	}
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return err
	}
	r := soapResp.Body.UpdateItemResponse.ResponseMessages.ResponseMessage
	if r.ResponseClass == "Error" {
		return errors.New(r.MessageText)
	}
	return nil
}

//func UpdateEvent(c Client, itemID, changeKey, newBody string, newStart, newEnd time.Time) error {
//	update := &UpdateItem{
//		Xmlns:                                 "http://schemas.microsoft.com/exchange/services/2006/messages",
//		MessageDisposition:                    "SaveOnly",
//		ConflictResolution:                    "AutoResolve",
//		SendMeetingInvitationsOrCancellations: "SendToNone",
//		ItemChanges: []ItemChange{
//			{
//				ItemId: ItemId{
//					Id:        itemID,
//					ChangeKey: changeKey,
//				},
//				Updates: []SetItemField{
//					{
//						FieldURI: FieldURI{FieldURI: "calendar:Start"},
//						CalendarItem: &CalendarItem{
//							Start: &newStart,
//						},
//					},
//					{
//						FieldURI: FieldURI{FieldURI: "calendar:End"},
//						CalendarItem: &CalendarItem{
//							End: &newEnd,
//						},
//					},
//					{
//						FieldURI: FieldURI{FieldURI: "item:Body"},
//						Message: &Message{
//							Body: Body{
//								BodyType: "Text",
//								Value:    newBody,
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//
//	xmlBytes, err := xml.MarshalIndent(update, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	resp, err := c.SendAndReceive(xmlBytes)
//	if err != nil {
//		return err
//	}
//
//	return checkUpdateItemResponseForErrors(resp)
//}

func DeleteEvent(c Client, itemID, changeKey string) error {
	deleteReq := &DeleteItem{
		Xmlns:                    "http://schemas.microsoft.com/exchange/services/2006/messages",
		DeleteType:               "HardDelete",
		SendMeetingCancellations: "SendToNone",
		AffectedTaskOccurrences:  "AllOccurrences",
		PerformReminderAction:    false,
		ItemIds: []ItemId{
			{
				Id:        itemID,
				ChangeKey: changeKey,
			},
		},
	}

	xmlBytes, err := xml.MarshalIndent(deleteReq, "", "  ")
	if err != nil {
		return err
	}

	resp, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return err
	}

	return checkDeleteItemResponseForErrors(resp)
}

func checkDeleteItemResponseForErrors(bb []byte) error {
	var soapResp deleteItemResponseBodyEnvelop
	if err := xml.Unmarshal(bb, &soapResp); err != nil {
		return err
	}

	resp := soapResp.Body.DeleteItemResponse.ResponseMessages.DeleteItemResponseMessage
	if resp.ResponseClass == "Error" {
		return errors.New(resp.MessageText)
	}
	return nil
}
