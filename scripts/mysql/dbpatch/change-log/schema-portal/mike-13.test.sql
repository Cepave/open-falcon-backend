/**
 * 提供 mike-13.sql Patch 前的測試資料
 */

INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_name_tag)
VALUES(2001, 'tg-1', '1.2.3.4', null),
  (2002, 'tg-2', '1.2.3.5', '機房一'),
  (2003, 'tg-3', '1.2.3.6', '機房二'),
  (2004, 'tg-4', '1.2.3.7', '機房三'),
  (2005, 'tg-5', '1.2.3.8', '機房三');

INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
VALUES(1501, 'conn-1', 'a@host1', '20.30.40.1'),
	(1502, 'conn-2', 'b@host2', '20.30.40.2'),
	(1503, 'conn-3', 'b@host3', '20.30.40.3');

INSERT INTO nqm_ping_task(pt_ag_id, pt_period)
VALUES(1501, 60), (1502, 60), (1503, 60);

INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_ag_id, tfnt_name_tag)
VALUES(1501, '機房一'),
	(1502, '機房一條件一'),
	(1503, '機房一條件一');
