# Workshop

## Deliverables

A sketch of a front end web-page that shows query results on the aggregated data in with
graphical tools. [Link](https://st0rrer.hotgloo.io/share/cpnQdegasCpbFWp)

### Client

Benchmark application which collect reports of activity and visit types with specified format and send them to web service

Arguments:
1) `concurrent-clients`: number of clients who will send the reports concurrently
2) `report-interval`: interval between sending the reports
3) `web-service`: URL to which to send reporting request
   
### Web-Service

Http server which listening on specified port, and handle POST request `/api/visit/v1/`&`/api/activity/v1/`
Received data push to the certain kafka topic.

Arguments:
1) `activity-topic`: topic name for `/api/visit/v1/` endpoint, if topic is absent in kafka will be created
2) `visit-topic`: topic name for `/api/visit/v1/` endpoint, if topic is absent in kafka will be created
3) `brokers value`: host and port one of the kafka cluster 
4) `port`:

### Data Collector

Application polls the messages from kafka, and saving them to file
Arguments:
1) `brokers`: host and port one of the kafka cluster
2) `concurrent-count`: number of processing operation which will be handle report concurrently
3) `groups`:  group and topic name for which start handle kafka messages
4) `record-directory` directory in which saving reports


### Launching
1) Create new network using command ```docker network create workshop-net```

2) Go to the kafka directory and start kafka cluster via docker-compose using command ```docker-compose up -d -V```

3) In project root directory build docker images using command: ```docker-compose build --no-cache```

4) Launch application using command ```docker-compose up -d```. There is an opportunity to specified programs arguments in 
docker-compose.yml 

## Extra Credit

1. What can be improved in your system? How would you do it?
    - Handle write errors. When failed to write messages to file, application need to handle these messages and send them to retry topic that would write in some time.
2. Is your system scalable? How? 
    - Data collector can scale horizontally, for that you need to launch new instance of Data collector application
  with the same groups config. If not each record from kafka will be broadcast to instance with different groups config.
    - Web service can be scale horizontally but for it needed to configure load balancing server for example `Nginx`.
    Also you can use service discovery for dynamic register new instance in load balancing server
    
    Each of the above services can scale vertically also. For it, needed to increase  involving the addition of CPUs, memory or storage
3. What is the bottleneck in the system? Why?  
    - Data collector has some bottlenecks such as limiting disk IO, and network. There may be a problem during recording in 
    file due to limit of file descriptors
    - Web services also has bottlenecks in network, and memory. There may be problems during handle 
    http request because application writes to kafka synchronized and 
    if connection to kafka is bad, it can take some time and memory for http connection
 
4. Can the system be debugged easily? Can its debuggability be improved?
   System is not debugged easily, because use ````golang.org/pkg/log/````, which doesn't have any logging levels. 
   For improved debuggability can add profiling for go application ```pprof```. Which allows debug some memory problems in an application.
   Also to structure logging message using log libraries such as ```logrus```