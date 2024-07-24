// Service user is an example implementation of an Encore service.
package users

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type User struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

//encore:api public method=GET path=/users
func List(ctx context.Context) (*ListResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT id, name
		FROM users
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &ListResponse{Users: users}, nil
}

type ListResponse struct {
	Users []User `json:"users"`
}

//encore:api auth method=GET path=/users/:id
func Get(ctx context.Context, id int) (*GetResponse, error) {
	eb := errs.B().Meta("id", id)

	var user User
	err := sqldb.QueryRow(ctx, `
		SELECT id, name
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name)

	if errors.Is(err, sqldb.ErrNoRows) {
		return nil, eb.Code(errs.InvalidArgument).Msg("no user found").Err()
	}

	if err != nil {
		return nil, err
	}

	return &GetResponse{User: user}, nil
}

type GetResponse struct {
	User User `json:"user"`
}

//encore:api auth method=POST path=/users
func Create(ctx context.Context, params CreateParams) (*CreateResponse, error) {
	eb := errs.B().Meta("params", params)

	if len(params.Name) == 0 {
		return nil, eb.Code(errs.InvalidArgument).Msg("name is empty").Err()
	}

	var user User
	err := sqldb.QueryRow(ctx, `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id, name
	`, params.Name).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, err
	}

	return &CreateResponse{User: user}, nil
}

type CreateParams struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	User User `json:"user"`
}
