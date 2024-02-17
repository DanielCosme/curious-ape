use secrecy::ExposeSecret;
use sqlx::PgPool;
use std::net::TcpListener;

use ape::configuration::get_configuration;
use ape::startup::run;
use ape::telemmetry::{get_tracing_subscriber, init_subscriber};

#[tokio::main]
async fn main() -> Result<(), std::io::Error> {
    // Tracing Setup.
    let subscriber = get_tracing_subscriber("Ape".into(), "info".into(), std::io::stdout);
    init_subscriber(subscriber);

    // Parse config.
    let configuration = get_configuration().expect("Failed to read configuraion");

    // Setup Database Connection.
    let conn_pool =
        PgPool::connect_lazy(configuration.database.connection_string().expose_secret())
            .expect("Failed to connect to Postgres.");

    // Setup TCP connection.
    let addr = format!("127.0.0.1:{}", configuration.application_port);
    let listener = TcpListener::bind(addr)?;

    run(listener, conn_pool)?.await
}
