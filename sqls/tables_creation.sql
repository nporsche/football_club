CREATE TABLE if not exists `players` (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL,
    `cellphone` varchar(20) NOT NULL,
    `status` tinyint NOT NULL,
    `tag` smallint NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE if not exists `match_log` (
    `id` int NOT NULL AUTO_INCREMENT,
    `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `competitor` varchar(64) NOT NULL,
    `cost` int NOT NULL,
    `goal` tinyint NOT NULL,
    `loss` tinyint NOT NULL,
    `author` varchar(64),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE if not exists `duration_log`(
    `match_id` int NOT NULL,
    `player_id` int NOT NULL,
    `duration` smallint NOT NULL DEFAULT 0,
    `status` tinyint NOT NULL DEFAULT 0,
    `author` varchar(64),
    PRIMARY KEY (`match_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE if not exists `goal_log` (
    `match_id` int NOT NULL,
    `player_id` int NOT NULL,
    `goal_type` varchar(32) NOT NULL
    `author` varchar(64),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE if not exists `revenue_log` (
    `id` int NOT NULL AUTO_INCREMENT,
    `player_id` int NOT NULL,
    `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `amount` int NOT NULL,
    `reason` varchar(128) NOT NULL,
    `author` varchar(64),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;






