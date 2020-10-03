SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

USE clothes;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` int(11) NOT NULL,
    `state` varchar(200) NOT NULL DEFAULT 'main',
    `recent` int(11) NOT NULL DEFAULT 0,
    `topcolor` varchar(200),
    `bottomcolor` varchar(200),
    `season` varchar(200),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `data`;
CREATE TABLE `data` (
     `user` int(11) NOT NULL,
     `id` int(11) NOT NULL DEFAULT 0,
     `photo` varchar(200),
     `name` varchar(200),
     `type` varchar(200),
     `cold` boolean DEFAULT false,
     `normal` boolean DEFAULT false,
     `warm` boolean DEFAULT false,
     `hot` boolean DEFAULT false,
     `color` varchar(200),
     `purity` varchar(200) DEFAULT 'clean',
     PRIMARY KEY (`user`, `id`),
     INDEX (`type`),
     INDEX (`cold`),
     INDEX (`normal`),
     INDEX (`warm`),
     INDEX (`hot`),
     INDEX (`color`),
     INDEX (`purity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;