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
            parsed = json.loads(message['data'])['message']

            # Check if 'msg=' is in the message
            if 'msg=' in parsed:
                # Get msg and make sure nothing else besides msg is in the message
                msg = json.loads(json.loads(parsed.split('msg=')[1].split(' source=console')[0]))
                
                print(msg)
        except Exception as e:
            print(e)