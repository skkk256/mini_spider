# Mini Spider

A simple spider developed with Go.
It's only for personal learning practice

- `go run main.go -h` show help
- load config from `conf/`, 
- supports parallel crawling of multiple routines

# Structure
```
├── src
│   ├── config       # config loader 
│   ├── request      # encapsulated request class
│   ├── response     # encapsulated response
│   ├── downloader   # Client that actually makes the request
│   ├── parser       # Parses the results in the response
│   ├── queue        # task queue and a hash table records processed url
│   ├── manager      # record the current number of unprocessed tasks
│    ...
├── conf
├── data
├── output
├── log


```