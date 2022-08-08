<h4>
K6 Worker
</h4>

K6 Worker exposes a redis client that can accept and run multiple test jobs concurrently and supports streaming and reporting of results back to redis. The setup is aimed to make it easy to run tests in a distributed and headless environment.

The underlying K6 execution engine has been changed as little as possible to support the K6 Javascript API.
For more information on K6, please see the upstream repo <a>https://github.com/grafana/k6</a>

NOTE: This is very very early stage and not functional

<h3>
Usage
</h3>

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

job = {
    'id': id,
    "sourceName": file_name,
    "source": str(file),
    "status": "pending",
    "options": json.dumps({
        "scenarios": {
            "contacts": {
                "executor": 'constant-vus',
                "exec": 'contacts',
                "vus": 50,
                "duration": '30s',
            },
            "news": {
                "executor": 'per-vu-iterations',
                "exec": 'news',
                "vus": 50,
                "iterations": 100,
                "startTime": '30s',
                "maxDuration": '1m',
            },
        },
    })
}

# Add job to redis
print(f"Adding job id {id} to redis")

for key, value in job.items():
    r.hset(id, key, value)

r.publish('k6:execution', id)

print(f"Job {id} added to redis")

# Listen for updates on the job
print(f"Listening for updates on job {id}:")

while True:
    sub = r.pubsub()
    sub.subscribe(f"k6:executionUpdates:{id}")
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
import { sleep } from 'k6';

export function contacts() {
  const res = http.get('https://test.k6.io/contacts.php', {
    tags: { my_custom_tag: 'contacts' },
  });
  console.log('contacts');
  sleep(1);
}

export function news() {
  const res = http.get('https://test.k6.io/news.php', { tags: { my_custom_tag: 'news' } });
  sleep(1);
}
```