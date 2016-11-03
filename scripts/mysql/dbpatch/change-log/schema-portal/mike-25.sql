INSERT INTO `host`(hostname, ip, agent_version, plugin_version)
SELECT nqm_agent.ag_name, INET_NTOA(CONV(HEX(ag_ip_address), 16, 10)), '', ''
FROM nqm_agent
	LEFT OUTER JOIN
	`host`
	ON nqm_agent.ag_hostname = `host`.hostname
WHERE `host`.id IS NULL;

ALTER TABLE nqm_agent
ADD COLUMN ag_hs_id INT AFTER ag_id;

UPDATE nqm_agent, `host`
SET ag_hs_id = `host`.`id`
WHERE `host`.hostname = nqm_agent.ag_hostname;

ALTER TABLE nqm_agent
MODIFY COLUMN ag_hs_id INT NOT NULL,
ADD CONSTRAINT fk_nqm_agent__host
	FOREIGN KEY (ag_hs_id)
	REFERENCES `host`(`id`)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT;
