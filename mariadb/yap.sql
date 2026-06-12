-- Create YAP tables

CREATE TABLE `user_group`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `group_name` VARCHAR(100) NOT NULL,
  `date_added` DATETIME DEFAULT NOW() NOT NULL,
  `group_active` BOOLEAN NOT NULL,
  `pbx_limit` SMALLINT UNSIGNED DEFAULT 1 NOT NULL,
  `new_pbx_sip_extension_default_limit` SMALLINT UNSIGNED NOT NULL,
  `new_pbx_sip_trunk_default_limit` SMALLINT UNSIGNED NOT NULL,
  `new_pbx_phone_number_default_limit` SMALLINT UNSIGNED NOT NULL,
  `new_pbx_cdr_default_limit` SMALLINT UNSIGNED NOT NULL,
  `new_pbx_voicemail_default_megabyte_limit` INT UNSIGNED NOT NULL,
  `new_pbx_call_recording_default_megabyte_limit` INT UNSIGNED NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `group_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75) NOT NULL,
  `city_town_village` VARCHAR(75) NOT NULL,
  `county_state_region` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(16) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `group_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75) NOT NULL,
  `city_town_village` VARCHAR(75) NOT NULL,
  `county_state_region` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(16) NOT NULL,
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
  `pbx_sip_extension_limit` SMALLINT UNSIGNED NOT NULL,
  `pbx_sip_trunk_limit` SMALLINT UNSIGNED NOT NULL,
  `pbx_phone_number_limit` SMALLINT UNSIGNED NOT NULL,
  `pbx_cdr_limit` SMALLINT UNSIGNED NOT NULL,
  `pbx_voicemail_megabyte_limit` INT UNSIGNED NOT NULL,
  `pbx_call_recording_megabyte_limit` INT UNSIGNED NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_invoice_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75) NOT NULL,
  `city_town_village` VARCHAR(75) NOT NULL,
  `county_state_region` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(16) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `pbx_site_address`
