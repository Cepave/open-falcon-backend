# UIC API list
`Don't forget do URL encoding..`
* `POST` /api/v1/auth/register
  * params:
    * `name` string [useraccount]
    * `password` string
    * `repeatPassword` string
    * `email` string
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
	  "name": "cepave",
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
        "version": "v1",
        "method": "POST",
        "status": "failed",
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
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
	  "name": "cepave",
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
        "version": "v1",
        "method": "POST",
        "status": "failed",
        "error": {
          "message": "password error"
        }
      }
      ```
* `GET` `POST` /api/v1/auth/logout
  * params:
    * `cName` string [useraccount] option
    * `cSig` string option
  * if params are null, will check [name & sig]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "message": "session is removed"
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
          "message": "name or sig is empty, please check again"
        }
      }
      ```
* `GET` `POST` /api/v1/auth/sessioncheck
  * params:
    * `cName` string [useraccount] option
    * `cSig` string option
  * if params are null, will check [name & sig]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "expired": 1458889692,
          "message": "session passed!",
          "sig": "5dda5bacdac511e5b83b001500c6ca5a"
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
          "message": "can not find this kind of session"
        }
      }
      ```
* `POST` /api/v1/auth/user
  * params:
    * `cName` string [useraccount] option
    * `cSig` string option
  * if params are null, will check [name & sig]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
        "status": "success",
        "data": {
          "cnname": "masato",
          "email": "",
          "im": "masato",
          "name": "cepavetest",
          "qq": "masato",
          "phone": ""
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
          "message": "can not find this kind of session"
        }
      }
      ```
* `POST` /api/v1/auth/user/update
  * params:
    * `cName` string [useraccount] option
    * `cSig` string option
    * `email` string option
    * `password` string option
      * if this is inputed, the original password of current user is required.
    * `oldpassword` string option
      * current password of this user
    * `cnname` string option
    * `im` string option
    * `qq` string option
    * `phone` string option
  * if params are null, will check [name & sig]
  * response:
    * ok

      ```
      {
        "version": "v1",
        "method": "POST",
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
          "message": "can not find this kind of session"
        }
      }
      ```
