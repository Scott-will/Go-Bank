CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()',
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

ALTER TABLE "accounts"
ADD CONSTRAINT "accounts_owner_fkey"
FOREIGN KEY ("owner") REFERENCES "users" ("username");
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

