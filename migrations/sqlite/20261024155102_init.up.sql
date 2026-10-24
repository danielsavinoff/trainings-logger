-- Create "preferences" table
CREATE TABLE `preferences` (
  `id` integer NULL,
  `weight_unit` text NOT NULL DEFAULT 'kg',
  PRIMARY KEY (`id`),
  CHECK (id = 1),
  CHECK (weight_unit IN ('g', 'kg', 'lb'))
);
-- Create "exercises" table
CREATE TABLE `exercises` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `body_weight_share` real NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  CHECK (length(trim(name)) > 0),
  CHECK (body_weight_share BETWEEN 0 AND 1)
);
-- Create index "exercises_name" to table: "exercises"
CREATE UNIQUE INDEX `exercises_name` ON `exercises` (`name`);

CREATE TABLE identities (
    id TEXT PRIMARY KEY NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('Dropbox'))
);

CREATE TABLE identity_credentials (
    identity_id TEXT PRIMARY KEY NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    access_token_expires_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE sessions (
 id text NOT NULL PRIMARY KEY,
 "index" integer NOT NULL UNIQUE CHECK("index" > 0),
 created_at datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE performances (
 id text NOT NULL PRIMARY KEY,
 session_id text NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
 exercise_id text NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
 created_at datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
 UNIQUE(session_id, exercise_id)
);
CREATE INDEX performances_session_created_at ON performances(session_id, created_at ASC, id ASC);
CREATE INDEX performances_exercise_id ON performances(exercise_id);

CREATE TABLE sets (
 id text NOT NULL PRIMARY KEY,
 performance_id text NOT NULL REFERENCES performances(id) ON DELETE CASCADE,
 exercise_weight_in_grams integer NOT NULL DEFAULT 0 CHECK(exercise_weight_in_grams >= 0),
 repetitions integer NOT NULL DEFAULT 0 CHECK(repetitions BETWEEN 0 AND 255),
 rest_in_seconds integer NOT NULL DEFAULT 0 CHECK(rest_in_seconds >= 0),
 created_at datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);
CREATE INDEX sets_performance_created_at ON sets(performance_id, created_at DESC, id DESC);
