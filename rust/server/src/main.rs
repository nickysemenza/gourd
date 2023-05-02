use std::env;

use gourd::server::Application;
use gourd::{configuration::get_configuration, usda_loader};
use opentelemetry::sdk;
use opentelemetry::sdk::propagation::TraceContextPropagator;
use opentelemetry::sdk::trace::Sampler;

use tracing::{info, span};
use tracing_subscriber::util::SubscriberInitExt;
use tracing_subscriber::{layer::SubscriberExt, EnvFilter};

extern crate clap;
use clap::{Arg, Command};

#[tokio::main]
async fn main() -> std::io::Result<()> {
    env::set_var("RUST_BACKTRACE", "1");
    let matches = Command::new("gourd CLI")
        .version("1.0")
        .arg(
            Arg::new("config")
                .short('c')
                .long("config")
                .value_name("FILE")
                .help("Sets a custom config file")
                .num_args(1),
        )
        .subcommand(
            Command::new("server").about("runs the http server").arg(
                Arg::new("debug")
                    .short('d')
                    .help("print debug information verbosely"),
            ),
        )
        .subcommand(
            Command::new("load_mappings")
                .about("load unit mappings")
                .arg(
                    Arg::new("MAPPING")
                        .short('m')
                        .long("mapping")
                        .value_name("FILE")
                        .required(true)
                        .help("Sets a custom config file")
                        .num_args(1),
                ),
        )
        .subcommand(Command::new("load_usda").about("load usda"))
        .get_matches();

    // Gets a value for config if supplied by user, or defaults to "default.conf"
    let _config = matches
        .get_one::<String>("config")
        .map(|s| s.as_str())
        .unwrap_or("default.conf");

    let configuration = get_configuration().expect("Failed to read configuration.");
    info!("confiig: {:?}", configuration);

    if matches.subcommand_matches("server").is_some() {
        initialize_tracing("gourd-rs");
        Application::run(configuration).await.unwrap();
    }
    initialize_tracing("gourd-cli");
    if let Some(_m) = matches.subcommand_matches("load_mappings") {
        let root = span!(tracing::Level::TRACE, "load_mappings",);
        let _enter = root.enter();
    }
    if let Some(_m) = matches.subcommand_matches("load_usda") {
        let root = span!(tracing::Level::TRACE, "load_mappings",);
        let _enter = root.enter();
        usda_loader::load_json_into_search().await.unwrap();
    }

    Ok(())
}

fn initialize_tracing(s: &str) {
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
