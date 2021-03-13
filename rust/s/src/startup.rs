use crate::configuration::{DatabaseSettings, Settings};
// use crate::email_clhealth_checkient::EmailClient;
use crate::routes::parser;
use actix_web::web::Data;
use actix_web::{dev::Server, web};
use actix_web::{App, HttpServer};
// use actix_web_opentelemetry::RequestTracing;
use sqlx::postgres::PgPoolOptions;
use sqlx::PgPool;
use std::net::TcpListener;
// use tracing_actix_web::TracingLogger;

pub struct Application {
    port: u16,
    server: Server,
}

impl Application {
    pub async fn build(configuration: Settings) -> Result<Self, std::io::Error> {
        let connection_pool = get_connection_pool(&configuration.database)
            .await
            .expect("Failed to connect to Postgres.");

        // let sender_email = configuration
        //     .email_client
        //     .sender()
        //     .expect("Invalid sender email address.");
        // let email_client = EmailClient::new(
        //     configuration.email_client.base_url,
        //     sender_email,
        //     configuration.email_client.authorization_token,
        // );

        let address = format!(
            "{}:{}",
            configuration.application.host, configuration.application.port
        );
        let listener = TcpListener::bind(&address)?;
        let port = listener.local_addr().unwrap().port();
        let server = run(listener, connection_pool)?;

        Ok(Self { port, server })
    }

    pub fn port(&self) -> u16 {
        self.port
    }

    pub async fn run_until_stopped(self) -> Result<(), std::io::Error> {
        self.server.await
    }
}

pub async fn get_connection_pool(configuration: &DatabaseSettings) -> Result<PgPool, sqlx::Error> {
    PgPoolOptions::new()
        .connect_timeout(std::time::Duration::from_secs(2))
        .connect_with(configuration.with_db())
        .await
}

fn run(
    listener: TcpListener,
    db_pool: PgPool,
    // email_client: EmailClient,
) -> Result<Server, std::io::Error> {
    let db_pool = Data::new(db_pool);

    // let email_client = Data::new(email_client);
    let server = HttpServer::new(move || {
        App::new()
            // .wrap(RequestTracing::new())
            .route("/parse", web::get().to(parser))
            .service(web::resource("/parse2").route(web::get().to(parser::index)))
            .app_data(db_pool.clone())
        // .app_data(email_client.clone())
    })
    .listen(listener)?
    .run();
    Ok(server)
}
