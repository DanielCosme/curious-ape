package data

//
// import (
// 	"database/sql"
// 	"errors"
//
// 	"github.com/danielcosme/curious-ape/internal/validator"
// )
//
// type FoodHabit struct {
// 	Habit
// }
//
// type FoodHabitModel struct {
// 	DB *sql.DB
// }
//
// func (fh *FoodHabitModel) GetAll() ([]*FoodHabit, error) {
// 	// food, work, sleep, fitness
// 	query := `
// 		SELECT id, state, date FROM food_habits`
// 	rows, err := fh.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	habits := []*FoodHabit{}
// 	for rows.Next() {
// 		var habit FoodHabit
//
// 		err := rows.Scan(
// 			&habit.ID,
// 			&habit.State,
// 			&habit.Date,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		habits = append(habits, &habit)
// 	}
//
// 	return habits, nil
// }
//
// func (fh *FoodHabitModel) Insert(habit *FoodHabit) error {
// 	query := `
// 		INSERT INTO food_habits (state, "date")
// 		VALUES ($1, $2)
// 		RETURNING id `
// 	args := []interface{}{habit.State, habit.Date}
// 	return fh.DB.QueryRow(query, args...).Scan(&habit.ID)
// }
//
// func (fh *FoodHabitModel) Get(date string) (*FoodHabit, error) {
// 	query := `
// 		SELECT id, state, date FROM food_habits
// 		WHERE "date" = $1`
// 	var habit FoodHabit
// 	err := fh.DB.QueryRow(query, date).Scan(
// 		&habit.ID,
// 		&habit.State,
// 		&habit.Date,
// 	)
//
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
//
// 	return &habit, nil
// }
//
// func (fh *FoodHabitModel) Update(habit *FoodHabit) error {
// 	stm := `
// 		UPDATE food_habits
// 		SET state = $1, "date" = $2
// 		WHERE "date" = $2
// 		RETURNING state, "date"`
//
// 	args := []interface{}{
// 		habit.State,
// 		habit.Date,
// 	}
//
// 	return fh.DB.QueryRow(stm, args...).Scan(&habit.State, &habit.Date)
// }
//
// func (fh *FoodHabitModel) Delete(id int) error {
// 	query := `
// 		DELETE FROM food_habits
// 		WHERE id = $1`
// 	result, err := fh.DB.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}
//
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
//
// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}
//
// 	return nil
// }
//
// func ValidateFoodHabit(v *validator.Validator, habit *FoodHabit) {
// 	v.Check(habit.Date != "", "date", "must be provided")
// 	v.Check(len([]rune(habit.Date)) == 10, "date", "must be exactly 10 characters long")
// }
