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
        "Iterations": 100,
        "VUs": 1,
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
            try :
                print(message['data'].decode('utf-8'))
            except Exception as e:
                pass
    
