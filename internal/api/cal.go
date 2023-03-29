package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

func (a *API) CalTest(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Cal")
	defer span.End()

	res, err := a.getEnhancedNotion(ctx, 20*time.Hour*24, "")
	if err != nil {
		return handleErr(c, err)
	}

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetTzid("America/Los_Angeles")

	for _, each := range res {
		r := each.Notion
		t := *r.Time
		event := ics.NewEvent(r.PageID)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(t)
		event.SetEndAt(t.Add(time.Hour))
		event.SetSummary(fmt.Sprintf("%s - %s", r.Title, strings.Join(r.Tags, ", ")))
		event.SetURL(r.NotionURL)
		event.SetDescription(spew.Sdump(each.Details))
		event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))

		cal.AddVEvent(event)
	}
	return c.Blob(http.StatusOK, "text/calendar", []byte(cal.Serialize()))

}
