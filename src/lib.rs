use actix_web::dev::Server;
use actix_web::{web, App, HttpRequest, HttpResponse, HttpServer};
use std::net::TcpListener;

#[derive(serde::Deserialize)]
struct FormData {
    name: String,
    description: String,
}

async fn create_habit(_form: web::Form<FormData>) -> HttpResponse {
    HttpResponse::Ok().finish()
}

async fn health_check(_req: HttpRequest) ->  HttpResponse {
    HttpResponse::Ok().finish()
}

pub fn run(listener: TcpListener) -> Result<Server, std::io::Error> {
    let server = HttpServer::new(|| App::new()
        .route("/health_check", web::get().to(health_check))
        .route("/create_habit", web::post().to(create_habit)))
        .listen(listener)?
        .run();

    Ok(server)
}
