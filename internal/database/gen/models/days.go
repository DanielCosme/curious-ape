// Code generated by BobGen sqlite v0.28.1. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
)

// Day is an object representing the database table.
type Day struct {
	ID   int32     `db:"id,pk" `
	Date time.Time `db:"date" `

	R dayR `db:"-" `
}

// DaySlice is an alias for a slice of pointers to Day.
// This should almost always be used instead of []*Day.
type DaySlice []*Day

// Days contains methods to work with the days table
var Days = sqlite.NewTablex[*Day, DaySlice, *DaySetter]("", "days")

// DaysQuery is a query on the days table
type DaysQuery = *sqlite.ViewQuery[*Day, DaySlice]

// DaysStmt is a prepared statment on days
type DaysStmt = bob.QueryStmt[*Day, DaySlice]

// dayR is where relationships are stored.
type dayR struct {
	DeepWorkLogs DeepWorkLogSlice // fk_deep_work_logs_0
	FitnessLogs  FitnessLogSlice  // fk_fitness_logs_0
	Habits       HabitSlice       // fk_habits_0
	SleepLogs    SleepLogSlice    // fk_sleep_logs_0
}

// DaySetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type DaySetter struct {
	ID   omit.Val[int32]     `db:"id,pk" `
	Date omit.Val[time.Time] `db:"date" `
}

func (s DaySetter) SetColumns() []string {
	vals := make([]string, 0, 2)
	if !s.ID.IsUnset() {
		vals = append(vals, "id")
	}

	if !s.Date.IsUnset() {
		vals = append(vals, "date")
	}

	return vals
}

func (s DaySetter) Overwrite(t *Day) {
	if !s.ID.IsUnset() {
		t.ID, _ = s.ID.Get()
	}
	if !s.Date.IsUnset() {
		t.Date, _ = s.Date.Get()
	}
}

func (s DaySetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 0, 2)
	if !s.ID.IsUnset() {
		vals = append(vals, sqlite.Arg(s.ID))
	}

	if !s.Date.IsUnset() {
		vals = append(vals, sqlite.Arg(s.Date))
	}

	return im.Values(vals...)
}

func (s DaySetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s DaySetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 2)

	if !s.ID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "id")...),
			sqlite.Arg(s.ID),
		}})
	}

	if !s.Date.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "date")...),
			sqlite.Arg(s.Date),
		}})
	}

	return exprs
}

type dayColumnNames struct {
	ID   string
	Date string
}

var DayColumns = buildDayColumns("days")

type dayColumns struct {
	tableAlias string
	ID         sqlite.Expression
	Date       sqlite.Expression
}

func (c dayColumns) Alias() string {
	return c.tableAlias
}

func (dayColumns) AliasedAs(alias string) dayColumns {
	return buildDayColumns(alias)
}

func buildDayColumns(alias string) dayColumns {
	return dayColumns{
		tableAlias: alias,
		ID:         sqlite.Quote(alias, "id"),
		Date:       sqlite.Quote(alias, "date"),
	}
}

type dayWhere[Q sqlite.Filterable] struct {
	ID   sqlite.WhereMod[Q, int32]
	Date sqlite.WhereMod[Q, time.Time]
}

func (dayWhere[Q]) AliasedAs(alias string) dayWhere[Q] {
	return buildDayWhere[Q](buildDayColumns(alias))
}

func buildDayWhere[Q sqlite.Filterable](cols dayColumns) dayWhere[Q] {
	return dayWhere[Q]{
		ID:   sqlite.Where[Q, int32](cols.ID),
		Date: sqlite.Where[Q, time.Time](cols.Date),
	}
}

type dayJoins[Q dialect.Joinable] struct {
	typ          string
	DeepWorkLogs func(context.Context) modAs[Q, deepWorkLogColumns]
	FitnessLogs  func(context.Context) modAs[Q, fitnessLogColumns]
	Habits       func(context.Context) modAs[Q, habitColumns]
	SleepLogs    func(context.Context) modAs[Q, sleepLogColumns]
}

