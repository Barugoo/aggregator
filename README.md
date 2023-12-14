>I was pretty short on time so i decided to skip implementing:
>- Graceful shutdown of server
>- Third-party logger for structured JSON logs
>- 100% test coverage

>These are NTH (IMHO) and can be added later.

# Server
This program implements HTTP-server as well as `aggregator` lib to execute aggregations on input data in a flexible way

Required envs:  
- `SERVER_ADDR` - this is the addr on which http server runs. You should pick `0.0.0.0:<port>` if you run this in docker. 

To run the server just execute:
```bash
make run
```
You can access the only endpoint of the server at route `/analyze?nweeks={ N past weeks to do aggregation on }`

To run unit-tests use:
```bash
make test
```

