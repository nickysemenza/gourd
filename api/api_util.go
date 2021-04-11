package api

import (
	"fmt"
	"math"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"go.opentelemetry.io/otel/trace"
)

func parsePagination(o *OffsetParam, l *LimitParam) ([]db.SearchOption, *Items) {
	offset := 0
	limit := 20
	if o != nil {
		offset = int(*o)
	}
	if l != nil {
		limit = int(*l)
	}
	return []db.SearchOption{db.WithOffset(uint64(offset)), db.WithLimit(uint64(limit))}, &Items{Offset: offset, Limit: limit, PageNumber: (offset/limit + 1)}
}

func (l *Items) setTotalCount(count uint64) {
	c := int(count)
	l.TotalCount = c
	l.PageCount = int(math.Ceil(float64(c) / float64(l.Limit)))
}

func sendErr(ctx echo.Context, code int, err error) error {
	trace.SpanFromContext(ctx.Request().Context()).AddEvent(fmt.Sprintf("error: %s", err))
	return ctx.JSON(code, Error{
		Message: err.Error(),
	})
}
