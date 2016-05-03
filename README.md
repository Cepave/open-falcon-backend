# Introduction

NQM is a module for [Open-Falcon](https://github.com/open-falcon/). It enables the feature of **N**etwork **Q**uality **M**easurement



# Makefile

*  `make`

   Build the binary

*  `make pack`

   Pack the necessary files into a tarball for deployment

*  `make clean`

   Remove the tarball and the excutable binary file





# Unit Test

> $ go test -v



# Configuration

You can modify `cfg.example.json` for creating your own configuration file:

```json
{
	"agentPushURL": "http://127.0.0.1:1988/v1/push",
	"hbsRPCServer": "127.0.0.1:6030",
	"hostname": "",
	"ipAddress": "",
	"connectionID": ""
}
```

Here are the explanations of the fields:

*  *agentPushURL*

  **Required**. The RESTful API URL where NQM agent pushes data to.


* *hbsRPCServer*

  **Required**. The RPC server where NQM agent gets configurations and probing commands.


* *hostname*

  **Optional**. If not set, NQM agent will use the system's hostname.

* *ipAddress*

  **Optional**. If not set, NQM agent will try to get the public IP address of the network interface. If failed, this field will be "UNKNOWN".

* *connectionID*

  **Optional**. If not set, NQM agent will generate a string combined by the hostname and the IP address.  



## Run

The default configuration file is `cfg.json`. You can run with the default configuration:

> $ ./nqm

Or, you can specify the configuration file by `-c`:

> $ ./nqm -c your.config.json
