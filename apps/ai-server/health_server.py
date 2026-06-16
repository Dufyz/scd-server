"""
Health check HTTP server
"""

import logging
from fastapi import FastAPI
import uvicorn

logger = logging.getLogger(__name__)

app = FastAPI(title="AI Server - Language Detection")


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "message": "Health status endpoint is operational",
        "status": "up"
    }


def start_health_server(port: int):
    """
    Start FastAPI health check server

    Args:
        port: Port to listen on
    """
    logger.info(f"Starting health server on port {port}")
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=port,
        log_level="info",
        access_log=False
    )
