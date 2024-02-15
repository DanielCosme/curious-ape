use actix_web::{web, HttpResponse};

#[derive(serde::Deserialize)]
pub struct FormData {
    name: String,
    description: String,
}

pub async fn create_habit(_form: web::Form<FormData>) -> HttpResponse {
    HttpResponse::Ok().finish()
}
