package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Init(serviceName string) (func(context.Context) error, error) {
	println("📊 OpenTelemetry başlatılıyor, servis:", serviceName)

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		println("❌ Telemetry exporter oluşturulamadı:", err.Error())
		return nil, err
	}

	println("✅ Telemetry exporter oluşturuldu")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	println("✅ Tracer provider oluşturuldu")
	otel.SetTracerProvider(tp)

	println("✅ OpenTelemetry başlatıldı")
	return tp.Shutdown, nil
}
