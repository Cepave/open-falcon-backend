# Dashboard API list
`Don't forget do URL encoding.. will check session automatically`
### `GET` `POST` /api/v1/dashboard/endpoints
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
### `GET` `POST` /api/v1/dashboard/endpointcounters
  * `required login session`
  * params:
    * `endpoints` []string (list of host name)
      * ex. ["docker-agent","testmachine"]
    * `metricQuery` string (Regex Query)
      * ex. "net.+"
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

### `GET` `POST` /api/v1/dashboard/endpointplugins
  * `required login session`
  * `required root login`
  * params:
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "GET",
        "status": "success",
        "data": {
          "Endpoints": [
            {
              "id": 1433110,
              "hostname": "docker",
              "ip": "10.0.0.167",
              "agent_version": "5.1.4",
              "plugin_version": "12155256cec3926186de22e282e67f4ce11cdbf7",
              "maintain_begin": 0,
              "maintain_end": 0,
              "update_at": "2016-06-30T07:53:56Z"
            },
            {
              "id": 1433111,
              "hostname": "foo",
              "ip": "10.0.0.168",
              "agent_version": "5.1.4",
              "plugin_version": "e5dd60e31698471431546a9a96434053adaa6c59",
              "maintain_begin": 0,
              "maintain_end": 0,
              "update_at": "2016-06-03T07:33:26Z"
            },
            ...
          ],
        "SessionFlag": false
        }
      }
      ```
      * Endpoints: the plugins information list of the machines.
    * failed

      ```
      {
        "version": "v1",
        "method": "GET",
        "status": "failed",
        "error": {
          "message": "name or sig is empty, please check again"
        }
      }
      ```

### `GET` `POST` /api/v1/dashboard/endpointrunningplugins
  * `required login session`
  * `required root login`
  * params:
    * `addr` string
      * ex. "http://10.0.0.167:1988/plugins"
  * response:
    * ok

      ```
        {
          "version": "v1",
          "method": "GET",
          "status": "success",
          "data": {
            "dataFromAgent": {
              "basic/chk/120_net_ping_gateway_loss.sh": {
                "Cycle": 120,
                "FilePath": "basic/chk/120_net_ping_gateway_loss.sh",
                "MTime": 1.468237012e+09
              },
              "basic/chk/60_check_heka_file.sh": {
                "Cycle": 60,
                "FilePath": "basic/chk/60_check_heka_file.sh",
                "MTime": 1.467890476e+09
              },
              ...
            },
            "msgFromAgent": "success",
            "requestAddr": "http://10.0.0.167:1988/plugins"
          }
        }
      ```
      * dataFromAgent: the running plugin information of the target machine.
    * ok

      ```
        {
          "version": "v1",
          "method": "GET",
          "status": "success",
          "data": {
            "errorFromAgent": "Get http://10.0.0.167:1988/plugins: dial tcp 10.0.0.167:1988: i/o timeout",
            "requestAddr": "http://10.0.0.167:1988/plugins"
          }
        }
      ```
      * errorFromAgent: the error message of the target machine.
    * failed

      ```
      {
        "version": "v1",
        "method": "GET",
        "status": "failed",
        "error": {
          "message": "name or sig is empty, please check again"
        }
      }
      ```

### `GET` `POST` /api/v1/dashboard/latestplugin
  * params:
  * response:
    * ok

      ```
        {
          "version": "v1",
          "method": "GET",
          "status": "success",
          "data": {
            "latestCommitHash": "9a4d709fd4e6511441d96281a1d5e392afba40b4"
          }
        }
      ```
### `GET` `POST` /api/v1/dashboard/counters
  * `required login session`
  * params:
    * `queryStr` string  [regex query string]
      * ex. `.+`
    * `limit`   integer [the maximum number of output]
      * ex. 20
  * response:
    * ok
    ```
      {
      "version": "v1",
      "method": "GET",
      "status": "success",
      "data": {
        "counters": [
          "check.heka.file",
          "check.heka.alived"
        ]
      }
    ```
    * failed
    ```
      {
        "version": "v1",
        "method": "GET",
        "status": "failed",
        "error": {
          "message": "query string is empty, please check it"
        }
      }
    ```

### `GET` `POST` /api/v1/dashboard/counterendpoints
  * `required login session`
  * params:
    * `counters` string array [the counter string array]
      * ex. `["agent.alive", "task"]`
    * `limit`   integer [the maximum number of output]
      * ex. 20
    * `filter` string [regex string describing the output]
      * ex. "bgp-.*"
  * response:
    * ok

      ```
        {
          "version": "v1",
          "method": "GET",
          "status": "success",
          "data": {
            "endpoints": [
              "bgp-bj-jjj",
              "bgp-bj-kkk",
              "bgp-bj-hhh",
              ...
            ]
          }
        }
      ```
      * endpoints: the list of endpoints that contains the counter in `counters` string array
    * failed

      ```
      {
        "version": "v1",
        "method": "GET",
        "status": "failed",
        "error": {
          "message": "query string counters is empty, please check it"
        }
      }
      ```

### `GET` `POST` /api/v1//hostgroup/query
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
### `GET` `POST` /api/v1//hostgroup/hosts
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

### `GET` `POST` /api/v1//hostgroup/hostgroupscounters
  * `required login session`
  * params:
    * `hostgroups` list of hostgroup
      * ex. ["docker-agent","testmachine"]
    * `metricQuery` string (Regex Query)
      * ex. "net.+"

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

