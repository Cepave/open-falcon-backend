SET NAMES 'utf8';

DROP TABLE IF EXISTS event;
CREATE TABLE IF NOT EXISTS event
(
    id VARCHAR(50),
    endpoint VARCHAR(100) NOT NULL,
    metric VARCHAR(200) NOT NULL,
    func VARCHAR(50),
    cond VARCHAR(200) NOT NULL,
    note VARCHAR(500),
    max_step int(10) unsigned,
    current_step int(10) unsigned,
    priority INT(6) NOT NULL,
    status VARCHAR(20) NOT NULL,
    timestamp Timestamp NOT NULL,
    update_at Timestamp,
    closed_at Timestamp,
    user_modified int(10) unsigned,
    expression_id int(10) unsigned,
    strategy_id int(10) unsigned,
    template_id int(10) unsigned,
    PRIMARY KEY (id),
    INDEX (endpoint),
    FOREIGN KEY (strategy_id) REFERENCES strategy(id),
    FOREIGN KEY (template_id) REFERENCES tpl(id)
)
  ENGINE =InnoDB
  DEFAULT CHARSET =utf8;
