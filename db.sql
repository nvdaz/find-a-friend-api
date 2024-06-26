CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `username` TEXT UNIQUE NOT NULL,
    `password` BLOB NOT NULL,
    `name` TEXT NOT NULL,
    `avatar` TEXT,
    `updated_at` DATETIME NOT NULL,
    `profile` JSONB,
    `generated_at` DATETIME
);

CREATE TABLE IF NOT EXISTS `matches` (
    `id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36),
    `other_id` VARCHAR(36),
    `reason` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_other_id` FOREIGN KEY (`other_id`) REFERENCES `users`(`id`)
    CONSTRAINT `unique_match` UNIQUE (`user_id`, `other_id`)
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
CREATE INDEX idx_sender_id ON messages (sender_id);
CREATE INDEX idx_receiver_id ON messages (receiver_id);
CREATE INDEX idx_created_at ON messages (created_at);
CREATE INDEX idx_sender_receiver_created_at ON messages (sender_id, receiver_id, created_at);
