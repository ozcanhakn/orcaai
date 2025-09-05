#!/usr/bin/env python3
"""
OrcaAI Python AI Worker
Handles AI provider communication and advanced routing logic
"""

import os
import asyncio
import aiohttp
import json
import time
import hashlib
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass
from fastapi import FastAPI, HTTPException, BackgroundTasks
from pydantic import BaseModel
import uvicorn
import redis.asyncio as redis
from contextlib import asynccontextmanager

# Configuration
REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost:6379')
PORT = int(os.getenv('PYTHON_WORKER_PORT', '8001'))

# API Keys
OPENAI_API_KEY = os.getenv('OPENAI_API_KEY', '')
CLAUDE_API_KEY = os.getenv('CLAUDE_API_KEY', '')
GEMINI_API_KEY = os.getenv('GEMINI_API_KEY', '')

@dataclass
class ProviderConfig:
    name: str
    base_url: str
    headers: Dict[str, str]
    cost_per_1k_input: float
    cost_per_1k_output: float
    max_tokens: int
    rate_limit: int  # requests per minute

# Provider configurations
PROVIDERS = {
    'openai': ProviderConfig(
        name='openai',
        base_url='https://api.openai.com/v1/chat/completions',
        headers={'Authorization': f'Bearer {OPENAI_API_KEY}', 'Content-Type': 'application/json'},
        cost_per_1k_input=0.03,
        cost_per_1k_output=0.06,
        max_tokens=4000,
        rate_limit=100
    ),
    'claude': ProviderConfig(
        name='claude',
        base_url='https://api.anthropic.com/v1/messages',
        headers={'x-api-key': CLAUDE_API_KEY, 'Content-Type': 'application/json', 'anthropic-version': '2023-06-01'},
        cost_per_1k_input=0.015,
        cost_per_1k_output=0.075,
        max_tokens=100000,
        rate_limit=50
    ),
    'gemini': ProviderConfig(
        name='gemini',
        base_url=f'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key={GEMINI_API_KEY}',
        headers={'Content-Type': 'application/json'},
        cost_per_1k_input=0.001,
        cost_per_1k_output=0.002,
        max_tokens=30000,
        rate_limit=60
    )
}

# Request/Response Models
class AIRequest(BaseModel):
    prompt: str
    provider: str
    model: str
    max_tokens: Optional[int] = 1000
    temperature: Optional[float] = 0.7
    task_type: Optional[str] = "text-generation"

class AIResponse(BaseModel):
    content: str
    provider: str
    model: str
    tokens_used: Dict[str, int]
    cost: float
    latency_ms: int
    timestamp: datetime

class HealthResponse(BaseModel):
    status: str
    providers: Dict[str, str]
    uptime: float

# Global variables
redis_client: Optional[redis.Redis] = None
start_time = time.time()
request_counts = {provider: 0 for provider in PROVIDERS.keys()}
error_counts = {provider: 0 for provider in PROVIDERS.keys()}

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    global redis_client
    redis_client = redis.from_url(REDIS_URL, decode_responses=True)
    
    # Test Redis connection
    try:
        await redis_client.ping()
        print("✅ Connected to Redis")
    except Exception as e:
        print(f"❌ Redis connection failed: {e}")
    
    yield
    
    # Shutdown
    if redis_client:
        await redis_client.close()

app = FastAPI(
    title="OrcaAI Python Worker",
    description="AI Provider Worker Service",
    version="1.0.0",
    lifespan=lifespan
)

@app.get("/health", response_model=HealthResponse)
async def health_check():
    provider_status = {}
    
    for name, config in PROVIDERS.items():
        # Simple connectivity check
        try:
            async with aiohttp.ClientSession() as session:
                # For health check, just verify we can create a session
                # In production, you might want to make actual test requests
                provider_status[name] = "healthy"
        except:
            provider_status[name] = "unhealthy"
    
    return HealthResponse(
        status="healthy",
        providers=provider_status,
        uptime=time.time() - start_time
    )

