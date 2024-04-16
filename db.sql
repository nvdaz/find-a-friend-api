CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `name` TEXT NOT NULL,
    `updated_at` DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS `user_profiles` (
    `id` VARCHAR(36) PRIMARY KEY,
    `bio` TEXT NOT NULL,
    `extroversion` FLOAT NOT NULL,
    `agreeableness` FLOAT NOT NULL,
    `conscientiousness` FLOAT NOT NULL,
    `neuroticism` FLOAT NOT NULL,
    `openness` FLOAT NOT NULL,
    `updated_at` DATETIME NOT NULL,
    CONSTRAINT `fk_id` FOREIGN KEY (`id`) REFERENCES `users`(`id`)
);

CREATE TABLE IF NOT EXISTS `interests` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user_id` VARCHAR(36),
    `interest` TEXT,
    `intensity` FLOAT NOT NULL,
    `skill` FLOAT NOT NULL,
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
    CONSTRAINT `unique_interest` UNIQUE (`user_id`, `interest`)
);

CREATE TABLE IF NOT EXISTS `service_conversations` (
    `id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36),
    `question` TEXT NOT NULL,
    `answer` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);
