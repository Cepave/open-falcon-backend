INSERT INTO host(hostname, ip, agent_version, plugin_version)
VALUES('no-nqm-1', '1.2.30.1', '', ''),
	('no-nqm-2', '1.2.30.2', '', ''),
	('nqm-1', '12.77.121.70', '', ''),
	('nqm-2', '12.77.121.71', '', '');

INSERT INTO nqm_agent(ag_name, ag_connection_id, ag_hostname, ag_ip_address)
VALUES('nqm-1', 'nqm-1@ip-01', 'nqm-1', x'5B172D0B'),
	('nqm-2', 'nqm-1@ip-02', 'nqm-2', x'5B172D0C'),
	('nqm-3', 'nqm-1@ip-03', 'nqm-3', x'5B172D0D'),
	('nqm-4', 'nqm-1@ip-04', 'nqm-4', x'5B172D0E');
