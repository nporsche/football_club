CREATE TABLE `players` (
    `id` int NOT NULL,
    `name` varchar(64) NOT NULL,
    `cellphone` varchar(20) NOT NULL,
    `status` tinyint NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `match_log` (
    `id` int NOT NULL,
    `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `competitor` varchar(64) NOT NULL,
    `goal` tinyint NOT NULL,
    `loss` tinyint NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `duration_log` (
    `match_id` int NOT NULL,
    `player_id` int NOT NULL,
    `duration` smallint NOT NULL DEFAULT 0,
    `status` tinyint NOT NULL DEFAULT 0,
    PRIMARY KEY (`match_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `goal_log` (
    `match_id` int NOT NULL,
    `player_id` int NOT NULL,
    `goal` tiny NOT NULL DEFAULT 0,
    `goal_type` tinyint NOT NULL DEFAULT 0,
    PRIMARY KEY (`match_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `revenue_log` (
    `id` int NOT NULL,
    `player_id` int NOT NULL,
    `datetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `amount` int NOT NULL,
    `reason` varchar(128) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8






