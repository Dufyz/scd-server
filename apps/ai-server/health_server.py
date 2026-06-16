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
        "status": "healthy",
        "service": "ai-server-language-detection"
    }


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": "AI Server - Language Detection",
        "version": "1.0.0"
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
