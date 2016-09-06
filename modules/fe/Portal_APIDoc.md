# Portal API list
`Don't forget do URL encoding.. will check session automatically`
### `GET` `POST` /api/v1/portal/eventcases/get
* `required login session`
* params:
  * `startTime` timestamp [if set then can't skip endTime]
    * ex: 1457450919
    * if not specific, means get all
  * `endTime` timestamp [if set then can't skip startTime]
    * ex: 1477450919
    * if not specific, means get all
  * `priority` int
    * ex: 0
    * -1 means no specific any priority level, get all.
    * default: -1
  * `status` string options
    * ex: "PROBLEM", "OK"
    * "ALL" means no specific any status, get all.
    * support mutiple status query. ex: "PRBOEM,OK"
    * 'OK' means 'Recovery', this wording is from open-falcon[judge], so I didn't change it.
  * `process_status` string options
    * ex: "in progress", "unresolved", "resolved", "ignored"
    * support mutiple status query. ex: "resolved,ignored"
  * `metrics` string options
    * ex: "cpu.idle", "df.statistics.total"
    * support mutiple status query. ex: "cpu.idle,df.statistics.total,net.if.in.bits/iface=eth0"
  * `cName` options (get from cookie refer)
    * default will only get case of current user, except admin role (ex. root).
  * `limit` int
    * set the return limit of eventCases
  * `elimit` int
    * set the return limit of events
  * `caseId` string
    * return specific alarm_case
* response:
  * ok

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "success",
      "data": {
        "eventCases": [
          {
            "id": "s_2_3275011ce6ce8603f5f3917d3d060ea7",
            "endpoint": "docker-agent",
            "metric": "test.alarm",
            "func": "all(#2)",
            "cond": "1 != 0",
            "note": "test.alarm",
            "max_step": 3000,
            "current_step": 182,
            "priority": 0,
            "status": "PROBLEM",
            "start_at": "2016-06-03T12:51:00+08:00",
            "update_at": "2016-06-06T13:32:00+08:00",
            "process_note": 0,
            "process_status": "unresolved",
            "tpl_creator": "root",
            "expression_id": 0,
            "strategy_id": 2,
            "template_id": 1,
            "evevnts" [
              {
                "id": 1,
                "step": 1,
                "cond": "100 != 0",
                "timestamp": "2016-04-06T14:11:00+08:00",
                "event_caseId": null
              },
              {
                "id": 2,
                "step": 2,
                "cond": "100 != 0",
                "timestamp": "2016-04-06T14:16:00+08:00",
                "event_caseId": null
              }
            ]
          }
        ]
      }
    }
    ```

### `GET` `POST` /api/v1/portal/events/get
  * `required login session`
  * params:
    * `startTime` timestamp
      * ex: 1457450919
    * `endTime` timestamp
      * ex: 1457450919
    * `status` string options
      * ex. "OK", "PROBLEM"
    * `cName` options (get from cookie refer)
      * no permission control, will get top 500 events as default
    * `caseId` string
      * session control
    * `limit` int
      * set the return limit of events
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "events": [
            {
              "id": 3,
              "step": 3,
              "cond": "100 != 0",
              "timestamp": "2016-04-06T14:21:00+08:00",
              "event_caseId": "s_41_b66b973ef551e4e503fad475dfc9e418",
              "tpl_creator": "root",
              "metric": "cpu.idle",
              "endpoint": "docker-agent"
            },
            {
              "id": 2,
              "step": 2,
              "cond": "100 != 0",
              "timestamp": "2016-04-06T14:16:00+08:00",
              "event_caseId": "s_41_b66b973ef551e4e503fad475dfc9e418",
              "tpl_creator": "root",
              "metric": "cpu.idle",
              "endpoint": "docker-agent"
            },
            {
              "id": 1,
              "step": 1,
              "cond": "100 != 0",
              "timestamp": "2016-04-06T14:11:00+08:00",
              "event_caseId": "s_41_b66b973ef551e4e503fad475dfc9e418",
              "tpl_creator": "root",
              "metric": "cpu.idle",
              "endpoint": "docker-agent"
            }
          ]
        }
      }
      ```
    * failed

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "failed",
        "error": {}
      }
      ```

### `GET` `POST` /api/v1/portal/eventcases/close
* !! this is deprecated now !!
* `required login session`
* params:
  * `id` string
    * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
  * `closedNote` string
  * ok

    ```
    {
      "version": "v1",
      "method": "PUT",
      "status": "success"
    }
    ```
  * failed

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "failed",
      "error": {
        "message": "You can not skip closed note"
      }
    }
    ```

