# Dashboard API list
`Don't forget do URL encoding..`
`will check session automatically`
* `GET` `POST` /api/v1/dashbaord/endpoints
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
      * sig: user session token
      * expired: token expired time
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
