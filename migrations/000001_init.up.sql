CREATE TABLE IF NOT EXISTS "credentials" (
    "id" uuid not null,
    "email" varchar(320) not null,
    "password" varchar(64),
    "verified_email" boolean not null default false,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "google" (
    "id" uuid not null,
    "email" varchar(320) not null,
    "google_id" varchar(255) not null,
    PRIMARY KEY ("id"),
    CONSTRAINT fk_credentials
        FOREIGN KEY("id")
            REFERENCES credentials("id")
);