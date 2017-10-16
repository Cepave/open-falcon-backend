ALTER TABLE `host`
 ADD resource_object_id int(10);

UPDATE `falcon_portal`.`host` SET resource_object_id = id;

ALTER TABLE `grp`
  ADD objects_search_term varchar(150),
  ADD object_group_type varchar(50),
  ADD auto_sync BOOLEAN;
