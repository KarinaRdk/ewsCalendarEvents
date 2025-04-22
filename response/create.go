package response

import (
	"encoding/xml"
	"errors"
)

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
	ResponseClass string        `xml:"ResponseClass,attr"`
	MessageText   string        `xml:"MessageText"`
	Items         CalendarItems `xml:"Items"`
}

type ItemId struct {
	Id        string `xml:"Id,attr"`
	ChangeKey string `xml:"ChangeKey,attr"`
}

type CalendarItemResponse struct {
	ItemId ItemId `xml:"ItemId"`
}

type CalendarItems struct {
	CalendarItem []CalendarItemResponse `xml:"CalendarItem"`
}

func ParseCreateResponse(bb []byte) (string, string, error) {
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