@app.post("/ai/query", response_model=AIResponse)
async def process_ai_request(request: AIRequest, background_tasks: BackgroundTasks):
    start_time_request = time.time()
    
    # Validate provider
    if request.provider not in PROVIDERS:
        raise HTTPException(status_code=400, detail=f"Unsupported provider: {request.provider}")
    
    provider_config = PROVIDERS[request.provider]
    
    # Check rate limits
    if await check_rate_limit(request.provider):
        raise HTTPException(status_code=429, detail="Rate limit exceeded")
    
    try:
        # Make request to AI provider
        response_data = await make_ai_request(provider_config, request)
        
        # Calculate metrics
        latency = int((time.time() - start_time_request) * 1000)
        cost = calculate_cost(provider_config, response_data['tokens_used'])
        
        # Update success metrics
        background_tasks.add_task(update_metrics, request.provider, True, latency)
        
        return AIResponse(
            content=response_data['content'],
            provider=request.provider,
            model=request.model,
            tokens_used=response_data['tokens_used'],
            cost=cost,
            latency_ms=latency,
            timestamp=datetime.now()
        )
        
    except Exception as e:
        # Update error metrics
        background_tasks.add_task(update_metrics, request.provider, False, 0)
        raise HTTPException(status_code=500, detail=f"AI request failed: {str(e)}")

async def make_ai_request(config: ProviderConfig, request: AIRequest) -> Dict:
    """Make request to specific AI provider"""
    
    if config.name == 'openai':
        return await make_openai_request(config, request)
    elif config.name == 'claude':
        return await make_claude_request(config, request)
    elif config.name == 'gemini':
        return await make_gemini_request(config, request)
    else:
        raise ValueError(f"Unknown provider: {config.name}")

