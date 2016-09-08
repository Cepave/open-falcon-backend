INSERT INTO nqm_agent(ag_id, ag_connection_id, ag_hostname, ag_ip_address)
VALUES(1311, 'conn-1', 'a@host1', '20.30.40.1'),
	(1312, 'conn-2', 'b@host2', '20.30.40.2'),
	(1313, 'conn-3', 'b@host3', '20.30.40.3');

INSERT INTO nqm_ping_task(pt_ag_id, pt_period)
VALUES(1311, 60), (1312, 60), (1313, 60);

INSERT INTO owl_name_tag(nt_id, nt_value)
VALUES(3001, 'tag-1'), (3002, 'tag-2');

INSERT INTO nqm_pt_target_filter_name_tag(tfnt_pt_ag_id, tfnt_nt_id)
VALUES(1311, 3001),
	(1312, 3002);

INSERT INTO nqm_pt_target_filter_isp(tfisp_pt_ag_id, tfisp_isp_id)
VALUES(1312, 2),
	(1313, 3);

INSERT INTO nqm_pt_target_filter_province(tfpv_pt_ag_id, tfpv_pv_id)
VALUES(1311, 11),
	(1312, 12);

INSERT INTO nqm_pt_target_filter_city(tfct_pt_ag_id, tfct_ct_id)
VALUES(1312, 21),
	(1313, 22);
