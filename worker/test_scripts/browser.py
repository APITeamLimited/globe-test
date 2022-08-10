# Tests execution by sending a request to redis

import json
from uuid import uuid4 as uuid
import redis

# Connect to redis
r = redis.Redis(host='localhost', port=6379, db=0)
file_name = "browser.js"

# Load file from disk
with open(file_name, 'r') as f:
    file = f.read()

id = str(uuid())

job = {
    'id': id,
    "sourceName": file_name,
    "source": str(file),
    "status": "pending",
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
print(f"Listening for updates on job {id}:")

while True:
    sub = r.pubsub()
    sub.subscribe(f"k6:executionUpdates:{id}")
    for message in sub.listen():
        if message is not None:
            print(message)