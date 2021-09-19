use gourd::configuration::get_configuration;
use gourd::startup::Application;
use opentelemetry::{
    global,
    sdk::{
        propagation::TraceContextPropagator,
        trace::{self, Sampler},
    },
};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // global::set_text_map_propagator(opentelemetry_jaeger::Propagator::new());
    // not jaeger, this one: https://github.com/openebs/Mayastor/blob/master/control-plane/rest/service/src/main.rs#L64
    global::set_text_map_propagator(TraceContextPropagator::new());
    let tracer = opentelemetry_jaeger::new_pipeline()
        .with_service_name("actix_server")
        .with_trace_config(trace::config().with_sampler(Sampler::AlwaysOn))
        .install_simple()
        .expect("pipeline install error");
    println!("tracer: {:?}", tracer);

    let configuration = get_configuration().expect("Failed to read configuration.");
    println!("confiig {:?}", configuration);
    let application = Application::build(configuration).await?;
    application.run_until_stopped().await?;
    Ok(())
}
