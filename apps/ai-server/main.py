"""
AI Server - Language Detection Service
Consumes messages from Kafka, detects language using OpenAI API,
and publishes results back to Kafka.
"""

import os
import sys
import json
import signal
import logging
from threading import Thread, Event
from dotenv import load_dotenv

from kafka_consumer import start_consumer
from health_server import start_health_server

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger(__name__)

# Load environment variables
load_dotenv()

stop_event = Event()


def signal_handler(signum, frame):
    logger.info(f"Received signal {signum}, shutting down gracefully...")
    stop_event.set()


def main():
    """Main entry point"""
    logger.info("Starting AI Server - Language Detection Service")

    # Validate required environment variables
    required_env_vars = [
        "KAFKA_BROKERS",
        "OPENAI_API_KEY",
    ]

    missing_vars = [var for var in required_env_vars if not os.getenv(var)]
    if missing_vars:
        logger.error(f"Missing required environment variables: {', '.join(missing_vars)}")
        sys.exit(1)

    # Register signal handlers for graceful shutdown
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    # Start health check server in background thread
    port = int(os.getenv("PORT", "8070"))
    health_thread = Thread(target=start_health_server, args=(port,), daemon=True)
    health_thread.start()
    logger.info(f"Health check server started on port {port}")

    # Start Kafka consumer (blocking call)
    try:
        start_consumer(stop_event)
    except KeyboardInterrupt:
        logger.info("Interrupted by user")
    except Exception as e:
        logger.error(f"Fatal error in consumer: {e}", exc_info=True)
        sys.exit(1)

    logger.info("AI Server shutdown complete")


if __name__ == "__main__":
    main()