func (j dayJoins[Q]) aliasedAs(alias string) dayJoins[Q] {
	return buildDayJoins[Q](buildDayColumns(alias), j.typ)
}

func buildDayJoins[Q dialect.Joinable](cols dayColumns, typ string) dayJoins[Q] {
	return dayJoins[Q]{
		typ:          typ,
		DeepWorkLogs: daysJoinDeepWorkLogs[Q](cols, typ),
		FitnessLogs:  daysJoinFitnessLogs[Q](cols, typ),
		Habits:       daysJoinHabits[Q](cols, typ),
		SleepLogs:    daysJoinSleepLogs[Q](cols, typ),
	}
}

// FindDay retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindDay(ctx context.Context, exec bob.Executor, IDPK int32, cols ...string) (*Day, error) {
	if len(cols) == 0 {
		return Days.Query(
			ctx, exec,
			SelectWhere.Days.ID.EQ(IDPK),
		).One()
	}

	return Days.Query(
		ctx, exec,
		SelectWhere.Days.ID.EQ(IDPK),
		sm.Columns(Days.Columns().Only(cols...)),
	).One()
}

// DayExists checks the presence of a single record by primary key
func DayExists(ctx context.Context, exec bob.Executor, IDPK int32) (bool, error) {
	return Days.Query(
		ctx, exec,
		SelectWhere.Days.ID.EQ(IDPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the Day
func (o *Day) PrimaryKeyVals() bob.Expression {
	return sqlite.Arg(o.ID)
}

// Update uses an executor to update the Day
func (o *Day) Update(ctx context.Context, exec bob.Executor, s *DaySetter) error {
	return Days.Update(ctx, exec, s, o)
}

// Delete deletes a single Day record with an executor
func (o *Day) Delete(ctx context.Context, exec bob.Executor) error {
	return Days.Delete(ctx, exec, o)
}

// Reload refreshes the Day using the executor
func (o *Day) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Days.Query(
		ctx, exec,
		SelectWhere.Days.ID.EQ(o.ID),
	).One()
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

func (o DaySlice) UpdateAll(ctx context.Context, exec bob.Executor, vals DaySetter) error {
	return Days.Update(ctx, exec, &vals, o...)
}

func (o DaySlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Days.Delete(ctx, exec, o...)
}

func (o DaySlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	IDPK := make([]int32, len(o))

	for i, o := range o {
		IDPK[i] = o.ID
	}

	mods = append(mods,
		SelectWhere.Days.ID.In(IDPK...),
	)

	o2, err := Days.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.ID != old.ID {
				continue
			}
			new.R = old.R
			*old = *new
			break
		}
	}

	return nil
}

func daysJoinDeepWorkLogs[Q dialect.Joinable](from dayColumns, typ string) func(context.Context) modAs[Q, deepWorkLogColumns] {
	return func(ctx context.Context) modAs[Q, deepWorkLogColumns] {
		return modAs[Q, deepWorkLogColumns]{
			c: DeepWorkLogColumns,
			f: func(to deepWorkLogColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, DeepWorkLogs.Name(ctx).As(to.Alias())).On(
						to.DayID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

func daysJoinFitnessLogs[Q dialect.Joinable](from dayColumns, typ string) func(context.Context) modAs[Q, fitnessLogColumns] {
	return func(ctx context.Context) modAs[Q, fitnessLogColumns] {
		return modAs[Q, fitnessLogColumns]{
			c: FitnessLogColumns,
			f: func(to fitnessLogColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, FitnessLogs.Name(ctx).As(to.Alias())).On(
						to.DayID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

func daysJoinHabits[Q dialect.Joinable](from dayColumns, typ string) func(context.Context) modAs[Q, habitColumns] {
	return func(ctx context.Context) modAs[Q, habitColumns] {
		return modAs[Q, habitColumns]{
			c: HabitColumns,
			f: func(to habitColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, Habits.Name(ctx).As(to.Alias())).On(
						to.DayID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

func daysJoinSleepLogs[Q dialect.Joinable](from dayColumns, typ string) func(context.Context) modAs[Q, sleepLogColumns] {
	return func(ctx context.Context) modAs[Q, sleepLogColumns] {
		return modAs[Q, sleepLogColumns]{
			c: SleepLogColumns,
			f: func(to sleepLogColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, SleepLogs.Name(ctx).As(to.Alias())).On(
						to.DayID.EQ(from.ID),
					))
				}

				return mods
			},
		}
	}
}

// DeepWorkLogs starts a query for related objects on deep_work_logs
func (o *Day) DeepWorkLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) DeepWorkLogsQuery {
	return DeepWorkLogs.Query(ctx, exec, append(mods,
		sm.Where(DeepWorkLogColumns.DayID.EQ(sqlite.Arg(o.ID))),
	)...)
}

func (os DaySlice) DeepWorkLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) DeepWorkLogsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = sqlite.ArgGroup(o.ID)
	}

	return DeepWorkLogs.Query(ctx, exec, append(mods,
		sm.Where(sqlite.Group(DeepWorkLogColumns.DayID).In(PKArgs...)),
	)...)
}

// FitnessLogs starts a query for related objects on fitness_logs
func (o *Day) FitnessLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) FitnessLogsQuery {
	return FitnessLogs.Query(ctx, exec, append(mods,
		sm.Where(FitnessLogColumns.DayID.EQ(sqlite.Arg(o.ID))),
	)...)
}

func (os DaySlice) FitnessLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) FitnessLogsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = sqlite.ArgGroup(o.ID)
	}

	return FitnessLogs.Query(ctx, exec, append(mods,
		sm.Where(sqlite.Group(FitnessLogColumns.DayID).In(PKArgs...)),
	)...)
}

// Habits starts a query for related objects on habits
func (o *Day) Habits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) HabitsQuery {
	return Habits.Query(ctx, exec, append(mods,
		sm.Where(HabitColumns.DayID.EQ(sqlite.Arg(o.ID))),
	)...)
}

func (os DaySlice) Habits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) HabitsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = sqlite.ArgGroup(o.ID)
	}

	return Habits.Query(ctx, exec, append(mods,
		sm.Where(sqlite.Group(HabitColumns.DayID).In(PKArgs...)),
	)...)
}

