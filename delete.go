package ews

import (
	"encoding/xml"

	"github.com/mhewedy/ews/response"
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

	return response.CheckDeleteItemResponseForErrors(resp)
}
