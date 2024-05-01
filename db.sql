CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `name` TEXT NOT NULL,
    `updated_at` DATETIME NOT NULL
    `profile` TEXT NOT NULL
    `generated_at` DATETIME NOT NULL
);


CREATE TABLE IF NOT EXISTS `service_conversations` (
    `id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36),
    `question` TEXT NOT NULL,
    `answer` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);

CREATE TABLE IF NOT EXISTS `matches` (
    `id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36),
    `match_id` VARCHAR(36),
    `reason` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_match_id` FOREIGN KEY (`match_id`) REFERENCES `users`(`id`)
    CONSTRAINT `unique_match` UNIQUE (`user_id`, `match_id`)
);

CREATE TABLE IF NOT EXISTS `messages` (
    `id` VARCHAR(36) PRIMARY KEY,
    `sender_id` VARCHAR(36),
    `receiver_id` VARCHAR(36),
    `message` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    CONSTRAINT `fk_sender_id` FOREIGN KEY (`sender_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_receiver_id` FOREIGN KEY (`receiver_id`) REFERENCES `users`(`id`)
);
