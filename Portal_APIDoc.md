# Portal API list
`Don't forget do URL encoding.. will check session automatically`
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
      * ex: "PORBLEM"
      * "ALL" means no specific any status, get all.
    * `if you keep all params is empty`, will get all event and status matched "PROBLEM".
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
              "id": "s_1_b66b973ef551e4e503fad475dfc9e418",
              "endpoint": "docker-agent",
              "metric": "cpu.idle",
              "func": "all(#1)",
              "cond": "97.91666666666667 != 0",
              "note": "",
              "max_step": 30000,
              "current_step": 28,
              "priority": 0,
              "status": "PROBLEM",
              "timestamp": "2016-03-09T21:42:00+08:00",
              "update_at": "2016-03-09T23:37:00+08:00",
              "closed_at": "0001-01-01T00:00:00Z",
              "user_modified": "",
              "expression_id": 0,
              "strategy_id": 1,
              "template_id": 1
            },
            {
              "id": "s_2_f0b61e6805e88bd507c6a48ebe43aea6",
              "endpoint": "docker-agent",
              "metric": "cpu.nice",
              "func": "all(#1)",
              "cond": "0 == 0",
              "note": "",
              "max_step": 30000,
              "current_step": 28,
              "priority": 0,
              "status": "PROBLEM",
              "timestamp": "2016-03-09T21:42:00+08:00",
              "update_at": "2016-03-09T23:37:00+08:00",
              "closed_at": "0001-01-01T00:00:00Z",
              "user_modified": "",
              "expression_id": 0,
              "strategy_id": 2,
              "template_id": 1
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
          "message": "query string is empty, please it"
        }
      }
      ```
* `GET` `POST` /api/v1/portal/events/close
  * `required login session`
  * params:
    * `id` string
      * ex: "s_3_ac9b8a08a4cb2def0320fec7ebecf8c8"
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
          "message": "query string is empty, please it"
        }
      }
      ```
