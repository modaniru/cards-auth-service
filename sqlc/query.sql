-- name: GetUserByAuthTypeAndAuthId :one
SELECT user_id FROM users_auths WHERE auth_type = $1 AND auth_id = $2;

-- name: AddUserAuthType :exec
INSERT INTO users_auths (user_id, auth_type, auth_id) values ($1, $2, $3);

-- name: CreateEmptyUser :one
INSERT INTO users DEFAULT VALUES RETURNING id;