"""
Kafka consumer for message.created events
"""

import os
import json
import logging
from confluent_kafka import Consumer, Producer, KafkaError, KafkaException

from language_detector import detect_language

logger = logging.getLogger(__name__)

KAFKA_CONSUMER_GROUP_ID = os.getenv('KAFKA_CONSUMER_GROUP_ID', 'ai-server')
KAFKA_TOPIC_CONSUME = os.getenv('KAFKA_TOPIC_CONSUME', 'message')
KAFKA_TOPIC_PRODUCE = os.getenv('KAFKA_TOPIC_PRODUCE', 'message.language_detected')


def create_consumer():
    """Create and configure Kafka consumer"""
    brokers = os.getenv('KAFKA_BROKERS', 'localhost:9094')
    logger.info(f"Conectando ao Kafka em: {brokers}")
    config = {
        'bootstrap.servers': brokers,
        'group.id': KAFKA_CONSUMER_GROUP_ID,
        'auto.offset.reset': 'earliest',
        'enable.auto.commit': True,
    }
    return Consumer(config)


def create_producer():
    """Create and configure Kafka producer"""
    config = {
        'bootstrap.servers': os.getenv('KAFKA_BROKERS', 'localhost:9094'),
        'acks': 'all',
        'retries': 3,
    }
    return Producer(config)


def delivery_callback(err, msg):
    """Callback for producer delivery reports"""
    if err:
        logger.error(f"Failed to deliver message to {msg.topic()}: {err}")
    else:
        logger.info(f"Message delivered to {msg.topic()} [{msg.partition()}] at offset {msg.offset()}")


def process_message(producer, msg):
    """
    Process a single message from Kafka

    Expected message format (from 'message' topic):
    {
        "type": "message",
        "payload": {
            "action": "create",
            "message": {
                "id": 123,
                "chat_id": 456,
                "message": "Hello world",
                "user_name": "John",
                "created_at": "2024-06-14T10:00:00Z",
                "language": null
            }
        }
    }
    """
    try:
        # Parse message
        event_data = json.loads(msg.value().decode('utf-8'))

        # Filter only 'create' actions
        payload = event_data.get('payload', {})
        action = payload.get('action', '')

        if action != 'create':
            logger.debug(f"Skipping event with action: {action}")
            return

        # Extract message data from payload
        message_data = payload.get('message', {})
        message_id = message_data.get('id')

        logger.info(f"Processing message ID: {message_id}")

        # Extract text to analyze
        text = message_data.get('message', '')
        if not text:
            logger.warning(f"Empty message text for ID: {message_id}")
            return

        # Detect language using OpenAI
        detected_language = detect_language(text)

        # Build result payload
        result = {
            'message_id': message_id,
            'chat_id': message_data.get('chat_id'),
            'detected_language': detected_language,
            'original_message': text,
        }

        # Publish to output topic
        producer.produce(
            topic=KAFKA_TOPIC_PRODUCE,
            key=str(message_id).encode('utf-8'),
            value=json.dumps(result).encode('utf-8'),
            callback=delivery_callback
        )
        producer.poll(0)  # Trigger delivery callbacks

        logger.info(f"Language detected for message {message_id}: {detected_language}")

    except json.JSONDecodeError as e:
        logger.error(f"Failed to parse message JSON: {e}")
    except Exception as e:
        logger.error(f"Error processing message: {e}", exc_info=True)


def start_consumer(running_flag):
    """
    Start Kafka consumer loop

    Args:
        running_flag: Global flag to control consumer loop
    """
    consumer = create_consumer()
    producer = create_producer()

    consumer.subscribe([KAFKA_TOPIC_CONSUME])
    logger.info(f"Subscribed to topic: {KAFKA_TOPIC_CONSUME} (group: {KAFKA_CONSUMER_GROUP_ID})")

    try:
        while not running_flag.is_set():
            msg = consumer.poll(timeout=1.0)

            if msg is None:
                continue

            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event - not an error
                    logger.debug(f"Reached end of partition {msg.partition()}")
                elif msg.error().code() == KafkaError.UNKNOWN_TOPIC_OR_PART:
                    # Topic doesn't exist yet - log warning and continue
                    logger.warning(f"Topic not available yet: {msg.error()}. Waiting...")
                    continue
                else:
                    raise KafkaException(msg.error())
            else:
                process_message(producer, msg)

    except Exception as e:
        logger.error(f"Consumer error: {e}", exc_info=True)

    finally:
        logger.info("Closing consumer and producer...")
        producer.flush(timeout=10)
        consumer.close()
