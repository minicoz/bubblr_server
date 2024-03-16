
CREATE TABLE IF NOT EXISTS schools (
    id SERIAL PRIMARY KEY,
    school TEXT NOT NULL,
    tier INT NOT NULL
);

CREATE INDEX ON "schools" ("id");

CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(40),
    last_name VARCHAR(40),
    email VARCHAR(80),
    hashed_password VARCHAR(265),
    school_id INT,
    dob_day SMALLINT,
    dob_month SMALLINT,
    dob_year SMALLINT,
    is_male BOOLEAN NOT NULL DEFAULT(false),
    grad_year INT,
    verified BOOLEAN,
    about TEXT,
    created_at timestamptz NOT NULL DEFAULT(now())
);

ALTER TABLE "users" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");
CREATE INDEX ON "users" ("id");

CREATE TABLE IF NOT EXISTS pictures (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    url TEXT
);

ALTER TABLE "pictures" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
CREATE INDEX ON "pictures" ("id");


CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    matched_user_id UUID NOT NULL
);

ALTER TABLE "matches" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
CREATE INDEX ON "matches" ("id");

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT(now()),
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    txt_message TEXT
);
