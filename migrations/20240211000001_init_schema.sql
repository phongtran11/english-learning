-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar(255) NOT NULL,
  "password" varchar(255) NOT NULL,
  "first_name" varchar(100),
  "last_name" varchar(100),
  "phone_number" varchar(20),
  "birthdate" date,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz
);

CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");

CREATE TABLE "sessions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" text NOT NULL,
  "user_agent" text,
  "client_ip" varchar(45),
  "is_revoked" boolean DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_sessions_refresh_token" ON "sessions" ("refresh_token");
CREATE INDEX "idx_sessions_user_id" ON "sessions" ("user_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "sessions";
DROP TABLE "users";
-- +goose StatementEnd
