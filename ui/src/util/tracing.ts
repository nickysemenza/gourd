import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { ZoneContextManager } from "@opentelemetry/context-zone";
import { JaegerPropagator } from "@opentelemetry/propagator-jaeger";
import { Resource } from "@opentelemetry/resources";
import { diag, DiagConsoleLogger, DiagLogLevel } from "@opentelemetry/api";
import {
  BatchSpanProcessor,
  SimpleSpanProcessor,
} from "@opentelemetry/sdk-trace-base";

const serviceName = "gourd-ui";
const resource = new Resource({ "service.name": serviceName });
const provider = new WebTracerProvider({ resource });

export const registerTracing = (url: string, batch: boolean) => {
  if (url === "") return;
  console.info("enabled tracing", url);
  diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.ERROR);

  const collectorOptions = {
    url, // url is optional and can be omitted - default is http://localhost:4318/v1/traces
    headers: {}, // an optional object containing custom headers to be sent with each request
    concurrencyLimit: 10, // an optional limit on pending requests
  };

  const exporter = new OTLPTraceExporter(collectorOptions);

  // const exporter = new CollectorTraceExporter({
  //   url: url,
  //   // https://github.com/open-telemetry/opentelemetry-js/issues/2321#issuecomment-889861080
  //   headers: {
  //     "Content-Type": "application/json",
  //   },
  // });

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: /.+/,
      }),
    ],
  });
  provider.register({
    contextManager: new ZoneContextManager(),
    propagator: new JaegerPropagator(),
  });

  if (batch) {
    provider.addSpanProcessor(new BatchSpanProcessor(exporter));
  } else {
    provider.addSpanProcessor(new SimpleSpanProcessor(exporter));
  }
};

export const tracer = provider.getTracer(serviceName);
