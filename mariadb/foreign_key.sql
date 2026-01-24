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
FOREIGN KEY (`id`)
REFERENCES `user_account_type` (`id`)
