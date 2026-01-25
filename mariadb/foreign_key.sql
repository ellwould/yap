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
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_endpoints`
ADD INDEX `index___ps_endpoints__pbx_id` (`pbx_id`);

ALTER TABLE `ps_aors`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_aors`
ADD INDEX `index___ps_aors__pbx_id` (`pbx_id`);

ALTER TABLE `ps_auths`
ADD COLUMN `pbx_id` BIGINT UNSIGNED NOT NULL;

ALTER TABLE `ps_auths`
ADD INDEX `index___ps_auths__pbx_id` (`pbx_id`);

ALTER TABLE `ps_endpoints`
ADD CONSTRAINT fk___ps_endpoints___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___pbx
FOREIGN KEY (`pbx_id`)
REFERENCES `pbx` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `aors` varchar(255) NOT NULL;

ALTER TABLE `ps_endpoints`
MODIFY COLUMN `auth` varchar(255) NOT NULL;

ALTER TABLE ps_aors
MODIFY COLUMN `id` varchar(255) NOT NULL;

ALTER TABLE `ps_aors`
ADD INDEX `index___ps_endpoints__aors` (`aors`);

ALTER TABLE `ps_auths`
MODIFY COLUMN `id` varchar(255) NOT NULL;

ALTER TABLE `ps_auths`
ADD INDEX `index___ps_endpoints__auth` (`auth`);

ALTER TABLE `ps_aors`
ADD CONSTRAINT fk___ps_aors___ps_endpoints
FOREIGN KEY (`aors`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;

ALTER TABLE `ps_auths`
ADD CONSTRAINT fk___ps_auths___ps_endpoints
FOREIGN KEY (`auth`)
REFERENCES `ps_endpoints` (`id`)
ON DELETE CASCADE;
