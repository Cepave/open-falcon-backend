package test

var SetAutoIncForHost = `ALTER TABLE host AUTO_INCREMENT = 36091`
var ResetAutoIncForHost = `ALTER TABLE host AUTO_INCREMENT = 1`
var InsertHostSQL = `
INSERT INTO host(id, hostname, agent_version, plugin_version)
VALUES
  (36091, 'ct-agent-1', '', ''),
	(36092, 'ct-agent-2', '', ''),
	(36093, 'ct-agent-3', '', ''),
	(36094, 'ct-agent-4', '', '')
`
var DeleteHostSQL = `DELETE FROM host WHERE id >= 36091 AND id <= 36095`

var SetAutoIncForNqmAgent = `ALTER TABLE nqm_agent AUTO_INCREMENT = 24021`
var ResetAutoIncForNqmAgent = `ALTER TABLE nqm_agent AUTO_INCREMENT = 1`
var InsertNqmAgentSQL = `
INSERT INTO nqm_agent(
	ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address,
	ag_pv_id, ag_ct_id
)
VALUES
  (24021, 36091, 'ct-255-1', 'ct-255-1@201.3.116.1', 'ct-1', x'C9037401', 7, 255),
	(24022, 36092, 'ct-255-2', 'ct-255-2@201.3.116.2', 'ct-2', x'C9037402', 7, 255),
	(24023, 36093, 'ct-255-3', 'ct-255-3@201.4.23.3', 'ct-3', x'C9037403', 7, 255),
	(24024, 36094, 'ct-263-1', 'ct-63-1@201.77.23.3', 'ct-4', x'C9022403', 7, 263)
`
var DeleteNqmAgentSQL = `DELETE FROM nqm_agent WHERE ag_id >= 24021 AND ag_id <= 24025`

var InsertNqmtargetSQL = `
INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_available, tg_status)
VALUES
  (80921, 'pl-tg-1', '10.81.7.1', TRUE, TRUE),
  (80922, 'pl-tg-2', '10.81.7.2', TRUE, TRUE),
  (80923, 'pl-tg-3', '10.81.7.3', TRUE, TRUE)
`
var DeleteNqmtargetSQL = "DELETE FROM nqm_target WHERE tg_id >= 80921 AND tg_id <= 80923"

var InsertNqmPingtaskSQL = `
INSERT INTO nqm_ping_task(
  pt_id, pt_period
)
VALUES
  (83051, 1)
`
var DeletetNqmPingtaskSQL = `DELETE FROM nqm_ping_task WHERE pt_id = 83051`

var InsertNqmAgentPingtaskSQL = `
INSERT INTO nqm_agent_ping_task(apt_ag_id, apt_pt_id)
VALUES
  (24022, 83051)
`
var DeleteNqmAgentPingtaskSQL = "DELETE FROM nqm_agent_ping_task WHERE apt_ag_id = 24022"

var InitNqmAgent = []string{SetAutoIncForHost, SetAutoIncForNqmAgent, InsertHostSQL, InsertNqmAgentSQL}
var ClearNqmAgent = []string{DeleteNqmAgentSQL, DeleteHostSQL, ResetAutoIncForNqmAgent, ResetAutoIncForHost}
