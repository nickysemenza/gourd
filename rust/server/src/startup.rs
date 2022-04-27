use crate::configuration::Settings;
// use crate::email_clhealth_checkient::EmailClient;
use crate::routes::parser;
use actix_web::middleware::Logger;
use actix_web::{dev::Server, web};
use actix_web::{App, HttpServer};
use actix_web_opentelemetry::RequestTracing;
use std::net::TcpListener;
use tracing::info;

pub struct Application {
    port: u16,
    server: Server,
}

impl Application {
    pub async fn build(configuration: Settings) -> Result<Self, std::io::Error> {
        let address = format!(
            "{}:{}",
            configuration.application.host, configuration.application.port
        );
        let listener = TcpListener::bind(&address)?;
        let port = listener.local_addr().unwrap().port();
        let server = run(listener)?;

        println!("running on http://{}", address);
        Ok(Self { port, server })
    }

    pub fn port(&self) -> u16 {
        self.port
    }

    pub async fn run_until_stopped(self) -> Result<(), std::io::Error> {
        self.server.await
    }
}

fn run(listener: TcpListener) -> Result<Server, std::io::Error> {
    info!("starting up");
    let server = HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .wrap(RequestTracing::new())
            .route("/parse", web::get().to(parser))
            .route("/parse_amount", web::get().to(parser::amount_parser))
            .route("/decode_recipe", web::get().to(parser::decode_recipe))
            .route("/convert", web::post().to(parser::convert))
            .route("/scrape", web::get().to(parser::scrape))
            .route("/pans", web::get().to(parser::pans))
            .route("/debug/scrape", web::get().to(parser::debug_scrape))
            .route(
                "/codec/expand",
                web::post().to(parser::expand_compact_to_input),
            )
            .route(
                "/normalize_amount",
                web::post().to(parser::normalize_amount),
            )
    })
    .listen(listener)?
    .run();
    Ok(server)
}
