# Portal API list
`Don't forget do URL encoding.. will check session automatically`
* `GET` `POST` /api/v1/portal/eventcases/get
* `required login session`
* params:
  * `startTime` timestamp
    * ex: 1457450919
    * if not specific, means get all
  * `endTime` timestamp
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
    * ex: "PROBLEM", "in process", "unresolved", "resolved", "ignored"
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

* `GET` `POST` /api/v1/portal/events/get
  * `required login session`
  * params:
    * `startTime` timestamp
      * ex: 1457450919
    * `endTime` timestamp
    * `priority` int
      * ex: 0
      * -1 means no specific any priority level, get all.
    * `status` string options
      * ex: "PROBLEM"
      * "ALL" means no specific any status, get all.
      *  if no specific will get 'ALL' case as default
    * `cName` options (get from cookie refer)
      * default will only get case of current user, except admin role (ex. root).
    * `elimit` int
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
* `GET` `POST` /api/v1/portal/eventcases/close
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

* `GET` `POST` /api/v1/portal/eventcases/addnote
* add note for one event case
* `required login session`
* params:
  * `id` string
    * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
  * `note` string
    * max 300 varchar
  * `status` string options
    * ex: "in processing"
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

* `GET` `POST` /api/v1/portal/eventcases/notes
* get notes of one event case.
* `required login session`
* params:
  * `id` string
    * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
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
* `GET` `POST` /api/v1/portal/eventcases/note
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