// SleepLogs starts a query for related objects on sleep_logs
func (o *Day) SleepLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) SleepLogsQuery {
	return SleepLogs.Query(ctx, exec, append(mods,
		sm.Where(SleepLogColumns.DayID.EQ(sqlite.Arg(o.ID))),
	)...)
}

func (os DaySlice) SleepLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) SleepLogsQuery {
	PKArgs := make([]bob.Expression, len(os))
	for i, o := range os {
		PKArgs[i] = sqlite.ArgGroup(o.ID)
	}

	return SleepLogs.Query(ctx, exec, append(mods,
		sm.Where(sqlite.Group(SleepLogColumns.DayID).In(PKArgs...)),
	)...)
}

func (o *Day) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "DeepWorkLogs":
		rels, ok := retrieved.(DeepWorkLogSlice)
		if !ok {
			return fmt.Errorf("day cannot load %T as %q", retrieved, name)
		}

		o.R.DeepWorkLogs = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Day = o
			}
		}
		return nil
	case "FitnessLogs":
		rels, ok := retrieved.(FitnessLogSlice)
		if !ok {
			return fmt.Errorf("day cannot load %T as %q", retrieved, name)
		}

		o.R.FitnessLogs = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Day = o
			}
		}
		return nil
	case "Habits":
		rels, ok := retrieved.(HabitSlice)
		if !ok {
			return fmt.Errorf("day cannot load %T as %q", retrieved, name)
		}

		o.R.Habits = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Day = o
			}
		}
		return nil
	case "SleepLogs":
		rels, ok := retrieved.(SleepLogSlice)
		if !ok {
			return fmt.Errorf("day cannot load %T as %q", retrieved, name)
		}

		o.R.SleepLogs = rels

		for _, rel := range rels {
			if rel != nil {
				rel.R.Day = o
			}
		}
		return nil
	default:
		return fmt.Errorf("day has no relationship %q", name)
	}
}

