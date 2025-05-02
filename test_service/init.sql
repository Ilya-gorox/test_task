CREATE DATABASE users_db;

\c users_db;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    age INTEGER NOT NULL,
    created TIMESTAMP NOT NULL
); 