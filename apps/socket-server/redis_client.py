import os
import uuid
import redis.asyncio as aioredis

INSTANCE_ID = str(uuid.uuid4())

_client: aioredis.Redis | None = None

def get_redis() -> aioredis.Redis:
    global _client
    if _client is None:
        url = os.getenv("REDIS_URL", "redis://localhost:6379")
        password = os.getenv("REDIS_PASSWORD") or None
        _client = aioredis.from_url(url, password=password, decode_responses=True)
    return _client
