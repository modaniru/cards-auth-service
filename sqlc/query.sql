-- name: GetUserByAuthTypeAndAuthId :one
SELECT user_id FROM users_auths WHERE auth_type = $1 AND auth_id = $2;