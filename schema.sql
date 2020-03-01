DROP DATABASE IF EXISTS `gonnect`;
CREATE DATABASE `gonnect`;

USE `gonnect`;

CREATE TABLE IF NOT EXISTS `player`
(
    `id`    INT                                    NOT NULL AUTO_INCREMENT,
    `alias` VARCHAR(64) UNIQUE                     NOT NULL,
    `state` ENUM ('away', 'searching', 'reserved') NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `match`
(
    `id`    INT                                 NOT NULL AUTO_INCREMENT,
    `state` ENUM ('creating', 'ready', 'ended') NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `match_players`
(
    `id`        INT NOT NULL AUTO_INCREMENT,
    `match_id`  int,
    `player_id` int,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`match_id`) REFERENCES `match` (`id`) ON DELETE CASCADE,
    FOREIGN KEY (`player_id`) REFERENCES `player` (`id`) ON DELETE CASCADE
);
