use opentelemetry::sdk::{self, propagation::TraceContextPropagator, trace::Sampler};
use tracing::info;
use tracing_subscriber::{
    prelude::__tracing_subscriber_SubscriberExt, util::SubscriberInitExt, EnvFilter,
};

pub fn initialize_tracing(s: &str) {
    // global::set_text_map_propagator(opentelemetry_jaeger::Propagator::new());
    // not jaeger, this one: https://github.com/openebs/Mayastor/blob/master/control-plane/rest/service/src/main.rs#L64
    opentelemetry::global::set_text_map_propagator(TraceContextPropagator::new());
    let tracer = opentelemetry_jaeger::new_agent_pipeline()
        .with_service_name(s)
        .with_trace_config(sdk::trace::config().with_sampler(Sampler::AlwaysOn))
        .install_simple()
        .expect("pipeline install error");

    info!("tracer: {:?}", tracer);

    // cf
    // https://github.com/BartWillems/rustfuif/blob/master/src/main.rs#L54
    // https://github.com/cardbox/backend/blob/e92bae636ad4ed0668d609f24dbe36db9387764c/app/src/configure.rs
    // https://github.com/ex-howiebbq/sisi-spss/blob/f3ea31df2620a3ae0a152160539b516b2c670193/stargate/src/telemetry.rs
    let opentelemetry = tracing_opentelemetry::layer().with_tracer(tracer);

    // let console_layer = console_subscriber::spawn();

    tracing_subscriber::registry()
        // .with(console_layer)
        .with(tracing_subscriber::fmt::layer())
        .with(EnvFilter::from_default_env())
        .with(opentelemetry)
        .try_init()
        .expect("unable to initialize the tokio tracer");
}
