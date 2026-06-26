import os
import uuid
import redis.asyncio as aioredis

INSTANCE_ID = str(uuid.uuid4())

_client: aioredis.Redis | None = None

def get_redis() -> aioredis.Redis:
    global _client
    if _client is None:
        url = os.getenv("REDIS_URL", "redis://localhost:6379")
        _client = aioredis.from_url(url, decode_responses=True)
    return _client
