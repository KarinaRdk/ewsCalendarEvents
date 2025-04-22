package ews

import (
	"encoding/xml"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_marshal_CalendarItem(t *testing.T) {

	attendee := make([]Attendee, 0)
	attendee = append(attendee,
		Attendee{Mailbox: Mailbox{EmailAddress: "User1@example.com"}},
	)
	attendees := make([]Attendees, 0)
	attendees = append(attendees, Attendees{Attendee: attendee})

	start, _ := time.Parse(time.RFC3339, "2025-04-20T14:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2025-04-20T15:00:00Z")

	citem := &CalendarItem{
		Subject: "Release keeper",
		Body: Body{
			BodyType: "Text",
			Body:     "Новые хэндлеры",
		},
		Location: "meet PCI DSS",

		Start:                start,
		End:                  end,
		IsAllDayEvent:        false,
		LegacyFreeBusyStatus: "Busy",
		RequiredAttendees:    attendees,
	}

	xmlBytes, err := xml.MarshalIndent(citem, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, `<CalendarItem>
  <t:Subject>Release keeper</t:Subject>
  <t:Body BodyType="Text">Новые хэндлеры</t:Body>
  <t:Start>2025-04-20T14:00:00Z</t:Start>
  <t:End>2025-04-20T15:00:00Z</t:End>
  <t:IsAllDayEvent>false</t:IsAllDayEvent>
  <t:LegacyFreeBusyStatus>Busy</t:LegacyFreeBusyStatus>
  <t:Location>meet PCI DSS</t:Location>
  <t:RequiredAttendees>
    <t:Attendee>
      <t:Mailbox>
        <t:EmailAddress>User1@example.com</t:EmailAddress>
      </t:Mailbox>
    </t:Attendee>
  </t:RequiredAttendees>
</CalendarItem>`, string(xmlBytes))
}
