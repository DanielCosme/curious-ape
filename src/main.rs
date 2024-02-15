use std::net::TcpListener;

use ape::configuration::get_configuration;
use ape::startup::run;

#[tokio::main]
async fn main() -> Result<(), std::io::Error> {
    let configuration = get_configuration().expect("Failed to read configuraion");
    let addr = format!("127.0.0.1:{}", configuration.application_port);
    let listener = TcpListener::bind(addr)?;

    run(listener)?.await
}
