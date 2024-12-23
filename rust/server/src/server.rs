use crate::configuration::Settings;
use crate::routes::parser;
use anyhow::Result;
use axum::{
    routing::{get, post},
    Router,
};
use axum_tracing_opentelemetry::middleware::{OtelAxumLayer, OtelInResponseLayer};
use std::net::SocketAddr;

pub struct Application {
    port: u16,
}

impl Application {
    pub async fn run(configuration: Settings) -> Result<()> {
        let app = Router::new()
            .route("/parse_amount", get(parser::amount_parser))
            .route("/decode_recipe", get(parser::decode_recipe))
            .route("/convert", post(parser::convert))
            .route("/scrape", get(parser::scrape))
            .route("/debug/scrape", get(parser::debug_scrape))
            .route("/debug/search_usda", get(crate::usda_loader::search_usda))
            .route("/debug/get_usda", get(crate::usda_loader::get_usda))
            .route("/codec/expand", post(parser::expand_compact_to_input))
            .route("/normalize_amount", post(parser::normalize_amount))
            .route(
                "/index_recipe_detail",
                post(crate::search::index_recipe_detail),
            )
            // include trace context as header into the response
            .layer(OtelInResponseLayer::default())
            //start OpenTelemetry trace on incoming request
            .layer(OtelAxumLayer::default());

        let port = configuration.application.port;
        let addr = SocketAddr::from(([0, 0, 0, 0], port));
        tracing::debug!("listening on {}", addr);
        let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
        axum::serve(listener, app.into_make_service())
            .await
            .unwrap();
        Ok(())
    }
    pub fn port(&self) -> u16 {
        self.port
    }
}
