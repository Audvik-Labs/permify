package telemetry

import (
	"runtime"
	"time"

	"github.com/rs/xid"
	"go.opentelemetry.io/contrib/instrumentation/host"
	orn "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/attribute"
	omt "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.10.0"
)

// NewNoopMeter - Creates new noop meter
func NewNoopMeter() omt.Meter {
	return omt.NewNoopMeter()
}

// NewMeter - Creates new meter
func NewMeter(exporter metric.Exporter) (omt.Meter, error) {
	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("permify"),
			attribute.String("id", xid.New().String()),
			attribute.String("version", "0.2.1"),
			attribute.String("os", runtime.GOOS),
			attribute.String("arch", runtime.GOARCH),
		)),
	)

	global.SetMeterProvider(mp)

	if err := orn.Start(
		orn.WithMinimumReadMemStatsInterval(time.Second),
		orn.WithMeterProvider(mp),
	); err != nil {
		return nil, err
	}

	if err := host.Start(host.WithMeterProvider(mp)); err != nil {
		return nil, err
	}

	return mp.Meter("permify"), nil
}
