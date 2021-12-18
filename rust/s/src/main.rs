use std::env;

use gourd::configuration::get_configuration;
use gourd::startup::Application;
use opentelemetry::sdk::{
    propagation::TraceContextPropagator,
    trace::{self, Sampler},
};
use tracing::{info, span};
// use tracing_subscriber::subscribe::CollectExt;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

extern crate clap;
use clap::{App, Arg, SubCommand};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env::set_var("RUST_BACKTRACE", "1");
    let matches = App::new("gourd CLI")
        .version("1.0")
        .arg(
            Arg::with_name("config")
                .short("c")
                .long("config")
                .value_name("FILE")
                .help("Sets a custom config file")
                .takes_value(true),
        )
        .subcommand(
            SubCommand::with_name("server")
                .about("runs the http server")
                .arg(
                    Arg::with_name("debug")
                        .short("d")
                        .help("print debug information verbosely"),
                ),
        )
        .subcommand(
            SubCommand::with_name("load_mappings")
                .about("load unit mappings")
                .arg(
                    Arg::with_name("MAPPING")
                        .short("m")
                        .long("mapping")
                        .value_name("FILE")
                        .required(true)
                        .help("Sets a custom config file")
                        .takes_value(true),
                ),
        )
        .get_matches();

    // Gets a value for config if supplied by user, or defaults to "default.conf"
    let _config = matches.value_of("config").unwrap_or("default.conf");

    let configuration = get_configuration().expect("Failed to read configuration.");
    info!("confiig: {:?}", configuration);

    if let Some(_) = matches.subcommand_matches("server") {
        initialize_tracing("gourd-rs");
        let application = Application::build(configuration.clone()).await?;
        application.run_until_stopped().await?;
    }
    initialize_tracing("gourd-cli");
    if let Some(_m) = matches.subcommand_matches("load_mappings") {
        let root = span!(tracing::Level::TRACE, "load_mappings",);
        let _enter = root.enter();
    }
    opentelemetry::global::force_flush_tracer_provider();

    Ok(())
}

fn initialize_tracing(s: &str) {
    // global::set_text_map_propagator(opentelemetry_jaeger::Propagator::new());
    // not jaeger, this one: https://github.com/openebs/Mayastor/blob/master/control-plane/rest/service/src/main.rs#L64
    opentelemetry::global::set_text_map_propagator(TraceContextPropagator::new());
    let tracer = opentelemetry_jaeger::new_pipeline()
        .with_service_name(s)
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
}
