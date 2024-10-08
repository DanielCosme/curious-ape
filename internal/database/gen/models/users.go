// Code generated by BobGen sqlite v0.28.1. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"

	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dialect"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
	"github.com/stephenafamo/bob/expr"
)

// User is an object representing the database table.
type User struct {
	ID       int32  `db:"id,pk" `
	Username string `db:"username" `
	Password string `db:"password" `
	Role     string `db:"role" `
	Email    string `db:"email" `
}

// UserSlice is an alias for a slice of pointers to User.
// This should almost always be used instead of []*User.
type UserSlice []*User

// Users contains methods to work with the users table
var Users = sqlite.NewTablex[*User, UserSlice, *UserSetter]("", "users")

// UsersQuery is a query on the users table
type UsersQuery = *sqlite.ViewQuery[*User, UserSlice]

// UsersStmt is a prepared statment on users
type UsersStmt = bob.QueryStmt[*User, UserSlice]

// UserSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type UserSetter struct {
	ID       omit.Val[int32]  `db:"id,pk" `
	Username omit.Val[string] `db:"username" `
	Password omit.Val[string] `db:"password" `
	Role     omit.Val[string] `db:"role" `
	Email    omit.Val[string] `db:"email" `
}

func (s UserSetter) SetColumns() []string {
	vals := make([]string, 0, 5)
	if !s.ID.IsUnset() {
		vals = append(vals, "id")
	}

	if !s.Username.IsUnset() {
		vals = append(vals, "username")
	}

	if !s.Password.IsUnset() {
		vals = append(vals, "password")
	}

	if !s.Role.IsUnset() {
		vals = append(vals, "role")
	}

	if !s.Email.IsUnset() {
		vals = append(vals, "email")
	}

	return vals
}

func (s UserSetter) Overwrite(t *User) {
	if !s.ID.IsUnset() {
		t.ID, _ = s.ID.Get()
	}
	if !s.Username.IsUnset() {
		t.Username, _ = s.Username.Get()
	}
	if !s.Password.IsUnset() {
		t.Password, _ = s.Password.Get()
	}
	if !s.Role.IsUnset() {
		t.Role, _ = s.Role.Get()
	}
	if !s.Email.IsUnset() {
		t.Email, _ = s.Email.Get()
	}
}

func (s UserSetter) InsertMod() bob.Mod[*dialect.InsertQuery] {
	vals := make([]bob.Expression, 0, 5)
	if !s.ID.IsUnset() {
		vals = append(vals, sqlite.Arg(s.ID))
	}

	if !s.Username.IsUnset() {
		vals = append(vals, sqlite.Arg(s.Username))
	}

	if !s.Password.IsUnset() {
		vals = append(vals, sqlite.Arg(s.Password))
	}

	if !s.Role.IsUnset() {
		vals = append(vals, sqlite.Arg(s.Role))
	}

	if !s.Email.IsUnset() {
		vals = append(vals, sqlite.Arg(s.Email))
	}

	return im.Values(vals...)
}

func (s UserSetter) Apply(q *dialect.UpdateQuery) {
	um.Set(s.Expressions()...).Apply(q)
}

func (s UserSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 5)

	if !s.ID.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "id")...),
			sqlite.Arg(s.ID),
		}})
	}

	if !s.Username.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "username")...),
			sqlite.Arg(s.Username),
		}})
	}

	if !s.Password.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "password")...),
			sqlite.Arg(s.Password),
		}})
	}

	if !s.Role.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "role")...),
			sqlite.Arg(s.Role),
		}})
	}

	if !s.Email.IsUnset() {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			sqlite.Quote(append(prefix, "email")...),
			sqlite.Arg(s.Email),
		}})
	}

	return exprs
}

type userColumnNames struct {
	ID       string
	Username string
	Password string
	Role     string
	Email    string
}

var UserColumns = buildUserColumns("users")

