use sqlx::{Connection, PgConnection};
use std::net::TcpListener;

use ape::{configuration::get_configuration, startup};

#[tokio::test]
async fn health_check_works() {
    let addr = spawn_app().await;
    let client = reqwest::Client::new();

    let response = client
        .get(addr + "/health_check")
        .send()
        .await
        .expect("Failed to excecute request.");

    assert!(response.status().is_success());
    assert_eq!(Some(0), response.content_length());
}

#[tokio::test]
async fn create_habit_returns_200_for_valid_form_data() {
    let addr = spawn_app().await;
    let config = get_configuration().expect("Failed to read configuration");
    let connection_string = config.database.connection_string();
    let mut connection = PgConnection::connect(&connection_string)
        .await
        .expect("Failed to connect to Postgress.");
    let client = reqwest::Client::new();

    let body = "name=Wake%20Up&description=wake-up";
    let response = client
        .post(&format!("{}/create_habit", &addr))
        .header("Content-Type", "application/x-www-form-urlencoded")
        .body(body)
        .send()
        .await
        .expect("Failed to execute request.");

    assert_eq!(200, response.status().as_u16());

    let saved = sqlx::query!("SELECT name, description FROM habits",)
        .fetch_one(&mut connection)
        .await
        .expect("Failed to fetch saved subscription");

    assert_eq!(saved.name, "Wake Up");
    assert_eq!(saved.description.unwrap(), "wake-up");
}

#[tokio::test]
async fn create_habit_returns_400_when_data_is_missing() {
    let addr = spawn_app().await;
    let client = reqwest::Client::new();
    let test_cases = vec![("name=", "missing name")];

    for (invalid_body, error_message) in test_cases {
        let response = client
            .post(&format!("{}/create_habit", &addr))
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

async fn spawn_app() -> String {
    let listener = TcpListener::bind("127.0.0.1:0").expect("Failed to bind random port");
    let port = listener.local_addr().unwrap().port();
    let config = get_configuration().expect("Failed to read configuration");
    let connection_string = config.database.connection_string();
    let mut connection = PgConnection::connect(&connection_string)
        .await
        .expect("Failed to connect to Postgress.");

    let server = startup::run(listener, connection).expect("Failed to bind address");
    let _ = tokio::spawn(server);
    format!("http://127.0.0.1:{}", port)
}

// [Curious Ape - Automated habit tracker]
//
// As Daniel
// I want to create a habit,
// So that I can track it.
//    POST /create_habit.
//    Collect Data From HTML.
//    Parse the request body of a POST request.
//    Libraries to work with a database.
//    Setup Migrations for the database.
//    Get a database connection on API request handlers.
//    Test side effects in the integration tests.
//    Avoid "weird" ineractions between tests when working with a database.
//
//    What information do I need to create a new habit?
//          Post /create_habit
//          Habit
//          - name -> must exist
//          - description -> optional
//
//          Post /create_habit_log ? name/id
//          Habit Log
//          - date - YYYY-MM-DD -> must exist and be valid.
//          - status - Done :: Not done :: No Info -> Must be: Done :: Not Done
//
//          Habit Event/Source
//          - done - yes :: no
//          - is_automated: true :: false
//          - origin: "fitness_record" -> internal resource connection.
//          - provider: "fitbit" -> external place
//
//
// As Daniel,
// I want to set a habit as done or not done.
// So that I can audit my life.
