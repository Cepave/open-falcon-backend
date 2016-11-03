SET NAMES 'utf8';

INSERT INTO owl_city(ct_id, ct_pv_id, ct_name, ct_post_code)
VALUES
	(293, 22, '毕节市', '551700'),
	(294, -1, '国外其它', '')
ON DUPLICATE KEY UPDATE
  ct_pv_id = VALUES(ct_pv_id),
  ct_name = VALUES(ct_name),
  ct_post_code = VALUES(ct_post_code);
