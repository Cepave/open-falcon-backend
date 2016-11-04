/**
 * Remove non-linked data of host(in group-host table)
 */
DELETE gh
FROM grp_host AS gh
	LEFT OUTER JOIN
	host AS hs
	ON gh.host_id = hs.id
	LEFT OUTER JOIN
	grp AS gr
	ON gh.grp_id = gr.id
WHERE hs.id IS NULL OR
	gr.id IS NULL;

ALTER TABLE host
	DROP PRIMARY KEY,
	CHANGE COLUMN `id` hs_id_old INT(10) UNSIGNED NULL,
	AUTO_INCREMENT = 1,
	ADD COLUMN `id` INT AUTO_INCREMENT FIRST,
	ADD CONSTRAINT PRIMARY KEY (`id`);

ALTER TABLE grp_host
	CHANGE COLUMN host_id gh_hs_id_old INT(10) UNSIGNED NULL,
	DROP INDEX idx_grp_host_host_id,
	DROP INDEX idx_grp_host_grp_id,
	ADD COLUMN host_id INT AFTER grp_id;

UPDATE host AS hs, grp_host AS gh
SET gh.host_id = hs.id
WHERE gh.gh_hs_id_old = hs.hs_id_old;

ALTER TABLE grp_host
	ADD PRIMARY KEY(grp_id, host_id),
  ADD CONSTRAINT fk_grp_host__host FOREIGN KEY(host_id)
		REFERENCES host(`id`)
		ON DELETE CASCADE
		ON UPDATE RESTRICT,
	ADD CONSTRAINT fk_grp_host__grp FOREIGN KEY(grp_id)
		REFERENCES grp(`id`)
		ON DELETE CASCADE
		ON UPDATE RESTRICT;
