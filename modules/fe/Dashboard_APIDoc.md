# Dashboard API list
`Don't forget do URL encoding.. will check session automatically`
* `GET` `POST` /api/v1/dashboard/endpoints
  * `required login session`
  * params:
    * `queryStr` string [regex query string]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "endpoints": [
            {
              "id": 0,
              "endpoint": "docker-task",
              "ts": 1456330140
            },
            {
              "id": 0,
              "endpoint": "docker-agent",
              "ts": 1456330080
            },
            {
              "id": 0,
              "endpoint": "10.0.0.167",
              "ts": 1456330140
            }
          ]
        }
      }
      ```
      * endpoints: matched endpoints
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
* `GET` `POST` /api/v1/dashboard/endpointcounters
  * `required login session`
  * params:
    * `endpoints` []string (list of host name)
      * ex. ["docker-agent","testmachine"]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "counters": [
            "agent.alive",
            "cpu.guest",
            "cpu.idle",
            "cpu.iowait",
            "cpu.irq",
            "cpu.nice",
            "cpu.softirq",
            "cpu.steal",
            "cpu.switches",
            "cpu.system",
            "cpu.user"
          ]
        }
      }
      ```
      * counters: available metrics of the machines list.
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

* `GET` `POST` /api/v1//hostgroup/query
  * `required login session`
  * params:
    * `queryStr` string [regex query string]
  * response:
    * ok

    ```
    {
      "version": "v1",
      "method": "POST",
      "status": "success",
      "data": {
        "hostgroups": [
          {
            "id": 2,
            "grp_name": "docker",
            "create_user": "cepavetest",
            "create_at": "2016-03-07T13:41:52+08:00",
            "come_from": 1
          }
        ]
      }
    }
    ```
    * hostgroups: matched hostgroups
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
* `GET` `POST` /api/v1//hostgroup/hosts
  * `required login session`
  * params:
    * `hostgroups` list of hostgroup
      * ex. ["docker-agent","testmachine"]

  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "hosts": [
            {
              "id": 1,
              "hostname": "docker-agent",
              "ip": "172.17.0.17",
              "agent_version": "5.1.0",
              "plugin_version": "plugin not enabled",
              "maintain_begin": 0,
              "maintain_end": 0,
              "update_at": "2016-01-04T16:43:37+08:00"
            },
            {
              "id": 28988,
              "hostname": "docker-task",
              "ip": "",
              "agent_version": "",
              "plugin_version": "",
              "maintain_begin": 0,
              "maintain_end": 0,
              "update_at": "2016-03-07T13:41:21+08:00"
            }
          ]
        }
      }
      ```
      * hosts: hosts list of hostgorup
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

* `GET` `POST` /api/v1//hostgroup/hostgroupscounters
  * `required login session`
  * params:
    * `hostgroups` list of hostgroup
      * ex. ["docker-agent","testmachine"]

  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "counters": [
            "falcon.task.alive",
            "agent.alive",
            "cpu.guest",
            "cpu.idle",
            "cpu.iowait",
            "cpu.irq",
            "cpu.nice",
            "cpu.softirq",
            "cpu.steal",
            "cpu.switches",
            "cpu.system",
            "cpu.user",
            "df.bytes.free.percent/fstype=ext4,mount=/config",
            "df.inodes.free.percent/fstype=ext4,mount=/config",
            "df.statistics.total",
            "df.statistics.used"
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
