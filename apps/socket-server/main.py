import asyncio
import os
import websockets
import sys
sys.stdout.reconfigure(line_buffering=True)
from dotenv import load_dotenv
from room_manager import remove_client
from event_handler import handle_event

load_dotenv()

PORT = int(os.getenv("PORT", 8765))

async def handler(websocket):
    print(f"Cliente conectado: {websocket.remote_address}")
    try:
        async for message in websocket:
            await handle_event(websocket, message)
    finally:
        remove_client(websocket)
        print(f"Cliente desconectado: {websocket.remote_address}")

async def main():
    print(f"Socket Server rodando na porta {PORT}")
    async with websockets.serve(handler, "0.0.0.0", PORT):
        await asyncio.Future()

asyncio.run(main())