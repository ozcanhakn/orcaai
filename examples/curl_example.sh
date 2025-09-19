#!/bin/bash

# API anahtarınız
API_KEY="your-api-key"

# Basit sorgu örneği
curl -X POST http://localhost:8080/api/v1/ai/query \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Yapay zeka nedir?",
    "task_type": "explanation"
  }'

# Gelişmiş sorgu örneği
curl -X POST http://localhost:8080/api/v1/ai/query \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Python ile veri analizi nasıl yapılır?",
    "task_type": "tutorial",
    "options": {
      "cost_weight": 0.7,
      "latency_weight": 0.3,
      "max_budget": 0.05,
      "preferred_providers": ["openai", "anthropic"]
    }
  }'

# Metrik görüntüleme
curl -X GET http://localhost:8080/api/v1/metrics \
  -H "Authorization: Bearer $API_KEY"
