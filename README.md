# http-checker

Check http endpoint via GET method. Base on the configuration create metric.

Configuration file is yaml base:
```yaml
check:
  - url: https://google.com # url to check
    timeout: 3              # timeout to evaluate url ping
    check-period: 10        # period to check url
    metric: google_status   # prometheus metrics name
    metric-description: "Google reachability"
  - url: https://facebook.com
    timeout: 3
    check-period: 10
    metric: facebook_status
    metric-description: "Facebook reachability"
metrics:
  port: 8080                # port for exposing metrics

```
