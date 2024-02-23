CREATE TABLE `users` (
`uuid` varchar(255) PRIMARY KEY,
`email` varchar(255) NOT NULL,
`username` varchar(255) NOT NULL,
`password` varchar(255) NOT NULL,
`first_name` varchar(255) NOT NULL,
`last_name` varchar(255) NOT NULL,
`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updated_at` datetime
) ENGINE=InnoDB DEFAULT CHARSET=utf8;