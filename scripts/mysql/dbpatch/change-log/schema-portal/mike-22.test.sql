INSERT INTO host(id, hostname)
VALUES(788091, 'host-1'),(3356781, 'host-2'),(4351, 'host-3'),(2876541, 'host-4');

INSERT INTO grp(id, grp_name)
VALUES(41, 'group-1'),(42, 'group-2'),(43, 'group-3'),(44, 'group-4');

INSERT INTO grp_host(grp_id, host_id)
VALUES (41, 788091), (41, 3356781),
	(42, 3356781), (42, 4351),
	(43, 4351), (43, 2876541),
	(44, 3356781), (43, 788091),
	(42, 55670), (143, 3356781), (141, 4351); -- 要被刪除，無法聯結的資料表
