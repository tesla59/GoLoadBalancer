# Go Load Balancer

The project was started as assessment for the organization [GoComet](https://www.gocomet.com/). Now the project is made open-source for references.

## What is GoLoadBalancer?
GoLoadBalancer is an automatic scaling load balancer with round-robin functionality.

The application consists of 3 parts (modules)
1. Load-balancer: 
   - Takes all web requests and pass it to available workers with round-robin functionality.

2. WebApp:
	- /api/v1/hello : This endpoint will respond with a JSON message after an delay. The delay is specified in config.yaml.
	- /worker/stats : This endpoint will provide the statistics of all workers. It is important to note that the statistics is updated and saved on each request made to the /api/v1/hello endpoint, and are stored in the database provided in the config.yaml.

3. Configuration Manager:
   - Config Manager reads config.yaml as in input and depending on the contents of the config file, it spawns the number of load balancer or workers. Also contains the average failure rate and average delay
   - Load Balancer are spawned as child processes and Workers are spawned as docker containers

## Config Manager
Config Manager is the entrypoint in the application and is defined as main package. It is responsible for reading the config file (config.yaml), spawn load balancer and spawn worker nodes.

On init, Config Manager first builds the Worker as docker image and then runs them in container. The number of worker images to spawn can be modified in config.yaml. The worker application is built using docker to showcase the replicability of the program and to mimic the Kubernetes behavior in some sorts

After Worker, Config Manager is responsible for running the load balancer. Multiple load balancer can also be spawned using goroutine. This behavior is kept to showcase the difference between running an application as process and running an application as docker container.

### Config.yaml
Used to change the behavior of application. Contains fields like:

```
worker: 3 # number of workers to spawn
pool: 2 # number of load balancer to spawn
stats-dir: worker-stats.db # location of the SQL DB
avg-delay: 500 # Average Delay of the API Req (in ms)
failure: 12 # Failure Rate of API Req
load-balancer-port: 5050 # Initial Port of Load Balancer: if more than one LB is created, port is incremented by 1
worker-port: 6000 # Initial Port of Worker: if more than one worker is created, port is incremented by 1
```

Config Manager also creates the database responsible for storing the worker stats.

## Load Balancer
Stored in root/loadbalancer folder as loadbalancer package. It is imported in Config Manager and run from there.
Follows Round Robin method to distribute load evenly
Runs health check at every 5s to keep list of healthy Server Pool updated

## Worker
A backend server written in Go. Uses inbuilt Go libraries. It has 2 endpoints

`/api/v1/hello` returns `{"message":"hello from container: <container id>"}`

`/worker/stats` returns the json object

```{
	"SuccessfulRequests": [
		{
			"total": 0
		},
		{
			"fdbaa57c12d3": 0
		},
		{
			"7fd03362f4f4": 0
		},
		{
			"c1623efb21e8": 0
		}
	],
	"FailedRequests": [
		{
			"total": 0
		},
		{
			"fdbaa57c12d3": 0
		},
		{
			"7fd03362f4f4": 0
		},
		{
			"c1623efb21e8": 0
		}
	],
	"TotalRequests": [
		{
			"total": 0
		},
		{
			"fdbaa57c12d3": 0
		},
		{
			"7fd03362f4f4": 0
		},
		{
			"c1623efb21e8": 0
		}
	],
	"AverageResponseTime": [
		{
			"total": 0
		},
		{
			"fdbaa57c12d3": 0
		},
		{
			"7fd03362f4f4": 0
		},
		{
			"c1623efb21e8": 0
		}
	]
}
```

Uses `math/rand` library to set the average response time and failure rates as mentioned in the config.yaml

# How to run?
1. Make sure `go` (V1.21.0) and `docker` is installed. And current user has the permission to run docker without root (i.e. `docker ps` should not throw permission error)
2. `go run .`
This should get all the libraries necessary and start the server with default configs

The architecture is inspired from K3s and tries to implement kubernetes operator
Notes:
1. Please delete worker-stats.db incase of any bug 

