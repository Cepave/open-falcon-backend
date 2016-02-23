# API URL doc

## API list
* `POST` /api/v1/auth/register
  * params:
    * `name` string [useraccount]
    * `password` string
    * `repeat_password` string
    * `email` string
  * response:
    * ok

      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Success",
        "data": {
          "expired": 1458815516,
          "sig": "a9de7114da1811e5adb7001500c6ca5a"
        }
      }
      ```
      * sig: user session token
      * expired: token expired time
    * failed
      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Failed",
        "error": {
          "message": "name is already existent"
        }
      }
      ```
