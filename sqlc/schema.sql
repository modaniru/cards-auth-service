CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    email varchar(255) unique DEFAULT NULL,
    password varchar(255),
    username varchar(255)  
);

CREATE TABLE IF NOT EXISTS users_auths (
    user_id integer REFERENCES users(id),
    auth_type varchar(255),
    auth_id varchar(255),
    unique (user_id, auth_type)
);