type userColumns struct {
	tableAlias string
	ID         sqlite.Expression
	Username   sqlite.Expression
	Password   sqlite.Expression
	Role       sqlite.Expression
	Email      sqlite.Expression
}

func (c userColumns) Alias() string {
	return c.tableAlias
}

func (userColumns) AliasedAs(alias string) userColumns {
	return buildUserColumns(alias)
}

func buildUserColumns(alias string) userColumns {
	return userColumns{
		tableAlias: alias,
		ID:         sqlite.Quote(alias, "id"),
		Username:   sqlite.Quote(alias, "username"),
		Password:   sqlite.Quote(alias, "password"),
		Role:       sqlite.Quote(alias, "role"),
		Email:      sqlite.Quote(alias, "email"),
	}
}

type userWhere[Q sqlite.Filterable] struct {
	ID       sqlite.WhereMod[Q, int32]
	Username sqlite.WhereMod[Q, string]
	Password sqlite.WhereMod[Q, string]
	Role     sqlite.WhereMod[Q, string]
	Email    sqlite.WhereMod[Q, string]
}

func (userWhere[Q]) AliasedAs(alias string) userWhere[Q] {
	return buildUserWhere[Q](buildUserColumns(alias))
}

func buildUserWhere[Q sqlite.Filterable](cols userColumns) userWhere[Q] {
	return userWhere[Q]{
		ID:       sqlite.Where[Q, int32](cols.ID),
		Username: sqlite.Where[Q, string](cols.Username),
		Password: sqlite.Where[Q, string](cols.Password),
		Role:     sqlite.Where[Q, string](cols.Role),
		Email:    sqlite.Where[Q, string](cols.Email),
	}
}

// FindUser retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindUser(ctx context.Context, exec bob.Executor, IDPK int32, cols ...string) (*User, error) {
	if len(cols) == 0 {
		return Users.Query(
			ctx, exec,
			SelectWhere.Users.ID.EQ(IDPK),
		).One()
	}

	return Users.Query(
		ctx, exec,
		SelectWhere.Users.ID.EQ(IDPK),
		sm.Columns(Users.Columns().Only(cols...)),
	).One()
}

// UserExists checks the presence of a single record by primary key
func UserExists(ctx context.Context, exec bob.Executor, IDPK int32) (bool, error) {
	return Users.Query(
		ctx, exec,
		SelectWhere.Users.ID.EQ(IDPK),
	).Exists()
}

// PrimaryKeyVals returns the primary key values of the User
func (o *User) PrimaryKeyVals() bob.Expression {
	return sqlite.Arg(o.ID)
}

// Update uses an executor to update the User
func (o *User) Update(ctx context.Context, exec bob.Executor, s *UserSetter) error {
	return Users.Update(ctx, exec, s, o)
}

// Delete deletes a single User record with an executor
func (o *User) Delete(ctx context.Context, exec bob.Executor) error {
	return Users.Delete(ctx, exec, o)
}

// Reload refreshes the User using the executor
func (o *User) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := Users.Query(
		ctx, exec,
		SelectWhere.Users.ID.EQ(o.ID),
	).One()
	if err != nil {
		return err
	}

	*o = *o2

	return nil
}

func (o UserSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals UserSetter) error {
	return Users.Update(ctx, exec, &vals, o...)
}

func (o UserSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	return Users.Delete(ctx, exec, o...)
}

func (o UserSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	var mods []bob.Mod[*dialect.SelectQuery]

	IDPK := make([]int32, len(o))

	for i, o := range o {
		IDPK[i] = o.ID
	}

	mods = append(mods,
		SelectWhere.Users.ID.In(IDPK...),
	)

	o2, err := Users.Query(ctx, exec, mods...).All()
	if err != nil {
		return err
	}

	for _, old := range o {
		for _, new := range o2 {
			if new.ID != old.ID {
				continue
			}

			*old = *new
			break
		}
	}

	return nil
}
