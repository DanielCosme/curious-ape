use actix_web::{web, HttpResponse};
use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

#[derive(serde::Deserialize)]
pub struct FormData {
    name: String,
    description: String,
}

pub async fn create_habit(form: web::Form<FormData>, pg_pool: web::Data<PgPool>) -> HttpResponse {
    match sqlx::query!(
        r#"
            INSERT INTO habits (id, name, description, created_at)
            VALUES ($1, $2, $3, $4)
        "#,
        Uuid::new_v4(),
        form.name,
        form.description,
        Utc::now()
    )
    .execute(pg_pool.get_ref())
    .await
    {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(e) => {
            println!("Failed to execute query: {}", e);
            HttpResponse::InternalServerError().finish()
        }
    }
}
