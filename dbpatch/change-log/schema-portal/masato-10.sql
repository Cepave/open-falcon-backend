SET NAMES 'utf8';

ALTER TABLE falcon_portal.events ADD COLUMN status int(3) unsigned DEFAULT 0;
