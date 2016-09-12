INSERT INTO owl_group_tag(gt_id, gt_name)
VALUES(43, 'group-1'), (44, 'group-2'), (45, 'group-3'), (46, 'group-4');

INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
VALUES(10091, 'conn-01', 'host-01', 0x01020301),
	(326791, 'conn-02', 'host-02', 0x01020302),
	(998071, 'conn-03', 'host-03', 0x01020303),
	(44, 'conn-04', 'host-04', 0x01020304);

INSERT INTO nqm_ping_task(pt_id, pt_period)
VALUES(90, 60),(91, 60),(92, 60),(93, 60);

INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
VALUES(10091, 90), (326791, 91), (998071, 92), (44, 93);

INSERT INTO nqm_agent_group_tag(agt_ag_id, agt_gt_id)
VALUES(10091, 43), (10091, 44),
	(326791, 44),
	(998071, 45), (998071, 46),
	(44, 45), (44, 46);
