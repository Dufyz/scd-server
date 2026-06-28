import asyncio
import json
import os
from confluent_kafka import Consumer, KafkaError

KAFKA_BROKERS = os.getenv("KAFKA_BROKERS", "localhost:9094")

def _create_consumer() -> Consumer:
    print(f"Conectando ao Kafka em: {KAFKA_BROKERS}")
    consumer = Consumer({
        "bootstrap.servers": KAFKA_BROKERS,
        "group.id": "socket-server",
        "auto.offset.reset": "latest",
        "session.timeout.ms": 10000,
    })
    consumer.subscribe(["message", "chat"])
    return consumer

async def consume_kafka(callback) -> None:
    while True:
        consumer = None
        try:
            consumer = await asyncio.to_thread(_create_consumer)
            print("Kafka consumer iniciado nos tópicos: message, chat")
            while True:
                msg = await asyncio.to_thread(consumer.poll, 1.0)
                if msg is None:
                    continue
                if msg.error():
                    if msg.error().code() != KafkaError._PARTITION_EOF:
                        print(f"Kafka error: {msg.error()}")
                    continue
                try:
                    data = json.loads(msg.value().decode("utf-8"))
                    await callback(data)
                except Exception as e:
                    print(f"Erro processando mensagem Kafka: {e}")
        except Exception as e:
            print(f"Kafka consumer erro: {e}. Reconectando em 5s...")
            await asyncio.sleep(5)
        finally:
            if consumer:
                consumer.close()
