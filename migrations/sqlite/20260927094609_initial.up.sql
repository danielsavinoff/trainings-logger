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
-- Create "sets" table
CREATE TABLE `sets` (
  `id` text NOT NULL,
  `exercise_id` text NOT NULL,
  `exercise_weight_in_grams` integer NOT NULL DEFAULT 0,
  `repetitions` integer NOT NULL DEFAULT 0,
  `rest_in_seconds` integer NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`),
  CONSTRAINT `0` FOREIGN KEY (`exercise_id`) REFERENCES `exercises` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CHECK (exercise_weight_in_grams >= 0),
  CHECK (repetitions BETWEEN 0 AND 255),
  CHECK (rest_in_seconds >= 0)
);
-- Create index "sets_exercise_created_at" to table: "sets"
CREATE INDEX `sets_exercise_created_at` ON `sets` (`exercise_id`, `created_at`);
