use gourd::configuration::get_configuration;
use gourd::startup::Application;
use opentelemetry::sdk::{
    propagation::TraceContextPropagator,
    trace::{self, Sampler},
};
use tracing::info;
// use tracing_subscriber::subscribe::CollectExt;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // global::set_text_map_propagator(opentelemetry_jaeger::Propagator::new());
    // not jaeger, this one: https://github.com/openebs/Mayastor/blob/master/control-plane/rest/service/src/main.rs#L64
    opentelemetry::global::set_text_map_propagator(TraceContextPropagator::new());
    let tracer = opentelemetry_jaeger::new_pipeline()
        .with_service_name("gourd-rs")
        .with_trace_config(trace::config().with_sampler(Sampler::AlwaysOn))
        .install_simple()
        .expect("pipeline install error");
    info!("tracer: {:?}", tracer);

    // cf
    // https://github.com/BartWillems/rustfuif/blob/master/src/main.rs#L54
    // https://github.com/cardbox/backend/blob/e92bae636ad4ed0668d609f24dbe36db9387764c/app/src/configure.rs
    // https://github.com/ex-howiebbq/sisi-spss/blob/f3ea31df2620a3ae0a152160539b516b2c670193/stargate/src/telemetry.rs
    let opentelemetry = tracing_opentelemetry::layer().with_tracer(tracer);

    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer())
        .with(opentelemetry)
        .try_init()
        .expect("unable to initialize the tokio tracer");

    let configuration = get_configuration().expect("Failed to read configuration.");
    info!("confiig: {:?}", configuration);
    let application = Application::build(configuration).await?;
    application.run_until_stopped().await?;
    Ok(())
}
