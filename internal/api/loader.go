package api

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel/trace"
)

type Foo[T any] interface {
	All(ctx context.Context, exec boil.ContextExecutor) (T, error)
	Count(ctx context.Context, exec boil.ContextExecutor) (int64, error)
}

func countAndQuery[T any, V Foo[T]](ctx context.Context, exec boil.ContextExecutor, loader func(mods ...QueryMod) V, orderBy string, mods ...QueryMod) (T, int64, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("db").Start(ctx, "countAndQuery")
	defer span.End()

	var empty T
	items, err := loader(append(mods, OrderBy(orderBy))...).All(ctx, exec)
	if err != nil {
		return empty, 0, err
	}
	count, err := loader(mods...).Count(ctx, exec)
	if err != nil {
		return empty, 0, err
	}
	return items, count, nil
}

func qmWithPagination(base []QueryMod, pagination Items, mods ...QueryMod) []QueryMod {
	filters := []QueryMod{
		Limit(pagination.Limit),
		Offset(pagination.Offset),
	}
	return append(base, append(mods, filters...)...)
}
