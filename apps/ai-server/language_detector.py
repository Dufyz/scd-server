"""
Language detection using OpenAI API
"""

import os
import logging
from openai import OpenAI

logger = logging.getLogger(__name__)

_client: OpenAI | None = None


def _get_client() -> OpenAI:
    global _client
    if _client is None:
        _client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
    return _client


def detect_language(text: str) -> str:
    """
    Detect the language of the given text using OpenAI API

    Args:
        text: Text to analyze

    Returns:
        ISO 639-1 language code (e.g., 'en', 'pt', 'es', 'fr')
    """
    try:
        model = os.getenv("OPENAI_MODEL", "gpt-4o-mini")

        response = _get_client().chat.completions.create(
            model=model,
            messages=[
                {
                    "role": "system",
                    "content": (
                        "You are a language detection system. "
                        "Respond ONLY with the ISO 639-1 two-letter language code "
                        "(e.g., 'en' for English, 'pt' for Portuguese, 'es' for Spanish). "
                        "Nothing else."
                    )
                },
                {
                    "role": "user",
                    "content": f"Detect the language of this text: {text}"
                }
            ],
            temperature=0,
            max_tokens=10,
        )

        detected_lang = response.choices[0].message.content.strip().lower()
        logger.info(f"Language detected: {detected_lang}")

        return detected_lang

    except Exception as e:
        logger.error(f"Error calling OpenAI API: {e}", exc_info=True)
        return "unknown"
