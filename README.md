<h2>
K6 Worker
</h2>

K6 Worker exposes a redis client that can accept and run multiple test jobs concurrently and supports streaming and reporting of results back to redis. The setup is aimed to make it easy to run tests in a distributed and headless environment.

The underlying K6 execution engine has been changed as little as possible to support the K6 Javascript API.
For more information on K6, please see the upstream repo <a>https://github.com/grafana/k6</a>

NOTE: This is very very early stage and not functional.

<h3>
Usage
</h3>

To Start a worker, run the following command:

```
go run main.go worker
```

Or similar equivalent in docker etc. this will spin up a worker that will listen on the default redis port.


Currently only localhost redis hardcoded

Example python queue script, replace with your own distribution system as you see fit

```python
# Tests execution by sending a request to redis

import json
from uuid import uuid4 as uuid
import redis

# Connect to redis
r = redis.Redis(host='localhost', port=6379, db=0)
file_name = "demo.js"

# Load demo.js from disk
with open(file_name, 'r') as f:
    file = f.read()

id = str(uuid())

job = ***REMOVED***
    'id': id,
    "sourceName": file_name,
    "source": str(file),
    "status": "pending",
    "options": json.dumps(***REMOVED***
        "scenarios": ***REMOVED***
            "contacts": ***REMOVED***
                "executor": 'constant-vus',
                "exec": 'contacts',
                "vus": 50,
                "duration": '30s',
            ***REMOVED***,
            "news": ***REMOVED***
                "executor": 'per-vu-iterations',
                "exec": 'news',
                "vus": 50,
                "iterations": 100,
                "startTime": '30s',
                "maxDuration": '1m',
            ***REMOVED***,
        ***REMOVED***,
    ***REMOVED***)
***REMOVED***

# Add job to redis
print(f"Adding job id ***REMOVED***id***REMOVED*** to redis")

for key, value in job.items():
    r.hset(id, key, value)

r.publish('k6:execution', id)

print(f"Job ***REMOVED***id***REMOVED*** added to redis")

# Listen for updates on the job
print(f"Listening for updates on job ***REMOVED***id***REMOVED***:")

while True:
    sub = r.pubsub()
    sub.subscribe(f"k6:executionUpdates:***REMOVED***id***REMOVED***")
    for message in sub.listen():
        if message is not None:
            try:
                print(message['data'].decode('utf-8'))
            except Exception as e:
                pass
```

Example script, options must be specified in redis job config, not the script

```javascript
import http from 'k6/http';
import ***REMOVED*** sleep ***REMOVED*** from 'k6';

export function contacts() ***REMOVED***
  const res = http.get('https://test.k6.io/contacts.php', ***REMOVED***
    tags: ***REMOVED*** my_custom_tag: 'contacts' ***REMOVED***,
  ***REMOVED***);
  console.log('contacts');
  sleep(1);
***REMOVED***

export function news() ***REMOVED***
  const res = http.get('https://test.k6.io/news.php', ***REMOVED*** tags: ***REMOVED*** my_custom_tag: 'news' ***REMOVED*** ***REMOVED***);
  sleep(1);
***REMOVED***
```