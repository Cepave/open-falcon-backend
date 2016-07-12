falcon-alarm
============

judge把报警event写入redis，alarm从redis读取event，做相应处理，可能是发报警短信、邮件，可能是callback某个http地址。
生成的短信、邮件写入queue，sender模块专门负责来发送。


## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/alarm.git
cd alarm
go get ./...
./control build
./control start
```

## Configuration

- uicToken: 留空即可
- http: 监听的http端口
- queue: 要发送的短信、邮件写入的队列，需要与sender配置一致
- redis: highQueues和lowQueues区别是是否做报警合并，默认配置是P0/P1不合并，收到之后直接发出；>=P2做报警合并
- api: 其他各个组件的地址

## Create Event Table
* use falcon_portal

```
CREATE TABLE event_cases (
        id VARCHAR(50),
        endpoint VARCHAR(100) NOT NULL,
        metric VARCHAR(200) NOT NULL,
        func VARCHAR(50),
        cond VARCHAR(200) NOT NULL,
        note VARCHAR(200),
        max_step int(10) unsigned,
        current_step int(10) unsigned,
        priority INT(6) NOT NULL,
        status VARCHAR(20) NOT NULL,
        timestamp Timestamp NOT NULL,
        update_at Timestamp,
        process_note MEDIUMINT,
        process_status VARCHAR(20),
        tpl_creator VARCHAR(64),
        expression_id int(10) unsigned,
        strategy_id int(10) unsigned,
        template_id int(10) unsigned,
        PRIMARY KEY (id),
        INDEX (endpoint, strategy_id, template_id)
);

CREATE TABLE events (
  id MEDIUMINT NOT NULL AUTO_INCREMENT,
  event_caseId VARCHAR(50),
  step int(10) unsigned,
  cond VARCHAR(200) NOT NULL,
  status  int(3) unsigned DEFAULT 0,
  timestamp Timestamp,
  PRIMARY KEY (id),
  INDEX(event_caseId),
  FOREIGN KEY (event_caseId) REFERENCES event_cases(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE TABLE event_note (
  id MEDIUMINT NOT NULL AUTO_INCREMENT,
  event_caseId VARCHAR(50),
  note    VARCHAR(300),
  case_id VARCHAR(20),
  status VARCHAR(15),
  timestamp Timestamp,
  user_id int(10) unsigned,
  PRIMARY KEY (id),
  INDEX (event_caseId),
  FOREIGN KEY (event_caseId) REFERENCES event_cases(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES uic.user(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
```
