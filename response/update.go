package response

import (
	"encoding/xml"
	"errors"
)

type UpdateItemResponseEnvelope struct {
	XMLName xml.Name               `xml:"Envelope"`
	Body    UpdateItemResponseBody `xml:"Body"`
}

type UpdateItemResponseBody struct {
	UpdateItemResponse UpdateItemResponseContent `xml:"UpdateItemResponse"`
}

type UpdateItemResponseContent struct {
	ResponseMessages UpdateItemResponseMessages `xml:"ResponseMessages"`
}

type UpdateItemResponseMessages struct {
	UpdateItemResponseMessage UpdateItemResponseMessage `xml:"UpdateItemResponseMessage"`
}

type UpdateItemResponseMessage struct {
	ResponseClass string        `xml:"ResponseClass,attr"`
	MessageText   string        `xml:"MessageText"`
	Items         CalendarItems `xml:"Items"`
}

func ParseUpdateResponse(resp []byte) (string, string, error) {
	var soapResp UpdateItemResponseEnvelope
	if err := xml.Unmarshal(resp, &soapResp); err != nil {
		return "", "", err
	}

	respMsg := soapResp.Body.UpdateItemResponse.ResponseMessages.UpdateItemResponseMessage
	if respMsg.ResponseClass == "Error" {
		return "", "", errors.New(respMsg.MessageText)
	}

	if len(respMsg.Items.CalendarItem) == 0 {
		return "", "", errors.New("no CalendarItem returned")
	}

	item := respMsg.Items.CalendarItem[0]
	return item.ItemId.Id, item.ItemId.ChangeKey, nil
}