func ThenLoadDayDeepWorkLogs(queryMods ...bob.Mod[*dialect.SelectQuery]) sqlite.Loader {
	return sqlite.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadDayDeepWorkLogs(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load DayDeepWorkLogs", retrieved)
		}

		err := loader.LoadDayDeepWorkLogs(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadDayDeepWorkLogs loads the day's DeepWorkLogs into the .R struct
func (o *Day) LoadDayDeepWorkLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.DeepWorkLogs = nil

	related, err := o.DeepWorkLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Day = o
	}

	o.R.DeepWorkLogs = related
	return nil
}

// LoadDayDeepWorkLogs loads the day's DeepWorkLogs into the .R struct
func (os DaySlice) LoadDayDeepWorkLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	deepWorkLogs, err := os.DeepWorkLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.DeepWorkLogs = nil
	}

	for _, o := range os {
		for _, rel := range deepWorkLogs {
			if o.ID != rel.DayID {
				continue
			}

			rel.R.Day = o

			o.R.DeepWorkLogs = append(o.R.DeepWorkLogs, rel)
		}
	}

	return nil
}

func ThenLoadDayFitnessLogs(queryMods ...bob.Mod[*dialect.SelectQuery]) sqlite.Loader {
	return sqlite.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadDayFitnessLogs(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load DayFitnessLogs", retrieved)
		}

		err := loader.LoadDayFitnessLogs(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadDayFitnessLogs loads the day's FitnessLogs into the .R struct
func (o *Day) LoadDayFitnessLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.FitnessLogs = nil

	related, err := o.FitnessLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Day = o
	}

	o.R.FitnessLogs = related
	return nil
}

// LoadDayFitnessLogs loads the day's FitnessLogs into the .R struct
func (os DaySlice) LoadDayFitnessLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	fitnessLogs, err := os.FitnessLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.FitnessLogs = nil
	}

	for _, o := range os {
		for _, rel := range fitnessLogs {
			if o.ID != rel.DayID {
				continue
			}

			rel.R.Day = o

			o.R.FitnessLogs = append(o.R.FitnessLogs, rel)
		}
	}

	return nil
}

