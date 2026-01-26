-- Create YAP tables

CREATE TABLE `group`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `group_name` VARCHAR(100) NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `group_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `group_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75),
  `city_town_village` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7),
  `county_state_region` VARCHAR(75),
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(16),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `group_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75),
  `city_town_village` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7),
  `county_state_region` VARCHAR(75),
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(16),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `pbx_name` VARCHAR(75) NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `pbx_active` BOOLEAN NOT NULL,
  `note` VARCHAR(255),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75),
  `city_town_village` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7),
  `county_state_region` VARCHAR(75),
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(16),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75),
  `city_town_village` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7),
  `county_state_region` VARCHAR(75),
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255),
  `contact_number` VARCHAR(16),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `user_account_type_id` SMALLINT UNSIGNED NOT NULL,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `pbx_id` BIGINT UNSIGNED NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `account_active` BOOLEAN NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account_type`
(
  `id` SMALLINT UNSIGNED NOT NULL,
  `description` VARCHAR(100) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

----------------------------------------------------------------------------------------------------

-- Add pbx_id column to Asterisk tables

ALTER TABLE `ps_endpoints`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_aors`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_auths`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

----------------------------------------------------------------------------------------------------

-- Alter existing data types

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `aors` varchar(255) NOT NULL;

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `auth` varchar(255) NOT NULL;

ALTER TABLE `ps_aors`
MODIFY COLUMN `id` varchar(255) NOT NULL;

----------------------------------------------------------------------------------------------------

-- Add index to columns

ALTER TABLE `pbx`
ADD INDEX `index___pbx__group_id` (`group_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__user_account_type_id` (`user_account_type_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__group_id` (`group_id`);

ALTER TABLE `user_account`
ADD INDEX `index___user_account__pbx_id` (`pbx_id`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__aors` (`aors`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__auth` (`auth`);

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__pbx_id` (`pbx_id`);

ALTER TABLE `ps_aors`
ADD INDEX `index___ps_aors__pbx_id` (`pbx_id`);

ALTER TABLE `ps_auths`
ADD INDEX `index___ps_auths__pbx_id` (`pbx_id`);

----------------------------------------------------------------------------------------------------

-- Create foreign key constraints

ALTER TABLE `group_invoice_address`
ADD CONSTRAINT fk___group_invoice_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `group_site_address`
ADD CONSTRAINT fk___group_site_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_invoice_address`
ADD CONSTRAINT fk___pbx_invoice_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_site_address`
ADD CONSTRAINT fk___pbx_site_address___group
FOREIGN KEY (`id`)
REFERENCES `group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___user_account_type
FOREIGN KEY (`user_account_type_id`)
REFERENCES `user_account_type` (`id`);

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___group
FOREIGN KEY (`group_id`)
REFERENCES `group` (`id`);

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`);

ALTER TABLE `ps_endpoints`
ADD CONSTRAINT fk___ps_endpoints___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___ps_endpoints
FOREIGN KEY (`id`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___ps_endpoints
FOREIGN KEY (`id`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

----------------------------------------------------------------------------------------------------

-- Insert data to YAP tables

INSERT INTO `group` (`id`, `group_name`, `date_added`, `group_active`, `note`) VALUES (1, 'system', NOW(), 1, 'created during YAP install');

INSERT INTO `user_account_type` (`id`, `description`)
VALUES
(100, 'A YAP admin account can create, read, update and delete all user accounts, groups and PBXs'),
(101, 'A YAP regular account can read all user accounts, groups and PBXs'),
(200, 'A group admin can read and update thier own PBX(s) and group'),
(201, 'A group regular account can read thier own PBX(s) and group'),
(300, 'A PBX admin account can read and update thier own PBX'),
(301, 'A PBX regular account can read thier own PBX');
