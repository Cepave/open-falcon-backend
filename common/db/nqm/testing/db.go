package testing

var DeleteNqmAgentPingtaskSQL = `DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 24021 AND apt_ag_id <= 24025`

var InsertPingtaskSQL = `
INSERT INTO nqm_ping_task(pt_id, pt_name, pt_period)
VALUES(10119, 'test1-pingtask_name', 40),
(10120, 'test2-pingtask_name', 3)
`
var DeletePingtaskSQL = `DELETE FROM nqm_ping_task WHERE pt_id >= 10119`

var InsertHostSQL = `
INSERT INTO host(id, hostname, agent_version, plugin_version)
VALUES(36091, 'ct-agent-1', '', ''),
	(36092, 'ct-agent-2', '', ''),
	(36093, 'ct-agent-3', '', ''),
	(36094, 'ct-agent-4', '', '')
`
var DeleteHostSQL = `DELETE FROM host WHERE id >= 36091 AND id <= 36095`

var InsertOwlNameTag = `
INSERT INTO owl_name_tag(nt_id, nt_value)
VALUES(3081, '美國 IP 群')
`

var DeleteOwlNameTag = `DELETE FROM owl_name_tag WHERE nt_id = 3081`

var InsertOwlGroupTag = `
INSERT INTO owl_group_tag(gt_id, gt_name)
VALUES
	(23401, '串流一網'),
	(23402, '串流二 網')
`

var DeleteOwlGroupTag = `DELETE FROM owl_group_tag WHERE gt_id >= 23401 AND gt_id <= 23402`

var InsertNqmTarget = `
INSERT INTO nqm_target(tg_id, tg_name, tg_host, tg_available, tg_status, tg_isp_id, tg_nt_id)
VALUES
	(40201, 'tg-name-1', 'tg-host-1', false, true, 3, -1),
	(40202, 'tg-name-2', 'tg-host-2', true, true, 4, 3081),
	(40203, 'tg-name-3', 'tg-host-3', true, false, 5, -1)
`
var DeleteNqmTarget = `DELETE FROM nqm_target WHERE tg_id >= 40201 AND tg_id <= 40203`

var InsertNqmTargetGroupTag = `
INSERT INTO nqm_target_group_tag(tgt_tg_id, tgt_gt_id)
VALUES
	(40201, 23401),
	(40201, 23402)
`

var InitNqmTarget = []string{InsertOwlNameTag, InsertOwlGroupTag, InsertNqmTarget, InsertNqmTargetGroupTag}
var ClearNqmTarget = []string{DeleteNqmTarget, DeleteOwlNameTag, DeleteOwlGroupTag}

var InsertNqmAgentSQL = `
INSERT INTO nqm_agent(
	ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address,
	ag_pv_id, ag_ct_id
)
VALUES(24021, 36091, 'ct-255-1', 'ct-255-1@201.3.116.1', 'ct-1', x'C9037401', 7, 255),
	(24022, 36092, 'ct-255-2', 'ct-255-2@201.3.116.2', 'ct-2', x'C9037402', 7, 255),
	(24023, 36093, 'ct-255-3', 'ct-255-3@201.4.23.3', 'ct-3', x'C9037403', 7, 255),
	(24024, 36094, 'ct-263-1', 'ct-63-1@201.77.23.3', 'ct-4', x'C9022403', 7, 263)
`
var DeleteNqmAgentSQL = `DELETE FROM nqm_agent WHERE ag_id >= 24021 AND ag_id <= 24025`

var InitNqmAgentAndPingtaskSQL = []string{InsertPingtaskSQL, InsertHostSQL, InsertNqmAgentSQL}
var CleanNqmAgentAndPingtaskSQL = []string{DeleteNqmAgentPingtaskSQL, DeleteNqmAgentSQL, DeleteHostSQL, DeletePingtaskSQL}

var InitNqmAgent = []string{InsertHostSQL, InsertNqmAgentSQL}
var ClearNqmAgent = []string{DeleteNqmAgentSQL, DeleteHostSQL}

var InsertNqmCacheAgentPingListLog = `
INSERT INTO nqm_cache_agent_ping_list_log(
	apll_ag_id, apll_number_of_targets, apll_time_access, apll_time_refresh
)
VALUES
	(24021, 3, '2017-01-01', '2017-01-01'),
	(24022, 0, '2017-01-01', '2017-01-01')
`

var DeleteNqmCacheAgentPingListLog = `
DELETE FROM nqm_cache_agent_ping_list_log
WHERE apll_ag_id>= 24021 AND  apll_ag_id<= 24022
`

var InsertNqmCacheAgentPingList = `
INSERT INTO nqm_cache_agent_ping_list(
	apl_apll_ag_id, apl_tg_id, apl_min_period, apl_time_access
)
VALUES
	(24021, 40201, 1, '2017-01-01'),
	(24021, 40202, 1, '2017-01-01'),
	(24021, 40203, 1, '2017-01-01')
`

var InitNqmCacheAgentPingListLog = []string{InsertHostSQL, InsertNqmAgentSQL, InsertNqmCacheAgentPingListLog}

var ClearNqmCacheAgentPingListLog = []string{DeleteNqmCacheAgentPingListLog, DeleteNqmAgentSQL, DeleteHostSQL}

var InitNqmCacheAgentPingList = []string{InsertHostSQL, InsertNqmAgentSQL, InsertNqmCacheAgentPingListLog, InsertOwlNameTag, InsertOwlGroupTag, InsertNqmTarget, InsertNqmTargetGroupTag, InsertNqmCacheAgentPingList}

var ClearNqmCacheAgentPingList = []string{DeleteNqmCacheAgentPingListLog, DeleteNqmAgentSQL, DeleteHostSQL, DeleteNqmTarget, DeleteOwlNameTag, DeleteOwlGroupTag}
