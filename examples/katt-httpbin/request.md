

## Test this code

```
GET http://httpbin.org/get
> User-Agent: curl/7.64.1
> Accept: */*

< 200 OK
< Content-Type: application/json
< Connection: keep-alive
< Server: gunicorn/19.9.0
< Access-Control-Allow-Origin: *
< Access-Control-Allow-Credentials: true
{
  "args": {}, 
  "headers": {
    "Accept": "*/*", 
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.64.1"
  }, 
  "origin": "_", 
  "url": "http://httpbin.org/get"
}
```