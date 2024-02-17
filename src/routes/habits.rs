use actix_web::{web, HttpResponse};
use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

#[derive(serde::Deserialize)]
pub struct FormData {
    name: String,
    description: String,
}

#[tracing::instrument(
    name = "Creating a new Habit",
    skip(form, pg_pool),
    fields(
        habit_name = %form.name,
        habit_description = %form.description,
    )
)]
pub async fn create_habit(form: web::Form<FormData>, pg_pool: web::Data<PgPool>) -> HttpResponse {
    match insert_habit(&pg_pool, &form).await {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

#[tracing::instrument(name = "Saving a new habit in the database", skip(form, pool))]
pub async fn insert_habit(pool: &PgPool, form: &FormData) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO habits (id, name, description, created_at)
            VALUES ($1, $2, $3, $4)
        "#,
        Uuid::new_v4(),
        form.name,
        form.description,
        Utc::now()
    )
    .execute(pool)
    .await
    .map_err(|e| {
        tracing::error!("Failed to execute query: {:?}", e);
        e
    })?;
    Ok(())
}
