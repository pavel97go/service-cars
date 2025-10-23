package tracing

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.1"
)

func Init(ctx context.Context, serviceName string, endpoint string) (*sdktrace.TracerProvider, error) {
	if serviceName == "" {
		serviceName = "cars-service"
	}
	if endpoint == "" {
		if env := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); env != "" {
			endpoint = env
		} else {
			endpoint = "localhost:4318"
		}
	}

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
func Middleware() fiber.Handler {
	tr := otel.Tracer("http")
	return func(c *fiber.Ctx) error {
		ctx, span := tr.Start(c.Context(), c.Method()+" "+c.Route().Path)
		defer span.End()
		c.SetUserContext(ctx)
		return c.Next()
	}
}
