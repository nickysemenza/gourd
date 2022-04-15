package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func parsePagination(o *OffsetParam, l *LimitParam) ([]db.SearchOption, Items) {
	offset := 0
	limit := 20
	if o != nil {
		offset = int(*o)
	}
	if l != nil {
		limit = int(*l)
	}
	items := Items{
		Offset: offset,
		Limit:  limit,
	}
	if limit == 0 {
		items.PageNumber = 0
	} else {
		items.PageNumber = offset/limit + 1
	}

	return []db.SearchOption{db.WithOffset(uint64(offset)), db.WithLimit(uint64(limit))}, items
}

func (l *Items) setTotalCount(count uint64) {
	c := int(count)
	l.TotalCount = c
	l.PageCount = int(math.Ceil(float64(c) / float64(l.Limit)))
}

func sendErr(c echo.Context, code int, err error) error {
	trace.SpanFromContext(c.Request().Context()).AddEvent(fmt.Sprintf("error: %v", err))
	logrus.WithField("code", code).WithField("route", c.Request().URL).Errorf("http err: %v", err)
	return c.JSON(code, Error{
		Message: err.Error(),
	})
}

func handleErr(c echo.Context, err error) error {
	if errors.Is(err, common.ErrNotFound) {
		return sendErr(c, http.StatusNotFound, err)
	}
	return sendErr(c, http.StatusInternalServerError, err)

}
