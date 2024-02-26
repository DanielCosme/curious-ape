use actix_web::{web, HttpResponse};
use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::domain::{HabitName, NewHabit};

#[derive(serde::Deserialize)]
pub struct FormData {
    name: String,
    description: String,
}

impl TryFrom<FormData> for NewHabit {
    type Error = String;
    fn try_from(value: FormData) -> Result<Self, Self::Error> {
        let name = HabitName::parse(value.name)?;
        Ok(Self {
            name,
            description: value.description,
        })
    }
}

#[tracing::instrument(
    name = "Creating a new Habit",
    skip(form, pg_pool),
    fields(habit_name = %form.name)
)]
pub async fn create_habit(form: web::Form<FormData>, pg_pool: web::Data<PgPool>) -> HttpResponse {
    let h = match form.0.try_into() {
        Ok(h) => h,
        Err(_) => return HttpResponse::BadRequest().finish(),
    };

    match insert_habit(&pg_pool, h).await {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

#[tracing::instrument(name = "Saving a new habit in the database", skip(new_habit, pool))]
pub async fn insert_habit(pool: &PgPool, new_habit: NewHabit) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"
            INSERT INTO habits (id, name, description, created_at)
            VALUES ($1, $2, $3, $4)
        "#,
        Uuid::new_v4(),
        new_habit.name.as_ref(),
        new_habit.description,
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
