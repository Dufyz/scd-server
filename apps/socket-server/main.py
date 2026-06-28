import asyncio
import json
import os
import websockets
import sys
sys.stdout.reconfigure(line_buffering=True)
from dotenv import load_dotenv
load_dotenv()

from room_manager import remove_client, broadcast
from event_handler import handle_event
from redis_client import get_redis, INSTANCE_ID
from kafka_client import consume_kafka

PORT = int(os.getenv("PORT", 8765))

async def handler(websocket):
    print(f"Cliente conectado: {websocket.remote_address}")
    try:
        async for message in websocket:
            await handle_event(websocket, message)
    finally:
        remove_client(websocket)
        print(f"Cliente desconectado: {websocket.remote_address}")

async def process_request(path, request_headers):
    if path == "/health":
        body = {
            "message": "Health status endpoint is operational",
            "status": "up",
        }
        return 200, [("Content-Type", "application/json")], json.dumps(body).encode()
    return None

async def redis_subscriber():
    while True:
        try:
            redis = get_redis()
            pubsub = redis.pubsub()
            await pubsub.psubscribe("room:*")
            print("Redis subscriber iniciado no padrão: room:*")
            async for raw in pubsub.listen():
                if raw["type"] != "pmessage":
                    continue
                try:
                    data = json.loads(raw["data"])
                except (json.JSONDecodeError, TypeError):
                    continue

                # Ignora mensagens originadas nesta instância (já fizemos broadcast local)
                if data.get("instance_id") == INSTANCE_ID:
                    continue

                room_id = data.get("room_id")
                event = data.get("event")
                if room_id and event:
                    await broadcast(room_id, event)
        except Exception as e:
            print(f"Redis subscriber erro: {e}. Reconectando em 5s...")
            await asyncio.sleep(5)

async def kafka_handler(data: dict):
    event_type = data.get("type")
    payload = data.get("payload", {})
    action = payload.get("action")

    if event_type == "message" and action in ("create", "update"):
        msg = payload.get("message", {})
        room_id = str(msg.get("chat_id", ""))
        if not room_id:
            return
        await broadcast(room_id, {
            "type": "send_message",
            "room_id": room_id,
            "id": msg.get("id"),
            "message": msg.get("message", ""),
            "user_name": msg.get("user_name", ""),
        })

    elif event_type == "message" and action == "language_update":
        msg = payload.get("message", {})
        room_id = str(msg.get("chat_id", ""))
        if not room_id:
            return
        await broadcast(room_id, {
            "type": "message_action_update",
            "room_id": room_id,
            "id": msg.get("id"),
            "language": msg.get("language", ""),
        })

    elif event_type == "message" and action == "delete":
        msg = payload.get("message", {})
        room_id = str(msg.get("chat_id", ""))
        if not room_id:
            return
        await broadcast(room_id, {
            "type": "message_deleted",
            "room_id": room_id,
            "id": msg.get("id"),
        })

    elif event_type == "chat" and action in ("create", "update"):
        chat = payload.get("chat", {})
        room_id = str(chat.get("id", ""))
        if not room_id:
            return
        await broadcast(room_id, {
            "type": "chat_updated",
            "room_id": room_id,
            "name": chat.get("name", ""),
            "category": chat.get("category", ""),
        })

    elif event_type == "chat" and action == "delete":
        chat = payload.get("chat", {})
        room_id = str(chat.get("id", ""))
        if not room_id:
            return
        await broadcast(room_id, {
            "type": "chat_deleted",
            "room_id": room_id,
            "id": chat.get("id"),
        })

async def main():
    print(f"Socket Server rodando na porta {PORT} | instance={INSTANCE_ID}")
    async with websockets.serve(handler, "0.0.0.0", PORT, process_request=process_request):
        await asyncio.gather(
            redis_subscriber(),
            consume_kafka(kafka_handler),
        )

asyncio.run(main())
