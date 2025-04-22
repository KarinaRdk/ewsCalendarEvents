package response

import (
	"encoding/xml"
	"errors"
)

type responseMessage struct {
	ResponseClass string `xml:"ResponseClass,attr"`
	MessageText   string `xml:"MessageText"`
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

func CheckDeleteItemResponseForErrors(bb []byte) error {
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
