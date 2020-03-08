DROP DATABASE IF EXISTS `gonnect`;
CREATE DATABASE `gonnect`;

USE `gonnect`;

CREATE TABLE IF NOT EXISTS `players`
(
    `id`         INT                               NOT NULL AUTO_INCREMENT,
    `alias`      VARCHAR(64) UNIQUE                NOT NULL,
    `state`      ENUM ('searching', 'unavailable') NOT NULL,
    `created_at` DATETIME                          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME                          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `matches`
(
    `id`         INT                     NOT NULL AUTO_INCREMENT,
    `state`      ENUM ('ready', 'ended') NOT NULL,
    `created_at` DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `matches_players`
(
    `id`         INT      NOT NULL AUTO_INCREMENT,
    `match_id`   int,
    `player_id`  int,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`match_id`) REFERENCES `matches` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`player_id`) REFERENCES `players` (`id`) ON DELETE CASCADE
);
