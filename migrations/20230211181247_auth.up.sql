CREATE TABLE IF NOT EXISTS "credentials"
(
    "id"             uuid         not null,
    "email"          varchar(320) not null,
    "password"       varchar(255),
    "verified_email" boolean      not null default false,
    "iat"            bigint       not null default 0,
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "google"
(
    "user_id"        uuid         not null,
    "email"          varchar(320) not null,
    "google_id"      varchar(255) not null,
    PRIMARY KEY ("user_id"),
    CONSTRAINT fk_credentials
        FOREIGN KEY ("user_id")
            REFERENCES credentials ("id")
);