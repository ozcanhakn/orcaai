const { OrcaAI } = require('orcaai');

// İstemciyi oluştur
const client = new OrcaAI('your-api-key');

// Basit bir sorgu
client.query('JavaScript nedir?')
  .then(response => {
    console.log(response.text);
  })
  .catch(error => {
    console.error('Hata:', error);
  });

// Gelişmiş sorgu
async function advancedQuery() {
  try {
    const response = await client.query('React ile web uygulaması nasıl geliştirilir?', {
      taskType: 'tutorial',
      options: {
        costWeight: 0.7,
        latencyWeight: 0.3,
        maxBudget: 0.05,
        preferredProviders: ['openai', 'anthropic']
      }
    });

    console.log(`Yanıt: ${response.text}`);
    console.log(`Sağlayıcı: ${response.provider}`);
    console.log(`Maliyet: $${response.cost}`);
    console.log(`Gecikme: ${response.latencyMs}ms`);
  } catch (error) {
    console.error('Hata:', error);
  }
}

advancedQuery();
