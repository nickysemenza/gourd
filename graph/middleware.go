package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// Observability is a Middleware for tracing graphql queries
type Observability struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.FieldInterceptor
	graphql.OperationInterceptor
} = Observability{}

func (c Observability) ExtensionName() string {
	return "Observability"
}

func (c Observability) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (c Observability) InterceptField(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	fc := graphql.GetFieldContext(ctx)

	tr := global.Tracer("graphql")
	ctx, span := tr.Start(ctx, fmt.Sprintf("graphql: field.%s", fc.Field.Name))

	defer span.End()

	field := fc.Field
	span.SetAttributes(
		label.Key("resolver.path").String(fc.Path().String()),
		label.Key("resolver.object").String(field.ObjectDefinition.Name),
		label.Key("resolver.field").String(field.Name),
		label.Key("resolver.alias").String(field.Alias),
	)
	for _, arg := range field.Arguments {
		if arg.Value != nil {
			span.SetAttributes(
				label.Key(fmt.Sprintf("resolver.args.%s", arg.Name)).String(arg.Value.String()),
			)
		}
	}

	errs := graphql.GetErrors(ctx)
	if len(errs) != 0 {
		span.SetStatus(codes.Unknown, errs.Error())
		span.SetAttributes(label.Key("error").Bool(true))
		for i, err := range errs {
			span.SetAttributes(
				label.Key(fmt.Sprintf("resolver.error.%d.message", i)).String(err.Error()),
				label.Key(fmt.Sprintf("resolver.error.%d.kind", i)).String(fmt.Sprintf("%T", err)),
			)
		}
	}

	return next(ctx)
}

// InterceptOperation is adapted from https://github.com/aereal/hibi/blob/master/api/gqlopencensus/census.go
func (c Observability) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)

	span := trace.SpanFromContext(ctx)
	span.SetName(operationName(ctx))

	span.SetAttributes(
		label.String("request.query", oc.RawQuery),
	)

	if stats := extension.GetComplexityStats(ctx); stats != nil {
		span.SetAttributes(label.Key("request.complexity.actual").Int(stats.Complexity))
		span.SetAttributes(label.Key("request.complexity.limit").Int(stats.ComplexityLimit))
	}

	for k, v := range oc.Variables {
		span.SetAttributes(label.Key(fmt.Sprintf("request.variables.%s", k)).String(fmt.Sprintf("%+v", v)))
	}
	span.AddEvent(ctx, "graphql: processing begin")
	span.AddEventWithTimestamp(ctx, oc.Stats.Read.Start, "graphql read: start")
	span.AddEventWithTimestamp(ctx, oc.Stats.Read.Start, "graphql read: end")
	return next(ctx)
}

func operationName(ctx context.Context) string {
	oc := graphql.GetOperationContext(ctx)
	reqName := "nameless-operation"
	if oc.Doc != nil && len(oc.Doc.Operations) != 0 {
		op := oc.Doc.Operations[0]
		if op.Name != "" {
			reqName = op.Name
		}
	}
	return fmt.Sprintf("graphql: %s", reqName)
}