(
  `id` BIGINT UNSIGNED NOT NULL,
  `address_line_1` VARCHAR(75) NOT NULL,
  `address_line_2` VARCHAR(75) NOT NULL,
  `city_town_village` VARCHAR(75) NOT NULL,
  `county_state_region` VARCHAR(75) NOT NULL,
  `postcode_zip_code` VARCHAR(7) NOT NULL,
  `country` VARCHAR(75) NOT NULL,
  `contact_email` VARCHAR(255) NOT NULL,
  `contact_number` VARCHAR(16) NOT NULL,
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
UNIQUE (`email`),
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

CREATE TABLE `user_account_type`
(
  `id` SMALLINT UNSIGNED NOT NULL,
  `type` VARCHAR(20) NOT NULL,
  `permission` VARCHAR(4400) NOT NULL,
PRIMARY KEY(`id`)
)
ENGINE = InnoDB;

----------------------------------------------------------------------------------------------------

-- Add pbx_id and and endpoint_type column to Asterisk tables

ALTER TABLE `ps_endpoints`
ADD COLUMN `endpoint_type` ENUM ('sip_extension','sip_trunk','webrtc_extension') NOT NULL;
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

ALTER TABLE `pbx`
ADD CONSTRAINT fk___pbx___user_group
FOREIGN KEY (`group_id`)
REFERENCES `user_group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `group_invoice_address`
ADD CONSTRAINT fk___group_invoice_address___user_group
FOREIGN KEY (`id`)
REFERENCES `user_group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `group_site_address`
ADD CONSTRAINT fk___group_site_address___user_group
FOREIGN KEY (`id`)
REFERENCES `user_group` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_invoice_address`
ADD CONSTRAINT fk___pbx_invoice_address___pbx
FOREIGN KEY (`id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `pbx_site_address`
ADD CONSTRAINT fk___pbx_site_address___pbx
FOREIGN KEY (`id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___user_account_type
FOREIGN KEY (`user_account_type_id`)
REFERENCES `user_account_type` (`id`);

ALTER TABLE `user_account`
ADD CONSTRAINT fk___user_account___user_group
FOREIGN KEY (`group_id`)
REFERENCES `user_group` (`id`);

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

-- Create View(s)

CREATE VIEW `view___account_detail` AS
SELECT
  `user_account`.`user_account_type_id`,
  `user_account`.`first_name` AS 'user_account_first_name',
  `user_account`.`last_name` AS 'user_account_last_name',
  `user_account`.`email` AS 'user_account_email',
  `user_account_type`.`type` AS 'user_account_type',
  `user_account`.`date_added` AS 'user_account_date_added',
  `user_group`.`group_name`,
  `user_account`.`group_id`,
  `pbx`.`pbx_name`,
  `user_account`.`pbx_id`,
  `user_account_type`.`permission` AS 'user_account_type_permission',
  `group_site_address`.`address_line_1` AS 'group_site_address_line_1',
  `group_site_address`.`address_line_2` AS 'group_site_address_line_2',
  `group_site_address`.`city_town_village` AS 'group_site_city_town_village',
  `group_site_address`.`county_state_region` AS 'group_site_county_state_region',
  `group_site_address`.`postcode_zip_code` AS 'group_site_postcode_zip_code',
  `group_site_address`.`country` AS 'group_site_country',
  `group_site_address`.`contact_email` AS 'group_site_contact_email',
  `group_site_address`.`contact_number` AS 'group_site_contact_number',
  `group_invoice_address`.`address_line_1` AS 'group_invoice_address_line_1',
  `group_invoice_address`.`address_line_2` AS 'group_invoice_address_line_2',
  `group_invoice_address`.`city_town_village` AS 'group_invoice_city_town_village',
  `group_invoice_address`.`county_state_region` AS 'group_invoice_county_state_region',
  `group_invoice_address`.`postcode_zip_code` AS 'group_invoice_postcode_zip_code',
  `group_invoice_address`.`country` AS 'group_invoice_country',
  `group_invoice_address`.`contact_email` AS 'group_invoice_contact_email',
  `group_invoice_address`.`contact_number` AS 'group_invoice_contact_number',
  `pbx_site_address`.`address_line_1` AS 'pbx_site_address_line_1',
  `pbx_site_address`.`address_line_2` AS 'pbx_site_address_line_2',
  `pbx_site_address`.`city_town_village` AS 'pbx_site_city_town_village',
  `pbx_site_address`.`county_state_region` AS 'pbx_site_county_state_region',
  `pbx_site_address`.`postcode_zip_code` AS 'pbx_site_postcode_zip_code',
  `pbx_site_address`.`country` AS 'pbx_site_country',
  `pbx_site_address`.`contact_email` AS 'pbx_site_contact_email',
  `pbx_site_address`.`contact_number` AS 'pbx_site_contact_number',
  `pbx_invoice_address`.`address_line_1` AS 'pbx_invoice_address_line_1',
  `pbx_invoice_address`.`address_line_2` AS 'pbx_invoice_address_line_2',
  `pbx_invoice_address`.`city_town_village` AS 'pbx_invoice_city_town_village',
  `pbx_invoice_address`.`county_state_region` AS 'pbx_invoice_county_state_region',
  `pbx_invoice_address`.`postcode_zip_code` AS 'pbx_invoice_postcode_zip_code',
  `pbx_invoice_address`.`country` AS 'pbx_invoice_country',
  `pbx_invoice_address`.`contact_email` AS 'pbx_invoice_contact_email',
  `pbx_invoice_address`.`contact_number` AS 'pbx_invoice_contact_number'
FROM `user_account`
INNER JOIN `user_account_type`
ON `user_account`.`user_account_type_id` = `user_account_type`.`id`
INNER JOIN `user_group`
ON `user_account`.`group_id` = `user_group`.`id`
INNER JOIN `pbx`
ON `user_account`.`pbx_id` = `pbx`.`id`
INNER JOIN `group_site_address`
ON `user_account`.`group_id` = `group_site_address`.`id`
INNER JOIN `group_invoice_address`
ON `user_account`.`group_id` = `group_invoice_address`.`id`
INNER JOIN `pbx_site_address`
ON `user_account`.`pbx_id` = `pbx_site_address`.`id`
INNER JOIN `pbx_invoice_address`
ON `user_account`.`pbx_id` = `pbx_invoice_address`.`id`;

CREATE VIEW `view___group_detail` AS
SELECT
  `user_group`.`id` AS 'group_id',
  `user_group`.`group_name`,
  `user_group`.`date_added` AS 'group_date_added',
  `user_group`.`group_active`,
  `user_group`.`pbx_limit`,
  `user_group`.`new_pbx_sip_extension_default_limit`,
  `user_group`.`new_pbx_sip_trunk_default_limit`,
  `user_group`.`new_pbx_phone_number_default_limit`,
  `user_group`.`new_pbx_cdr_default_limit`,
  `user_group`.`new_pbx_voicemail_default_megabyte_limit`,
  `user_group`.`new_pbx_call_recording_default_megabyte_limit`,
  `group_site_address`.`address_line_1` AS 'group_site_address_line_1',
  `group_site_address`.`address_line_2` AS 'group_site_address_line_2',
  `group_site_address`.`city_town_village` AS 'group_site_city_town_village',
  `group_site_address`.`county_state_region` AS 'group_site_county_state_region',
  `group_site_address`.`postcode_zip_code` AS 'group_site_postcode_zip_code',
  `group_site_address`.`country` AS 'group_site_country',
  `group_site_address`.`contact_email` AS 'group_site_contact_email',
  `group_site_address`.`contact_number` AS 'group_site_contact_number',
  `group_invoice_address`.`address_line_1` AS 'group_invoice_address_line_1',
  `group_invoice_address`.`address_line_2` AS 'group_invoice_address_line_2',
  `group_invoice_address`.`city_town_village` AS 'group_invoice_city_town_village',
  `group_invoice_address`.`county_state_region` AS 'group_invoice_county_state_region',
  `group_invoice_address`.`postcode_zip_code` AS 'group_invoice_postcode_zip_code',
  `group_invoice_address`.`country` AS 'group_invoice_country',
  `group_invoice_address`.`contact_email` AS 'group_invoice_contact_email',
  `group_invoice_address`.`contact_number` AS 'group_invoice_contact_number'
FROM `user_group`
INNER JOIN `group_site_address`
ON `user_group`.`id` = `group_site_address`.`id`
INNER JOIN `group_invoice_address`
ON `user_group`.`id` = `group_invoice_address`.`id`;

CREATE VIEW `view___pbx_detail` AS
SELECT
  `pbx`.`id` AS 'pbx_id',
  `pbx`.`pbx_name`,
  `pbx`.`group_id`,
  `user_group`.`group_name`,
  `pbx`.`date_added` AS 'pbx_date_added',
  `pbx`.`pbx_active`,
  `pbx`.`pbx_sip_extension_limit`,
  `pbx`.`pbx_sip_trunk_limit`,
  `pbx`.`pbx_phone_number_limit`,
  `pbx`.`pbx_cdr_limit`,
  `pbx`.`pbx_voicemail_megabyte_limit`,
  `pbx`.`pbx_call_recording_megabyte_limit`,
  `pbx_site_address`.`address_line_1` AS 'pbx_site_address_line_1',
  `pbx_site_address`.`address_line_2` AS 'pbx_site_address_line_2',
  `pbx_site_address`.`city_town_village` AS 'pbx_site_city_town_village',
  `pbx_site_address`.`county_state_region` AS 'pbx_site_county_state_region',
  `pbx_site_address`.`postcode_zip_code` AS 'pbx_site_postcode_zip_code',
  `pbx_site_address`.`country` AS 'pbx_site_country',
  `pbx_site_address`.`contact_email` AS 'pbx_site_contact_email',
  `pbx_site_address`.`contact_number` AS 'pbx_site_contact_number',
  `pbx_invoice_address`.`address_line_1` AS 'pbx_invoice_address_line_1',
  `pbx_invoice_address`.`address_line_2` AS 'pbx_invoice_address_line_2',
  `pbx_invoice_address`.`city_town_village` AS 'pbx_invoice_city_town_village',
  `pbx_invoice_address`.`county_state_region` AS 'pbx_invoice_county_state_region',
  `pbx_invoice_address`.`postcode_zip_code` AS 'pbx_invoice_postcode_zip_code',
  `pbx_invoice_address`.`country` AS 'pbx_invoice_country',
  `pbx_invoice_address`.`contact_email` AS 'pbx_invoice_contact_email',
  `pbx_invoice_address`.`contact_number` AS 'pbx_invoice_contact_number'
FROM `pbx`
INNER JOIN `user_group`
ON `pbx`.`group_id` = `user_group`.`id`
INNER JOIN `pbx_site_address`
ON `pbx`.`id` = `pbx_site_address`.`id`
INNER JOIN `pbx_invoice_address`
ON `pbx`.`id` = `pbx_invoice_address`.`id`;

CREATE VIEW `view___sip_extension_detail` AS
SELECT DISTINCT
  `ps_auths`.`username` AS 'sip_username',
  `ps_auths`.`password` AS 'sip_password',
  IFNULL(`ps_endpoints`.`callerid`, '(NOT SET)') AS 'caller_id',
  IFNULL(`ps_endpoints`.`callerid_privacy`, 'allowed_not_screened (DEFAULT)') AS 'caller_id_privacy',
  IFNULL(`ps_endpoints`.`named_call_group`, '(NOT SET)') AS 'call_group',
  `ps_endpoints`.`allow` AS 'codec_allowed',
  IFNULL(`ps_endpoints`.`direct_media`, 'yes (DEFAULT)') AS 'direct_media',
  IFNULL(`ps_endpoints`.`direct_media_method`, 'invite (DEFAULT)') AS 'direct_media_method',
  IFNULL(`ps_endpoints`.`dtmf_mode`, 'rfc4733 (DEFAULT)') AS 'dtmf_mode',
  IFNULL(`ps_endpoints`.`force_rport`, 'yes (DEFAULT)') AS 'force_rport',  
  IFNULL(`ps_endpoints`.`from_user`, '(NOT SET)') AS 'from_sip_header_user',
  IFNULL(`ps_endpoints`.`from_domain`, '(NOT SET)') AS 'from_sip_header_domain',
  IFNULL(`ps_endpoints`.`permit`, '(NOT SET)') AS 'ip_address_allowed',
  IFNULL(`ps_endpoints`.`named_pickup_group`, '(NOT SET)') AS 'pickup_group',
  IFNULL(`ps_endpoints`.`media_encryption`, 'no (RECOMMENDED TO ENABLE TLS OR SETUP A VPN SERVER)') AS 'media_encryption_enabled',
  IFNULL(`ps_endpoints`.`stir_shaken`, 'no (DEFAULT)') AS 'stir_shaken_enabled',
  IFNULL(`ps_endpoints`.`stir_shaken_profile`, '(NOT SET)') AS 'stir_shaken_profile',
  `ps_contacts`.`endpoint` IS NOT NULL AS 'registered',
  `pbx`.`pbx_name`,
  `pbx`.`id` AS 'pbx_id',
  `user_group`.`group_name`,
  `user_group`.`id` AS 'group_id'
FROM `ps_endpoints`
INNER JOIN `ps_auths`
ON `ps_endpoints`.`id` = `ps_auths`.`id`
INNER JOIN `pbx`
ON `ps_endpoints`.`pbx_id` = `pbx`.`id`
LEFT JOIN `ps_contacts`
on `ps_endpoints`.`id` = `ps_contacts`.`endpoint`
INNER JOIN `user_group`
ON `pbx`.`group_id` = `user_group`.`id`
WHERE `ps_endpoints`.`endpoint_type` = 'sip_extension';

CREATE VIEW `view___sip_extension_registered` AS
SELECT
  `ps_auths`.`username` AS 'sip_username',
  `ps_contacts`.`uri`,
  `ps_contacts`.`user_agent`,
  `pbx`.`pbx_name` AS 'pbx_name',
  `pbx`.`id` AS 'pbx_id',
  `user_group`.`group_name` AS 'group_name',
  `user_group`.`id` AS 'group_id'
FROM `ps_endpoints`
INNER JOIN `ps_auths`
ON `ps_endpoints`.`id` = `ps_auths`.`id`
INNER JOIN `pbx`
ON `ps_endpoints`.`pbx_id` = `pbx`.`id`
INNER JOIN `ps_contacts`
on `ps_endpoints`.`id` = `ps_contacts`.`endpoint`
INNER JOIN `user_group`
ON `pbx`.`group_id` = `user_group`.`id`
WHERE `ps_endpoints`.`endpoint_type` = 'sip_extension';

----------------------------------------------------------------------------------------------------

-- Insert data to YAP tables

INSERT INTO `user_group` (`id`, `group_name`, `group_active`, `pbx_limit`, `new_pbx_sip_endpoint_default_limit`, `new_pbx_sip_trunk_default_limit`, `new_pbx_phone_number_default_limit`, `new_pbx_cdr_default_limit`, `new_pbx_voicemail_default_megabyte_limit`, `new_pbx_call_recording_default_megabyte_limit`)
VALUES (1, 'system', 1, 0, 0, 0, 0, 0, 0, 0);
  
INSERT INTO `group_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `group_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx` (`id`, `pbx_name`, `group_id`, `pbx_active`, `pbx_sip_endpoint_limit`, `pbx_sip_trunk_limit`, `pbx_phone_number_limit`, `pbx_cdr_limit`, `pbx_voicemail_megabyte_limit`, `pbx_call_recording_megabyte_limit`)
VALUES (1, 'system', 1, 0, 0, 0, 0, 0, 0, 0);
  
INSERT INTO `pbx_invoice_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `pbx_site_address` (`id`,	`address_line_1`,	`address_line_2`,	`city_town_village`, `postcode_zip_code`,	`county_state_region`, `country`,	`contact_email`, `contact_number`)
VALUES (1, 'system', 'system', 'system', 'system', 'system', 'system', 'system', 'system');

INSERT INTO `user_account_type` (`id`, `type`, `permission`)
VALUES
(100, 'YAP Admin (100)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#9989 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a YAP Admin (100) User Account<br>
      &#9989 View a YAP Admin (100) User Account<br>
      &#9989 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#9989 Create a Group Admin (200) User Account<br>
      &#9989 View a Group Admin (200) User Account<br>
      &#9989 Update a Group Admin (200) User Account<br>
      &#9989 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a Group Regular (201) User Account<br>
      &#9989 View a Group Regular (201) User Account<br>
      &#9989 Update a Group Regular (201) User Account<br>
      &#9989 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#9989 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#9989 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#9989 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#9989 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#9989 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Group<br>
      &#9940 Update Own Group<br>
      &#9940 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a Group<br>
      &#9989 View a Group<br>
      &#9989 Update a Group<br>
      &#9989 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#9989 Create a Group Invoice<br>
      &#9989 View a Group Invoice<br>
      &#9989 Update a Group Invoice<br>
      &#9989 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View YAP User Account Logs<br>
      &#9989 View Group Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#9989 View Server Information<br>
      &#9989 Download Server Information<br>
      &#9989 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(200, 'Group Admin (200)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Group Admin (200) User Account<br>
      &#9989 View a Group Admin (200) User Account<br>
      &#10060 Update a Group Admin (200) User Account<br>
      &#10060 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Group Regular (201) User Account<br>
      &#9989 View a Group Regular (201) User Account<br>
      &#9989 Update a Group Regular (201) User Account<br>
      &#10060 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View Own Group<br>
      &#9989 Update Own Group<br>
      &#10060 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Group<br>
      &#10060 View a Group<br>
      &#10060 Update a Group<br>
      &#10060 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Group Invoice<br>
      &#9989 View a Group Invoice<br>
      &#10060 Update a Group Invoice<br>
      &#10060 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#9989 View Group Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(201, 'Group Regular (201)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Group Admin (200) User Account<br>
      &#10060 View a Group Admin (200) User Account<br>
      &#10060 Update a Group Admin (200) User Account<br>
      &#10060 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Group Regular (201) User Account<br>
      &#10060 View a Group Regular (201) User Account<br>
      &#10060 Update a Group Regular (201) User Account<br>
      &#10060 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#9989 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 View Own Group<br>
      &#10060 Update Own Group<br>
      &#10060 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Group<br>
      &#10060 View a Group<br>
      &#10060 Update a Group<br>
      &#10060 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9940 View Own PBX<br>
      &#9940 Update Own PBX<br>
      &#9940 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#9989 Create a PBX<br>
      &#9989 View a PBX<br>
      &#9989 Update a PBX<br>
      &#9989 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Group Invoice<br>
      &#9989 View a Group Invoice<br>
      &#10060 Update a Group Invoice<br>
      &#10060 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Group Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(300, 'PBX Admin (300)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#9989 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Group Admin (200) User Account<br>
      &#10060 View a Group Admin (200) User Account<br>
      &#10060 Update a Group Admin (200) User Account<br>
      &#10060 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Group Regular (201) User Account<br>
      &#10060 View a Group Regular (201) User Account<br>
      &#10060 Update a Group Regular (201) User Account<br>
      &#10060 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#9989 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#9989 View a PBX Regular (301) User Account<br>
      &#9989 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Group<br>
      &#9940 Update Own Group<br>
      &#9940 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Group<br>
      &#10060 View a Group<br>
      &#10060 Update a Group<br>
      &#10060 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#9989 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Group Invoice<br>
      &#10060 View a Group Invoice<br>
      &#10060 Update a Group Invoice<br>
      &#10060 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Group Logs<br>
      &#9989 View PBX Logs<br>
      &#9989 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(301, 'PBX Regular (301)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Group Admin (200) User Account<br>
      &#10060 View a Group Admin (200) User Account<br>
      &#10060 Update a Group Admin (200) User Account<br>
      &#10060 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Group Regular (201) User Account<br>
      &#10060 View a Group Regular (201) User Account<br>
      &#10060 Update a Group Regular (201) User Account<br>
      &#10060 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#10060 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#10060 View a PBX Regular (301) User Account<br>
      &#10060 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#9989 View a PBX Read Only (302) User Account<br>
      &#9989 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Group<br>
      &#9940 Update Own Group<br>
      &#9940 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Group<br>
      &#10060 View a Group<br>
      &#10060 Update a Group<br>
      &#10060 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#10060 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9989 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#9989 Update a SIP Endpoint<br>
      &#9989 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Group Invoice<br>
      &#10060 View a Group Invoice<br>
      &#10060 Update a Group Invoice<br>
      &#10060 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Group Logs<br>
      &#10060 View PBX Logs<br>
      &#10060 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
),
(302, 'PBX Read Only (302)',
'<table>
  <tr>
    <td>
      <b>Key:</b><br>
      &#9989 = Allowed<br>
      &#10060 = Prohibited<br>
      &#9940 = Not Applicable<br>
    </td>
    <td>
      &#9989 View Own User Account<br>
      &#10060 Update Own User Account<br>
      &#10060 Delete Own User Account<br>
      <div class="main-menu-space"></div>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a YAP Admin (100) User Account<br>
      &#10060 View a YAP Admin (100) User Account<br>
      &#10060 Update a YAP Admin (100) User Account<br>
      &#10060 Delete a YAP Admin (100) User Account<br>
    </td>
    <td>
      &#10060 Create a Group Admin (200) User Account<br>
      &#10060 View a Group Admin (200) User Account<br>
      &#10060 Update a Group Admin (200) User Account<br>
      &#10060 Delete a Group Admin (200) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a Group Regular (201) User Account<br>
      &#10060 View a Group Regular (201) User Account<br>
      &#10060 Update a Group Regular (201) User Account<br>
      &#10060 Delete a Group Regular (201) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Admin (300) User Account<br>
      &#10060 View a PBX Admin (300) User Account<br>
      &#10060 Update a PBX Admin (300) User Account<br>
      &#10060 Delete a PBX Admin (300) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a PBX Regular (301) User Account<br>
      &#10060 View a PBX Regular (301) User Account<br>
      &#10060 Update a PBX Regular (301) User Account<br>
      &#10060 Delete a PBX Regular (301) User Account<br>
    </td>
    <td>
      &#10060 Create a PBX Read Only (302) User Account<br>
      &#10060 View a PBX Read Only (302) User Account<br>
      &#10060 Update a PBX Read Only (302) User Account<br>
      &#10060 Delete a PBX Read Only (302) User Account<br>
    </td>
  </tr>
  <tr>
    <td>
      &#9940 View Own Group<br>
      &#9940 Update Own Group<br>
      &#9940 Delete Own Group<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a Group<br>
      &#10060 View a Group<br>
      &#10060 Update a Group<br>
      &#10060 Delete a Group<br>
    </td>
  <tr>
    <td>
      &#9989 View Own PBX<br>
      &#10060 Update Own PBX<br>
      &#10060 Delete Own PBX<br>
      <div class="main-menu-space"></div>
    </td>
    <td>
      &#10060 Create a PBX<br>
      &#10060 View a PBX<br>
      &#10060 Update a PBX<br>
      &#10060 Delete a PBX<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 Create a SIP Endpoint<br>
      &#9989 View a SIP Endpoint<br>
      &#10060 Update a SIP Endpoint<br>
      &#10060 Delete a SIP Endpoint<br>
    </td>
    <td>
      &#10060 Create a Group Invoice<br>
      &#10060 View a Group Invoice<br>
      &#10060 Update a Group Invoice<br>
      &#10060 Delete a Group Invoice<br>
    </td>
  </tr>
  <tr>
    <td>
      &#10060 View YAP User Account Logs<br>
      &#10060 View Group Logs<br>
      &#10060 View PBX Logs<br>
      &#10060 Download Logs<br>
    </td>
    <td>
      &#10060 View Server Information<br>
      &#10060 Download Server Information<br>
      &#10060 Set Resource Limits<br>
      &#9989 View Resource Limits<br>
    </td>
  </tr>
</table>'
);