async def make_openai_request(config: ProviderConfig, request: AIRequest) -> Dict:
    """OpenAI API request"""
    payload = {
        "model": request.model,
        "messages": [{"role": "user", "content": request.prompt}],
        "max_tokens": request.max_tokens,
        "temperature": request.temperature
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.post(config.base_url, json=payload, headers=config.headers) as response:
            if response.status != 200:
                error_text = await response.text()
                raise Exception(f"OpenAI API error: {response.status} - {error_text}")
            
            data = await response.json()
            
            return {
                'content': data['choices'][0]['message']['content'],
                'tokens_used': {
                    'input': data['usage']['prompt_tokens'],
                    'output': data['usage']['completion_tokens']
                }
            }

async def make_claude_request(config: ProviderConfig, request: AIRequest) -> Dict:
    """Claude API request"""
    payload = {
        "model": request.model,
        "max_tokens": request.max_tokens,
        "messages": [{"role": "user", "content": request.prompt}],
        "temperature": request.temperature
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.post(config.base_url, json=payload, headers=config.headers) as response:
            if response.status != 200:
                error_text = await response.text()
                raise Exception(f"Claude API error: {response.status} - {error_text}")
            
            data = await response.json()
            
            # Claude API response format
            return {
                'content': data['content'][0]['text'],
                'tokens_used': {
                    'input': data['usage']['input_tokens'],
                    'output': data['usage']['output_tokens']
                }
            }

async def make_gemini_request(config: ProviderConfig, request: AIRequest) -> Dict:
    """Gemini API request"""
    payload = {
        "contents": [{"parts": [{"text": request.prompt}]}],
        "generationConfig": {
            "maxOutputTokens": request.max_tokens,
            "temperature": request.temperature
        }
    }
    
    async with aiohttp.ClientSession() as session:
        async with session.post(config.base_url, json=payload, headers=config.headers) as response:
            if response.status != 200:
                error_text = await response.text()
                raise Exception(f"Gemini API error: {response.status} - {error_text}")
            
            data = await response.json()
            
            return {
                'content': data['candidates'][0]['content']['parts'][0]['text'],
                'tokens_used': {
                    'input': data['usageMetadata']['promptTokenCount'],
                    'output': data['usageMetadata']['candidatesTokenCount']
                }
            }

def calculate_cost(config: ProviderConfig, tokens_used: Dict[str, int]) -> float:
    """Calculate cost based on token usage"""
    input_cost = (tokens_used['input'] / 1000) * config.cost_per_1k_input
    output_cost = (tokens_used['output'] / 1000) * config.cost_per_1k_output
    return round(input_cost + output_cost, 6)

async def check_rate_limit(provider: str) -> bool:
    """Check if provider rate limit is exceeded"""
    if not redis_client:
        return False
    
    key = f"rate_limit:{provider}:{int(time.time() // 60)}"
    current_count = await redis_client.get(key)
    
    if current_count is None:
        await redis_client.setex(key, 60, 1)
        return False
    
    current_count = int(current_count)
    provider_config = PROVIDERS[provider]
    
    if current_count >= provider_config.rate_limit:
        return True
    
    await redis_client.incr(key)
    return False

async def update_metrics(provider: str, success: bool, latency: int):
    """Update provider metrics in Redis"""
    if not redis_client:
        return
    
    timestamp = int(time.time())
    
    # Update counters
    if success:
        await redis_client.incr(f"metrics:{provider}:success")
        await redis_client.lpush(f"metrics:{provider}:latencies", latency)
        await redis_client.ltrim(f"metrics:{provider}:latencies", 0, 99)  # Keep last 100 latencies
    else:
        await redis_client.incr(f"metrics:{provider}:errors")
    
    # Update hourly metrics
    hour_key = f"metrics:{provider}:hourly:{timestamp // 3600}"
    await redis_client.incr(hour_key)
    await redis_client.expire(hour_key, 7 * 24 * 3600)  # Keep for 7 days

@app.get("/metrics")
async def get_metrics():
    """Get provider performance metrics"""
    metrics = {}
    
    for provider in PROVIDERS.keys():
        success_count = await redis_client.get(f"metrics:{provider}:success") or 0
        error_count = await redis_client.get(f"metrics:{provider}:errors") or 0
        
        # Get average latency
        latencies = await redis_client.lrange(f"metrics:{provider}:latencies", 0, -1)
        avg_latency = sum(map(int, latencies)) / len(latencies) if latencies else 0
        
        # Calculate reliability
        total_requests = int(success_count) + int(error_count)
        reliability = int(success_count) / total_requests if total_requests > 0 else 1.0
        
        metrics[provider] = {
            "success_count": int(success_count),
            "error_count": int(error_count),
            "avg_latency_ms": round(avg_latency, 2),
            "reliability": round(reliability, 4),
            "total_requests": total_requests
        }
    
    return {"metrics": metrics, "timestamp": datetime.now()}

@app.get("/providers")
async def get_providers():
    """Get available providers and their configurations"""
    providers_info = {}
    
    for name, config in PROVIDERS.items():
        providers_info[name] = {
            "name": config.name,
            "cost_per_1k_input": config.cost_per_1k_input,
            "cost_per_1k_output": config.cost_per_1k_output,
            "max_tokens": config.max_tokens,
            "rate_limit": config.rate_limit,
            "status": "active" if globals().get(f"{name.upper()}_API_KEY") else "inactive"
        }
    
    return {"providers": providers_info}

@app.post("/routing/smart")
async def smart_routing(request: dict):
    """Advanced routing logic based on historical performance"""
    prompt = request.get('prompt', '')
    task_type = request.get('task_type', 'text-generation')
    user_preferences = request.get('preferences', {})
    
    # Get current metrics for all providers
    provider_scores = {}
    
    for provider_name in PROVIDERS.keys():
        # Get recent performance data
        success_count = int(await redis_client.get(f"metrics:{provider_name}:success") or 0)
        error_count = int(await redis_client.get(f"metrics:{provider_name}:errors") or 0)
        
        # Get average latency
        latencies = await redis_client.lrange(f"metrics:{provider_name}:latencies", 0, 9)  # Last 10 requests
        avg_latency = sum(map(int, latencies)) / len(latencies) if latencies else 2000
        
        # Calculate reliability
        total_requests = success_count + error_count
        reliability = success_count / total_requests if total_requests > 0 else 0.9
        
        # Get provider config
        config = PROVIDERS[provider_name]
        
        # Calculate composite score
        score = calculate_provider_score(
            cost=config.cost_per_1k_input,
            latency=avg_latency,
            reliability=reliability,
            task_type=task_type,
            user_preferences=user_preferences
        )
        
        provider_scores[provider_name] = {
            'score': score,
            'reliability': reliability,
            'avg_latency': avg_latency,
            'cost': config.cost_per_1k_input
        }
    
    # Sort by score and return best provider
    best_provider = max(provider_scores.items(), key=lambda x: x[1]['score'])
    
    # Get fallback providers
    sorted_providers = sorted(provider_scores.items(), key=lambda x: x[1]['score'], reverse=True)
    fallbacks = [p[0] for p in sorted_providers[1:3]]  # Top 2 alternatives
    
    return {
        "recommended_provider": best_provider[0],
        "confidence": min(best_provider[1]['score'], 1.0),
        "reasoning": generate_routing_reason(best_provider, task_type),
        "fallbacks": fallbacks,
        "all_scores": provider_scores
    }

def calculate_provider_score(cost: float, latency: float, reliability: float, 
                           task_type: str, user_preferences: dict) -> float:
    """Calculate composite score for provider selection"""
    
    # Normalize metrics (0-1 scale)
    cost_score = max(0, 1 - (cost / 0.1))  # Normalize to $0.1 per 1k tokens
    latency_score = max(0, 1 - (latency / 5000))  # Normalize to 5 seconds
    reliability_score = reliability
    
    # Task-specific quality scores
    quality_scores = {
        'openai': {
            'text-generation': 0.9,
            'code-generation': 0.95,
            'summarization': 0.85,
            'conversation': 0.9
        },
        'claude': {
            'text-generation': 0.95,
            'code-generation': 0.9,
            'summarization': 0.95,
            'reasoning': 0.98
        },
        'gemini': {
            'text-generation': 0.8,
            'multimodal': 0.95,
            'summarization': 0.75
        }
    }
    
    # Default weights
    weights = {
        'cost': user_preferences.get('cost_weight', 0.25),
        'latency': user_preferences.get('latency_weight', 0.25),
        'reliability': user_preferences.get('reliability_weight', 0.3),
        'quality': user_preferences.get('quality_weight', 0.2)
    }
    
    # Calculate weighted score
    quality_score = 0.8  # Default quality score
    
    final_score = (
        cost_score * weights['cost'] +
        latency_score * weights['latency'] +
        reliability_score * weights['reliability'] +
        quality_score * weights['quality']
    )
    
    return final_score

def generate_routing_reason(best_provider_data, task_type: str) -> str:
    """Generate human-readable routing explanation"""
    provider_name, metrics = best_provider_data
    
    reasons = []
    
    if metrics['reliability'] > 0.95:
        reasons.append("high reliability")
    if metrics['avg_latency'] < 1500:
        reasons.append("fast response")
    if metrics['cost'] < 0.01:
        reasons.append("cost-effective")
    
    if not reasons:
        reasons.append("balanced performance")
    
    return f"Selected {provider_name} for {task_type}: {', '.join(reasons)}"

if __name__ == "__main__":
    print("🚀 Starting OrcaAI Python Worker...")
    print(f"📡 Listening on port {PORT}")
    
    # Check API keys
    active_providers = []
    for name, key in [('OpenAI', OPENAI_API_KEY), ('Claude', CLAUDE_API_KEY), ('Gemini', GEMINI_API_KEY)]:
        if key:
            active_providers.append(name)
            print(f"✅ {name} API key configured")
        else:
            print(f"⚠️  {name} API key not found")
    
    print(f"🔧 Active providers: {', '.join(active_providers)}")
    
    uvicorn.run(
        "worker:app",
        host="0.0.0.0",
        port=PORT,
        reload=True,
        log_level="info"
    )