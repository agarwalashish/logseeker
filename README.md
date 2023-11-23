# logseeker
This project has been developed as part of Cribl's take home assignment for the role of Senior Software Engineer.

## How to run

### System requirements 

1. The application has been dockerized and so the only dependency to run the application is docker. The application was built and tested using Docker version 4.10.1
2. To run the unit tests, go 1.21+ will be required

### Generating sample log data
There is a python script to generate sample log data which is in the /scripts folder. The script can be run using the command - 

```
python ./scripts/generate_logs.py
```

Enter the filename and the size of the file in the input prompts and the file will be created in the logs/ folder locally. When the docker image is compiled, the files will be copied into the `/var/log` folder of the container.

### Build and run the container
To build and run the container, run the command below -
```
docker-compose build && docker-compose up
```

## Testing
### API Endpoint
To search for logs, a POST request has to be sent as - 
```
curl --location 'http://localhost:8080/logs/search' \
--header 'Content-Type: application/json' \
--data '{
    "numLines": 500,
    "filename": "/var/log/large_file.log",
    "keywords": "System update"
}'
```

* `filename` is the file where logs will be searched and this is a required field
* `numLines` is the number of lines that will be returned. This is an optional field and by default, the last 10 lines will be returned
* `keywords` are the keywords that will be searched for. This is an optional field

### Running all unit tests
To run all unit tests, run the command -
```
go test ./...
```

## Future work

Future improvements to this project include - 

### Custom Parsers
This will allow users to provide a parser type (e.g `AWS VPC Flow Logs`) as an optional input parameter. The service can then attempt to parse the log file based on the file format and the client calling the API will be able to filter based on more fields.


### Authentication
Integrate the service with either a custom authentication service or an external one (e.g Auth0). Additionally, logs should also be separated for different organizations and a user should only be able to make the api requests to the organization that they are a member of.

### Streaming
Implement web sockets for streaming. This will allow the client to enter a filename and observe the events being written to it in real time.

