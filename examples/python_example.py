from orcaai import Client

# İstemciyi oluştur
client = Client(api_key="your-api-key")

# Basit bir sorgu
response = client.query("Python programlama dili nedir?")
print(response.text)

# Gelişmiş sorgu (maliyet optimizasyonlu)
response = client.query(
    prompt="Python ile web scraping nasıl yapılır?",
    task_type="explanation",
    options={
        "cost_weight": 0.8,        # Maliyet optimizasyonuna öncelik ver
        "latency_weight": 0.2,     # Hıza daha az öncelik ver
        "max_budget": 0.05,        # Maksimum 0.05$ harca
        "preferred_providers": ["openai", "anthropic"]  # Tercih edilen sağlayıcılar
    }
)

# Başarılı yanıt
if response.success:
    print(f"Yanıt: {response.text}")
    print(f"Kullanılan sağlayıcı: {response.provider}")
    print(f"Maliyet: ${response.cost:.4f}")
    print(f"Gecikme süresi: {response.latency_ms}ms")
else:
    print(f"Hata: {response.error}")
