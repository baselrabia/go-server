# Go HTTP server 

Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds. The server should continue to return the correct numbers after restarting it, by
persisting data to a file.

## Configuration

Configuration is managed via a `.env` file. Below is an example:

```bash
PORT=8080
DATA_FILE=data.json
WINDOW_DURATION=60s
PERSIST_INTERVAL=3000ms
```

- <b>port:</b> The port on which the server listens.
- <b> data_file:</b>  The path to the file where request data is persisted.
- <b> window_duration:</b>  The time window for counting requests.
- <b> persist_interval:</b>  The interval at which request data is persisted to the file.

## Run the Server

```bash
go run cmd/serve.go
``` 

## Serve Requests 
You can use curl or any web browser to make requests to the server on port `8080`
```bash
curl localhost:8080
```

## Unit tests

in order to run the unit tests you can run the command:

 ```bash
go test -race -v ./...
``` 

## Stress test  
Use the following command to send 10,000 requests concurrently to the server:

 ```bash
seq 1 10000 | xargs -P 10000 -I {} curl 0.0.0.0:8080
``` 


## Load tests

You can test the application with the help of [Vegeta](https://github.com/tsenart/vegeta) tool.

After running the server on port 8080, you can run the command:


 ```bash
echo "GET http://localhost:8080" | vegeta attack -rate=50000 -duration=5s| vegeta report
``` 


This command will send 50K request in 5 seconds to our application. You can change those values to test

### Output Example

``` bash 
╰─ echo "GET http://localhost:8080" | vegeta attack -duration=5s -rate=50000 | vegeta report
Requests      [total, rate, throughput]         72231, 14396.49, 13059.50
Duration      [total, attack, wait]             5.531s, 5.017s, 513.653ms
Latencies     [min, mean, 50, 90, 95, 99, max]  101.852µs, 633.451ms, 273.701ms, 1.599s, 1.792s, 1.89s, 3.2s
Bytes In      [total, mean]                     2805903, 38.85
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:72231  
Error Set:
```