use ape::telemmetry::{get_tracing_subscriber, init_subscriber};
use once_cell::sync::Lazy;
use secrecy::ExposeSecret;
use sqlx::{Connection, Executor, PgConnection, PgPool};
use std::net::TcpListener;
use uuid::Uuid;

use ape::configuration::{get_configuration, DatabaseSettings};
use ape::startup::run;

static TRACING: Lazy<()> = Lazy::new(|| {
    let default_filter_level = "info".to_string();
    let subscriber_name = "test".to_string();

    if std::env::var("TEST_LOG").is_ok() {
        let subscriber =
            get_tracing_subscriber(subscriber_name, default_filter_level, std::io::stdout);
        init_subscriber(subscriber);
    } else {
        let subscriber =
            get_tracing_subscriber(subscriber_name, default_filter_level, std::io::sink);
        init_subscriber(subscriber);
    }
});

pub struct TestApp {
    pub addr: String,
    pub db_pool: PgPool,
}

async fn spawn_app() -> TestApp {
    Lazy::force(&TRACING);

    let tcp_listener = TcpListener::bind("127.0.0.1:0").expect("Failed to bind random port");
    let port = tcp_listener.local_addr().unwrap().port();
    let addr = format!("http://127.0.0.1:{}", port);

    let mut config = get_configuration().expect("Failed to read configuration");
    config.database.database_name = Uuid::new_v4().to_string();

    let db_pool = configure_database(&config.database).await;

    let server = run(tcp_listener, db_pool.clone()).expect("Failed to bind address");
    let _ = tokio::spawn(server);
    TestApp { addr, db_pool }
}

pub async fn configure_database(config: &DatabaseSettings) -> PgPool {
    let mut conn = PgConnection::connect(&config.connection_string_without_db().expose_secret())
        .await
        .expect("Failed to connect to Postgres");

    conn.execute(format!(r#"CREATE DATABASE "{}";"#, config.database_name).as_str())
        .await
        .expect("Failed to create database.");

    let conn_pool = PgPool::connect(&config.connection_string().expose_secret())
        .await
        .expect("Failed to connect to Postgress");

    sqlx::migrate!("./migrations")
        .run(&conn_pool)
        .await
        .expect("Failed to migrate the database");

    conn_pool
}

#[tokio::test]
async fn health_check_works() {
    let app = spawn_app().await;
    let client = reqwest::Client::new();

    let response = client
        .get(app.addr + "/health_check")
        .send()
        .await
        .expect("Failed to excecute request.");

    assert!(response.status().is_success());
    assert_eq!(Some(0), response.content_length());
}

#[tokio::test]
async fn create_habit_returns_200_for_valid_form_data() {
    let app = spawn_app().await;
    let client = reqwest::Client::new();

    let body = "name=Wake%20Up&description=wake-up";
    let response = client
        .post(&format!("{}/create_habit", &app.addr))
        .header("Content-Type", "application/x-www-form-urlencoded")
        .body(body)
        .send()
        .await
        .expect("Failed to execute request.");

    assert_eq!(200, response.status().as_u16());

    let saved = sqlx::query!("SELECT name, description FROM habits",)
        .fetch_one(&app.db_pool)
        .await
        .expect("Failed to fetch saved subscription");

    assert_eq!(saved.name, "Wake Up");
    assert_eq!(saved.description.unwrap(), "wake-up");
}

#[tokio::test]
async fn create_habit_returns_400_when_data_is_missing() {
    let app = spawn_app().await;
    let client = reqwest::Client::new();
    let test_cases = vec![("name=", "missing name")];

    for (invalid_body, error_message) in test_cases {
        let response = client
            .post(&format!("{}/create_habit", app.addr))
            .header("Content-Type", "application/x-www-form-urlencoded")
            .body(invalid_body)
            .send()
            .await
            .expect("Failed to execute request.");

        assert_eq!(
            400,
            response.status().as_u16(),
            "The API did not fail with 400 Bad Request when the payload was {}.",
            error_message
        )
    }
}
