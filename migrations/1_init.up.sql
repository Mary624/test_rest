CREATE TABLE IF NOT EXISTS people
(
    id        serial PRIMARY KEY,
    name     TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INTEGER NOT NULL,
    gender TEXT NOT NULL,
    nationality TEXT NOT NULL
);