### `GET` `POST` /api/v1/portal/eventcases/addnote
* add note for one event case
* `required login session`
* params:
  * `id` string
    * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
    * support batch update -> ex. "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8,s_4_ac9b8a08a4cb2def0320fec7ebecf8c8,s_9_ac9b8a08a4cb2def0320fec7ebecf8c2"
  * `note` string
    * max 300 varchar
  * `status` string options
    * ex: "in progress", "unresolved", "resolved", "ignored"
  * `eventId` string options
    * boss case id
  * ok

    ```
    {
      "version": "v1",
      "method": "PUT",
      "status": "success"
    }
    ```
  * failed

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "failed",
      "error": {
        "message": "You can not skip closed note"
      }
    }
    ```

### `GET` `POST` /api/v1/portal/eventcases/notes
* get notes of one event case.
* `required login session`
* params:
  * `id` string
    * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
  * `status` string options
    * in events status means process_status
    * "ALL" means no specific any status, get all.
    * ex: "in progress", "unresolved", "resolved", "ignored"
    * support mutiple status query. ex: "resolved,ignored"
  * `filterIgnored` bool
    * default: false
  * @Timeiflter
    * `startTime` unixTime
    * `endTime` unixTime
      * Accpet only set startTime and will make the currentTime as the endTime
      * if startTime & endTime did not specific, will use the lasted #1 of that alarm as startTime , the currentTime will be the endTime.
  * ok

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "success",
      "data": {
        "notes": [
          {
            "id": 3,
            "event_caseId": "s_116_b66b973ef551e4e503fad475dfc9e418",
            "note": "測試結果",
            "case_id": "99667",
            "status": "inprocess",
            "timestamp": "2016-06-01T11:15:50+08:00",
            "user_name": "root"
          },
          {
            "id": 2,
            "event_caseId": "s_116_b66b973ef551e4e503fad475dfc9e418",
            "note": "測試結果",
            "case_id": "99667",
            "status": "inprocess",
            "timestamp": "2016-06-01T11:15:38+08:00",
            "user_name": "root"
          },
          {
            "id": 1,
            "event_caseId": "s_116_b66b973ef551e4e503fad475dfc9e418",
            "note": "測試結果",
            "case_id": "99667",
            "status": "inprocess",
            "timestamp": "2016-06-01T11:15:29+08:00",
            "user_name": "root"
          }
        ]
      }
    }
    ```
  * failed

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "failed",
      "error": {
        "message": "You dosen't pick any event id"
      }
    }
    ```

### `GET` `POST` /api/v1/portal/eventcases/note
* get one note
* `required login session`
* params:
  * `id` int
    * ex: 18
  * ok

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "success",
      "data": {
        "note": {
          "id": 18,
          "event_caseId": "s_2_3275011ce6ce8603f5f3917d3d060ea7",
          "note": "測試結果",
          "case_id": "99667",
          "status": "resolved",
          "timestamp": "2016-06-03T15:41:21+08:00",
          "user_name": "root"
        }
      }
    }
    ```
  * failed

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "failed",
      "error": {
        "message": "You dosen't pick any note id"
      }
    }
    ```

### `GET` `POST` /api/v2/portal/eventcases/get
### `GET` `POST` /api/v3/portal/eventcases/get
* `required login session`
* v2: mapping get from boss api, v3: mapping get from sql cache
* params:
  * `includeEvents` string (boolean)
    * ex. false, true
    * default: `false`
    * if true will response events information (time-series alerts data by step).
  * `startTime` timestamp [if set then can't skip endTime]
    * ex: 1457450919
    * if not specific, means get all
  * `endTime` timestamp [if set then can't skip startTime]
    * ex: 1477450919
    * if not specific, means get all
  * `priority` int
    * ex: 0
    * -1 means no specific any priority level, get all.
    * default: -1
  * `status` string options
    * ex: "PROBLEM", "OK"
    * "ALL" means no specific any status, get all.
    * support mutiple status query. ex: "PRBOEM,OK"
    * 'OK' means 'Recovery', this wording is from open-falcon[judge], so I didn't change it.
  * `process` string options
    * ex: "in progress", "unresolved", "resolved", "ignored"
    * support mutiple status query. ex: "resolved,ignored"
  * `metric` string options
    * ex: "cpu.idle", "df.statistics.total"
    * support mutiple status query. ex: "cpu.idle,df.statistics.total,net.if.in.bits/iface=eth0"
  * `cName` options (get from cookie refer)
    * default will only get case of current user, except admin role (ex. root).
  * `limit` int
    * set the return limit of eventCases
  * `elimit` int
    * set the return limit of events
  * `caseId` string
    * return specific alarm_case
  * `show_all` bool
    * filter alarm case of inactive endpoints
* response:
  * ok

    ```
    {
    "version": "v2",
    "method": "GET",
    "status": "success",
    "data": {
      "eventCases": [
        {
          "contact": [
            {
              "phone": "-",
              "email": "-",
              "name": "-"
            }
          ],
          "idc": "not found",
          "ip": "not found",
          "platform": "not found",
          "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
          "hostname": "docker-agent",
          "metric": "cpu.idle",
          "author": "cepave",
          "templateID": 66,
          "priority": "2",
          "severity": "Low",
          "status": "Triggered",
          "statusRaw": "PROBLEM",
          "type": "proc",
          "content": "cpu idle !!",
          "timeStart": "2016-07-19 15:46",
          "timeUpdate": "2016-07-19 15:56",
          "duration": "6 days ago",
          "notes": [
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-20 14:24",
              "user": "root"
            },
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-20 14:25",
              "user": "root"
            },
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-20 14:25",
              "user": "root"
            },
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-20 14:25",
              "user": "root"
            },
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-20 14:25",
              "user": "root"
            },
            {
              "hash": "s_104_1dsasadsadasdwqwdqwdwqdw",
              "note": "ignored",
              "status": "ignored",
              "time": "2016-07-24 18:00",
              "user": "root"
            }
          ],
          "events": [
            {
              "id": 813262,
              "step": 3,
              "cond": "0 \u003c 1",
              "status": 0,
              "timestamp": "2016-07-19T15:56:00+08:00",
              "event_caseId": null
            },
            {
              "id": 813207,
              "step": 2,
              "cond": "0 \u003c 1",
              "status": 0,
              "timestamp": "2016-07-19T15:51:00+08:00",
              "event_caseId": null
            },
            {
              "id": 813134,
              "step": 1,
              "cond": "0 \u003c 1",
              "status": 0,
              "timestamp": "2016-07-19T15:46:00+08:00",
              "event_caseId": null
            }
          ],
          "process": "unresolved",
          "function": "all(#1)",
          "condition": "0 \u003c 1",
          "stepLimit": 3,
          "step": 3
        }]
      }
    }
    ```

### `GET` `POST` /api/v2/portal/eventcases/feed
* get one note
* `required login session`
* if neede refresh alarm page, the any_new will be true
* params:
  * `cName` string
  * `cSig` string
  * ok

    ```
    {
    "version": "v2",
    "method": "GET",
    "status": "success",
    "data": {
      "admin": true,
      "any_new": true,
      events": [
      {
        "id": "s_237_0867fdsfdsfsdfdsf1a68cb",
        "endpoint": "docker-agent",
        "metric": "service.logs.test",
        "func": "all(#3)",
        "cond": "5 \u003e= 5",
        "note": "this is note",
        "max_step": 1,
        "current_step": 1,
        "priority": 4,
        "status": "PROBLEM",
        "start_at": "2016-08-02T15:12:00+08:00",
        "update_at": "2016-08-02T15:12:00+08:00",
        "process_note": 0,
        "process_status": "unresolved",
        "tpl_creator": "weiqs",
        "expression_id": 0,
        "strategy_id": 237,
        "template_id": 111,
        "evevnts": [
          {
            "id": 923543,
            "step": 1,
            "cond": "5 \u003e= 5",
            "status": 0,
            "timestamp": "2016-08-02T15:12:00+08:00",
            "event_caseId": null
          },
          {
            "id": 890743,
            "step": 1,
            "cond": "1 \u003e= 5",
            "status": 1,
            "timestamp": "2016-07-29T11:08:00+08:00",
            "event_caseId": null
          },
          {
            "id": 890734,
            "step": 1,
            "cond": "5 \u003e= 5",
            "status": 0,
            "timestamp": "2016-07-29T11:07:00+08:00",
            "event_caseId": null
          },
          {
            "id": 890720,
            "step": 1,
            "cond": "4 \u003e= 5",
            "status": 1,
            "timestamp": "2016-07-29T11:03:00+08:00",
            "event_caseId": null
          }]
      }],
      "notes": [
        {
          "id": 7358,
          "event_caseId": "s_288_db27b4ff411966bfe6c3dsfds4add58c",
          "note": "ignored",
          "case_id": "",
          "status": "ignored",
          "timestamp": "2016-08-02T12:09:50+08:00",
          "user_name": "root"
        },
        {
          "id": 7359,
          "event_caseId": "s_148_dcf9347sfd21d73ea0299dsf39bd0225",
          "note": "ignored",
          "case_id": "",
          "status": "ignored",
          "timestamp": "2016-08-02T12:09:50+08:00",
          "user_name": "root"
        },
        {
          "id": 7360,
          "event_caseId": "s_128_e7be175bfsdfdsfffec2e083cfd7d22a",
          "note": "ignored",
          "case_id": "",
          "status": "ignored",
          "timestamp": "2016-08-02T12:09:50+08:00",
          "user_name": "root"
        }]
      }
    }
    ```
  * failed

    ```
    {
      "version": "v2",
      "method": "GET",
      "status": "failed",
      "error": {
        "message": "can not find this kind of session"
      }
    }
    ```