func ThenLoadDayHabits(queryMods ...bob.Mod[*dialect.SelectQuery]) sqlite.Loader {
	return sqlite.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadDayHabits(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load DayHabits", retrieved)
		}

		err := loader.LoadDayHabits(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadDayHabits loads the day's Habits into the .R struct
func (o *Day) LoadDayHabits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.Habits = nil

	related, err := o.Habits(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Day = o
	}

	o.R.Habits = related
	return nil
}

// LoadDayHabits loads the day's Habits into the .R struct
func (os DaySlice) LoadDayHabits(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	habits, err := os.Habits(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.Habits = nil
	}

	for _, o := range os {
		for _, rel := range habits {
			if o.ID != rel.DayID {
				continue
			}

			rel.R.Day = o

			o.R.Habits = append(o.R.Habits, rel)
		}
	}

	return nil
}

func ThenLoadDaySleepLogs(queryMods ...bob.Mod[*dialect.SelectQuery]) sqlite.Loader {
	return sqlite.Loader(func(ctx context.Context, exec bob.Executor, retrieved any) error {
		loader, isLoader := retrieved.(interface {
			LoadDaySleepLogs(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
		})
		if !isLoader {
			return fmt.Errorf("object %T cannot load DaySleepLogs", retrieved)
		}

		err := loader.LoadDaySleepLogs(ctx, exec, queryMods...)

		// Don't cause an issue due to missing relationships
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	})
}

// LoadDaySleepLogs loads the day's SleepLogs into the .R struct
func (o *Day) LoadDaySleepLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.SleepLogs = nil

	related, err := o.SleepLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, rel := range related {
		rel.R.Day = o
	}

	o.R.SleepLogs = related
	return nil
}

// LoadDaySleepLogs loads the day's SleepLogs into the .R struct
func (os DaySlice) LoadDaySleepLogs(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	sleepLogs, err := os.SleepLogs(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, o := range os {
		o.R.SleepLogs = nil
	}

	for _, o := range os {
		for _, rel := range sleepLogs {
			if o.ID != rel.DayID {
				continue
			}

			rel.R.Day = o

			o.R.SleepLogs = append(o.R.SleepLogs, rel)
		}
	}

	return nil
}

func insertDayDeepWorkLogs0(ctx context.Context, exec bob.Executor, deepWorkLogs1 []*DeepWorkLogSetter, day0 *Day) (DeepWorkLogSlice, error) {
	for i := range deepWorkLogs1 {
		deepWorkLogs1[i].DayID = omit.From(day0.ID)
	}

	ret, err := DeepWorkLogs.InsertMany(ctx, exec, deepWorkLogs1...)
	if err != nil {
		return ret, fmt.Errorf("insertDayDeepWorkLogs0: %w", err)
	}

	return ret, nil
}

func attachDayDeepWorkLogs0(ctx context.Context, exec bob.Executor, count int, deepWorkLogs1 DeepWorkLogSlice, day0 *Day) (DeepWorkLogSlice, error) {
	setter := &DeepWorkLogSetter{
		DayID: omit.From(day0.ID),
	}

	err := DeepWorkLogs.Update(ctx, exec, setter, deepWorkLogs1...)
	if err != nil {
		return nil, fmt.Errorf("attachDayDeepWorkLogs0: %w", err)
	}

	return deepWorkLogs1, nil
}

func (day0 *Day) InsertDeepWorkLogs(ctx context.Context, exec bob.Executor, related ...*DeepWorkLogSetter) error {
	if len(related) == 0 {
		return nil
	}

	deepWorkLogs1, err := insertDayDeepWorkLogs0(ctx, exec, related, day0)
	if err != nil {
		return err
	}

	day0.R.DeepWorkLogs = append(day0.R.DeepWorkLogs, deepWorkLogs1...)

	for _, rel := range deepWorkLogs1 {
		rel.R.Day = day0
	}
	return nil
}

func (day0 *Day) AttachDeepWorkLogs(ctx context.Context, exec bob.Executor, related ...*DeepWorkLog) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	deepWorkLogs1 := DeepWorkLogSlice(related)

	_, err = attachDayDeepWorkLogs0(ctx, exec, len(related), deepWorkLogs1, day0)
	if err != nil {
		return err
	}

	day0.R.DeepWorkLogs = append(day0.R.DeepWorkLogs, deepWorkLogs1...)

	for _, rel := range related {
		rel.R.Day = day0
	}

	return nil
}

func insertDayFitnessLogs0(ctx context.Context, exec bob.Executor, fitnessLogs1 []*FitnessLogSetter, day0 *Day) (FitnessLogSlice, error) {
	for i := range fitnessLogs1 {
		fitnessLogs1[i].DayID = omit.From(day0.ID)
	}

	ret, err := FitnessLogs.InsertMany(ctx, exec, fitnessLogs1...)
	if err != nil {
		return ret, fmt.Errorf("insertDayFitnessLogs0: %w", err)
	}

	return ret, nil
}

func attachDayFitnessLogs0(ctx context.Context, exec bob.Executor, count int, fitnessLogs1 FitnessLogSlice, day0 *Day) (FitnessLogSlice, error) {
	setter := &FitnessLogSetter{
		DayID: omit.From(day0.ID),
	}

	err := FitnessLogs.Update(ctx, exec, setter, fitnessLogs1...)
	if err != nil {
		return nil, fmt.Errorf("attachDayFitnessLogs0: %w", err)
	}

	return fitnessLogs1, nil
}

func (day0 *Day) InsertFitnessLogs(ctx context.Context, exec bob.Executor, related ...*FitnessLogSetter) error {
	if len(related) == 0 {
		return nil
	}

	fitnessLogs1, err := insertDayFitnessLogs0(ctx, exec, related, day0)
	if err != nil {
		return err
	}

	day0.R.FitnessLogs = append(day0.R.FitnessLogs, fitnessLogs1...)

	for _, rel := range fitnessLogs1 {
		rel.R.Day = day0
	}
	return nil
}

func (day0 *Day) AttachFitnessLogs(ctx context.Context, exec bob.Executor, related ...*FitnessLog) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	fitnessLogs1 := FitnessLogSlice(related)

	_, err = attachDayFitnessLogs0(ctx, exec, len(related), fitnessLogs1, day0)
	if err != nil {
		return err
	}

	day0.R.FitnessLogs = append(day0.R.FitnessLogs, fitnessLogs1...)

	for _, rel := range related {
		rel.R.Day = day0
	}

	return nil
}

