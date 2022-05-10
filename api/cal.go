package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/labstack/echo/v4"
)

func (a *API) CalTest(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Cal")
	defer span.End()

	res, err := a.Notion.GetAll(ctx, 20, "")
	if err != nil {
		return handleErr(c, err)
	}

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetTzid("America/Los_Angeles")

	for _, r := range res {
		t := *r.Time
		event := ics.NewEvent(r.PageID)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(t)
		event.SetEndAt(t.Add(time.Hour))
		// event.SetEndAt(time.Now())
		event.SetSummary(fmt.Sprintf("%s - %s", r.Title, strings.Join(r.Tags, ", ")))
		// event.SetLocation("Address")
		// event.SetDescription("Description")
		event.SetURL(r.NotionURL)
		event.AddRrule(fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", time.Now().Month(), time.Now().Day()))
		event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
		event.AddAttendee("reciever or participant", ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))

		cal.AddVEvent(event)
	}
	return c.Blob(http.StatusOK, "text/calendar", []byte(cal.Serialize()))

}
