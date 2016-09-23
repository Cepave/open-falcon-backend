CREATE TABLE `common_config` (
    `key` VARCHAR(255) NOT NULL DEFAULT '',
    `value` VARCHAR(255) NOT NULL DEFAULT ''
)
    ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `common_config`(`key`, `value`)
VALUES('git_repo', 'https://gitlab.com/Cepave/OwlPlugin.git');

INSERT INTO `common_config`(`key`, `value`)
VALUES('atom_addr', 'https://gitlab.com/Cepave/OwlPlugin/commits/master.atom');

