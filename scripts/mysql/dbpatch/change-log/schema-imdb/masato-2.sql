-- a template table made for building `host_tags_records` `tag_search_lists`
CREATE OR REPLACE VIEW host_tags AS
  SELECT resource_object_id, tag_id, value_id, hostname, host_id, name as tag_name, object_tag_id, tags.tag_type_id  from
      (SELECT ot.resource_object_id as resource_object_id, ot.tag_id as tag_id, ot.id as object_tag_id, ot.value_id as value_id, ho.hostname as hostname, ho.host_id as host_id from
        (SELECT h.id as host_id, h.hostname as hostname, h.resource_object_id as id FROM `falcon_portal`.host as h JOIN `resource_objects` as ro ON h.resource_object_id = ro.id) as ho
        JOIN `object_tags` as ot ON ho.id = ot.resource_object_id) as ht2 JOIN `tags` ON ht2.tag_id = tags.id;

-- build tag values select cache, only support str_values & value_models
CREATE OR REPLACE VIEW tag_search_lists AS
  SELECT st1.tag_id as tag_id, st1.value as value from
    (SELECT ht.tag_id as tag_id, sv.value as value from
      (SELECT value_id, tag_id, tag_name from host_tags where tag_type_id = 1) as ht
      JOIN str_values as sv on ht.value_id = sv.id) as st1 GROUP BY st1.value, st1.tag_id
  UNION
  SELECT tag_id, value from value_models;

-- build a quick search table for finding tag value matching by host
CREATE OR REPLACE VIEW host_tags_records AS
  SELECT * from
    (
      -- select string type
      SELECT ht.object_tag_id, hostname, host_id, sv.resource_object_id, sv.tag_id as tag_id, tag_name, value, value_id from
        (SELECT * from host_tags where tag_type_id = 1) as ht
        JOIN str_values as sv on ht.value_id = sv.id

      UNION

      -- select value_model type
      SELECT v2.object_tag_id, v2.hostname as hostname, v2.host_id as host_id, v2.resource_object_id as resource_object_id, v2.tag_id as tag_id, v2.tag_name as tag_name, vmd.value as value, v2.value_id as value_id from
        (SELECT ht.object_tag_id, hostname, host_id, sv.resource_object_id, sv.tag_id, tag_name, value_model_id, value_id from
          (SELECT * from host_tags where tag_type_id = 3) as ht
            JOIN vmodel_values as sv on ht.value_id = sv.id) as v2
        JOIN value_models as vmd ON v2.value_model_id = vmd.id

      UNION

      -- select integer type
      SELECT ht.object_tag_id, ht.hostname, ht.host_id, sv.resource_object_id, sv.tag_id, ht.tag_name, CONVERT(sv.value, char) as value, ht.value_id as value_id from
        (SELECT * from host_tags where tag_type_id = 2) as ht
        JOIN int_values as sv on ht.value_id = sv.id

      UNION

      -- select description type
      SELECT ht.object_tag_id, hostname, host_id, sv.resource_object_id, sv.tag_id as tag_id, tag_name, value, value_id from
        (SELECT * from host_tags where tag_type_id = 4) as ht
        JOIN description_values as sv on ht.value_id = sv.id
    ) as htr;
