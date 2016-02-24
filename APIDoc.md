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
* `POST` /api/v1/auth/login
  * params:
    * `name` string [useraccount]
    * `password` string
  * response:
    * ok

      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Success",
        "data": {
          "expired": 1458875558,
          "sig": "75981f64daa411e59eaa001500c6ca5a"
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
          "message": "password error"
        }
      }
      ```
* `POST` /api/v1/auth/logout
  * params:
    * `name` string [useraccount]
    * `token` string
  * response:
    * ok

      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Success",
        "data": {
          "message": "Session is deleted."
        }
      }
      ```
    * failed
    
      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Failed",
        "error": {
          "message": "name or token is empty, please check again"
        }
      }
      ```
* `POST` /api/v1/auth/sessioncheck
  * params:
    * `name` string [useraccount]
    * `token` string
  * response:
    * ok

      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Success",
        "data": {
          "expired": 1458889692,
          "message": "this token is works!",
          "token": "5dda5bacdac511e5b83b001500c6ca5a"
        }
      }
      ```
      * token: user session token
      * expired: token expired time
    * failed
    
      ```
      {
        "value": "v1",
        "method": "POST",
        "status": "Failed",
        "error": {
          "message": "can not find this kind of session"
        }
      }
      ```
