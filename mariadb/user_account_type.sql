CREATE TABLE `user_account_type`
(
  `id` SMALLINT UNSIGNED NOT NULL,
  `description` VARCHAR(100) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

INSERT INTO user_account_type (id, description)
VALUES
(100, 'A YAP admin account can create, read, update and delete all user accounts, groups and PBXs'),
(101, 'A YAP regular account can read all user accounts, groups and PBXs'),
(200, 'A group admin can read and update thier own PBX(s) and group'),
(201, 'A group regular account can read thier own PBX(s) and group'),
(300, 'A PBX admin account can read and update thier own PBX'),
(301, 'A PBX regular account can read thier own PBX');
