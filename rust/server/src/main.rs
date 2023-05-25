mod obs;

use std::env;

use gourd::server::Application;
use gourd::{configuration::get_configuration, usda_loader};

use tracing::{info, span};

extern crate clap;
use clap::{Arg, Command};

use crate::obs::initialize_tracing;

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
