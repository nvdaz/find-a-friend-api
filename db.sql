CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `name` TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS `personalities` (
    `id` VARCHAR(36) PRIMARY KEY,
    `extraversion` FLOAT NOT NULL,
    `agreeableness` FLOAT NOT NULL,
    `conscientiousness` FLOAT NOT NULL,
    `neuroticism` FLOAT NOT NULL,
    `openness` FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS `interests` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` VARCHAR(36),
    `interest` TEXT,
    CONSTRAINT `unique_interest` UNIQUE (`user_id`, `interest`)
)
