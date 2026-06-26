import json
from room_manager import add_client, broadcast
from redis_client import get_redis, INSTANCE_ID

async def handle_event(websocket, raw_message):
    try:
        data = json.loads(raw_message)
    except json.JSONDecodeError:
        print(f"Mensagem inválida recebida: {raw_message}")
        return

    event_type = data.get("type")
    if not event_type:
        print("Mensagem sem campo 'type', ignorando")
        return

    if event_type == "join_room":
        await handle_join_room(websocket, data)
    elif event_type == "send_message":
        await handle_send_message(websocket, data)
    elif event_type == "typing":
        await handle_typing(websocket, data)
    else:
        print(f"Tipo de evento desconhecido: {event_type}")

async def handle_join_room(websocket, data):
    room_id = str(data.get("room_id"))
    if not room_id:
        print("join_room sem room_id, ignorando")
        return
    add_client(room_id, websocket)
    print(f"Cliente entrou na sala {room_id}")

async def handle_send_message(websocket, data):
    room_id = str(data.get("room_id"))
    message = data.get("message")
    user_name = data.get("user_name")
    if not room_id or not message or not user_name:
        print("send_message com campos faltando, ignorando")
        return

    print(f"Mensagem de {user_name} na sala {room_id}: {message}")

    event = {
        "type": "send_message",
        "room_id": room_id,
        "message": message,
        "user_name": user_name,
    }

    await broadcast(room_id, event)

    try:
        redis = get_redis()
        await redis.publish(f"room:{room_id}", json.dumps({
            "instance_id": INSTANCE_ID,
            "room_id": room_id,
            "event": event,
        }))
    except Exception as e:
        print(f"Redis publish erro (mensagem já entregue localmente): {e}")

async def handle_typing(websocket, data):
    room_id = str(data.get("room_id"))
    user_name = data.get("user_name")
    typing = data.get("typing", False)
    if not room_id or not user_name:
        print("typing com campos faltando, ignorando")
        return
    print(f"{user_name} {'está digitando' if typing else 'parou de digitar'} na sala {room_id}")
    await broadcast(room_id, {
        "type": "typing",
        "room_id": room_id,
        "user_name": user_name,
        "typing": typing,
    }, exclude=websocket)
