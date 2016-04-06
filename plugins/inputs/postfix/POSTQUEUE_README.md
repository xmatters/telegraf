# Postqueue Input Plugin

The postqueue plugin gathers info about the Postfix mail queue.

Metrics
```
total_count 205
The total number of mails in the queue.
```

### Configuration:

```toml
# Get metrics regarding the Postfix mail queue.
[[inputs.postqueue]]
  # no configuration
```

### Measurements & Fields:

- postqueue
    - total_count (integer, `total_count`)

### Tags:

None

### Example Output:

```
$ telegraf -config ~/ws/telegraf.conf -input-filter postqueue -test
* Plugin: postqueue, Collection 1
> postqueue total_count=205i 1459967023286302545
```
