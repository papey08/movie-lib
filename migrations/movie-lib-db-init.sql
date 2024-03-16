CREATE TABLE "movies" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(150),
    "description" VARCHAR(1000),
    "release_date" DATE,
    "rating" FLOAT
);

CREATE TABLE "actors" (
    "id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR(100),
    "second_name" VARCHAR(100),
    "gender" VARCHAR(10)
);

CREATE TABLE "movie-actor" (
    "movie-id" INTEGER,
    "actor_id" INTEGER,
    UNIQUE ("movie-id", "actor_id")
);

CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "role" VARCHAR(10)
);

INSERT INTO "users" ("role")
VALUES
    ('admin'),
    ('regular')
;

SET TIME ZONE 'UTC';
