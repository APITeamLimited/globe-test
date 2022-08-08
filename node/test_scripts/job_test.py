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
