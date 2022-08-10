# Tests execution by sending a request to redis

import json
from uuid import uuid4 as uuid
import redis

# Connect to redis
r = redis.Redis(host='localhost', port=6379, db=0)
file_name = "simple.js"

# Load file from disk
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
                "vus": 1,
                "duration": '1s',
            },
            "news": {
                "executor": 'per-vu-iterations',
                "exec": 'news',
                "vus": 1,
                "iterations": 1,
                "startTime": '0s',
                "maxDuration": '5s',
            },
        },
    })
}

# Add job to redis
print(f"Adding job id {id} to redis")

for key, value in job.items():
    r.hset(id, key, value)

r.publish('k6:execution', id)

# Add to history in case no worker none is listening
r.sadd('k6:executionHistory', id)

print(f"Job {id} added to redis")

# Listen for updates on the job
print(f"Listening for updates on at:", f"k6:executionUpdates:{id}")

while True:
    sub = r.pubsub()
    sub.subscribe(f"k6:executionUpdates:{id}")
    for message in sub.listen():
        try:
            print(json.loads(message['data']))
        except Exception as e:
            print(e)