package ews

import (
	"encoding/xml"
	"errors"
)

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