func insertDayHabits0(ctx context.Context, exec bob.Executor, habits1 []*HabitSetter, day0 *Day) (HabitSlice, error) {
	for i := range habits1 {
		habits1[i].DayID = omit.From(day0.ID)
	}

	ret, err := Habits.InsertMany(ctx, exec, habits1...)
	if err != nil {
		return ret, fmt.Errorf("insertDayHabits0: %w", err)
	}

	return ret, nil
}

func attachDayHabits0(ctx context.Context, exec bob.Executor, count int, habits1 HabitSlice, day0 *Day) (HabitSlice, error) {
	setter := &HabitSetter{
		DayID: omit.From(day0.ID),
	}

	err := Habits.Update(ctx, exec, setter, habits1...)
	if err != nil {
		return nil, fmt.Errorf("attachDayHabits0: %w", err)
	}

	return habits1, nil
}

func (day0 *Day) InsertHabits(ctx context.Context, exec bob.Executor, related ...*HabitSetter) error {
	if len(related) == 0 {
		return nil
	}

	habits1, err := insertDayHabits0(ctx, exec, related, day0)
	if err != nil {
		return err
	}

	day0.R.Habits = append(day0.R.Habits, habits1...)

	for _, rel := range habits1 {
		rel.R.Day = day0
	}
	return nil
}

func (day0 *Day) AttachHabits(ctx context.Context, exec bob.Executor, related ...*Habit) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	habits1 := HabitSlice(related)

	_, err = attachDayHabits0(ctx, exec, len(related), habits1, day0)
	if err != nil {
		return err
	}

	day0.R.Habits = append(day0.R.Habits, habits1...)

	for _, rel := range related {
		rel.R.Day = day0
	}

	return nil
}

func insertDaySleepLogs0(ctx context.Context, exec bob.Executor, sleepLogs1 []*SleepLogSetter, day0 *Day) (SleepLogSlice, error) {
	for i := range sleepLogs1 {
		sleepLogs1[i].DayID = omit.From(day0.ID)
	}

	ret, err := SleepLogs.InsertMany(ctx, exec, sleepLogs1...)
	if err != nil {
		return ret, fmt.Errorf("insertDaySleepLogs0: %w", err)
	}

	return ret, nil
}

func attachDaySleepLogs0(ctx context.Context, exec bob.Executor, count int, sleepLogs1 SleepLogSlice, day0 *Day) (SleepLogSlice, error) {
	setter := &SleepLogSetter{
		DayID: omit.From(day0.ID),
	}

	err := SleepLogs.Update(ctx, exec, setter, sleepLogs1...)
	if err != nil {
		return nil, fmt.Errorf("attachDaySleepLogs0: %w", err)
	}

	return sleepLogs1, nil
}

func (day0 *Day) InsertSleepLogs(ctx context.Context, exec bob.Executor, related ...*SleepLogSetter) error {
	if len(related) == 0 {
		return nil
	}

	sleepLogs1, err := insertDaySleepLogs0(ctx, exec, related, day0)
	if err != nil {
		return err
	}

	day0.R.SleepLogs = append(day0.R.SleepLogs, sleepLogs1...)

	for _, rel := range sleepLogs1 {
		rel.R.Day = day0
	}
	return nil
}

func (day0 *Day) AttachSleepLogs(ctx context.Context, exec bob.Executor, related ...*SleepLog) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	sleepLogs1 := SleepLogSlice(related)

	_, err = attachDaySleepLogs0(ctx, exec, len(related), sleepLogs1, day0)
	if err != nil {
		return err
	}

	day0.R.SleepLogs = append(day0.R.SleepLogs, sleepLogs1...)

	for _, rel := range related {
		rel.R.Day = day0
	}

	return nil
}
