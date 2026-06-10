import json

rooms: dict[str, set] = {}

def get_or_create_room(room_id: str) -> set:
    if room_id not in rooms:
        rooms[room_id] = set()
    return rooms[room_id]

def add_client(room_id: str, websocket) -> None:
    room = get_or_create_room(room_id)
    room.add(websocket)
    print(f"Cliente adicionado na sala {room_id}. Total: {len(room)}")

def remove_client(websocket) -> None:
    empty_rooms = []
    for room_id, clients in rooms.items():
        clients.discard(websocket)
        if len(clients) == 0:
            empty_rooms.append(room_id)
    for room_id in empty_rooms:
        del rooms[room_id]
        print(f"Sala {room_id} removida por estar vazia")

def get_clients(room_id: str) -> set:
    return rooms.get(room_id, set()).copy()

async def broadcast(room_id: str, event: dict, exclude=None) -> None:
    clients = get_clients(room_id)
    if exclude:
        clients.discard(exclude)
    message = json.dumps(event)
    disconnected = []
    for client in clients:
        try:
            await client.send(message)
        except Exception:
            disconnected.append(client)
    for client in disconnected:
        remove_client(client)