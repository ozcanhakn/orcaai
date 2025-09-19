import { useState } from 'react';
import Layout from '../components/Layout';
import { apiFetch } from '../components/apiClient';

export default function AIPlayground() {
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [history, setHistory] = useState([]);

  const [provider, setProvider] = useState('');
  const [model, setModel] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
      const data = await apiFetch('/ai/query', {
        method: 'POST',
        token,
        body: {
          prompt: input,
          provider: provider || undefined,
          model: model || undefined,
          options: { cost_weight: 0.7, latency_weight: 0.3 },
        },
      });
      setResult(data);
      setHistory([data, ...history]);
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">AI Playground</h1>
          <p className="text-gray-400">Test your AI queries and see how the orchestrator works.</p>
        </div>

        <form onSubmit={handleSubmit} className="mb-8">
          <div className="flex gap-4">
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              className="flex-1 bg-gray-800 border border-gray-700 rounded-lg p-4 text-white"
              placeholder="Enter your query here..."
              rows={4}
            />
            <div className="flex flex-col gap-2">
              <input
                value={provider}
                onChange={(e) => setProvider(e.target.value)}
                className="bg-gray-800 border border-gray-700 rounded px-3 py-2"
                placeholder="Provider (optional)"
              />
              <input
                value={model}
                onChange={(e) => setModel(e.target.value)}
                className="bg-gray-800 border border-gray-700 rounded px-3 py-2"
                placeholder="Model (optional)"
              />
              <button
                type="submit"
                disabled={loading}
                className="px-6 py-3 bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? 'Processing...' : 'Send Query'}
              </button>
            </div>
          </div>
        </form>

        {result && (
          <div className="bg-gray-800 rounded-lg p-6 mb-8">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-bold">Result</h2>
              <div className="flex items-center gap-4">
                <span className="text-sm text-gray-400">
                  Provider: {result.provider}
                </span>
                <span className="text-sm text-gray-400">
                  Cost: ${result.cost.toFixed(4)}
                </span>
                <span className="text-sm text-gray-400">
                  Latency: {result.latency_ms}ms
                </span>
              </div>
            </div>
            <pre className="bg-gray-900 p-4 rounded-lg overflow-x-auto">
              {result.content}
            </pre>
          </div>
        )}

        {history.length > 0 && (
          <div>
            <h2 className="text-xl font-bold mb-4">Query History</h2>
            <div className="space-y-4">
              {history.map((item, index) => (
                <div key={index} className="bg-gray-800 rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm text-gray-400">
                      Provider: {item.provider}
                    </span>
                    <div className="flex items-center gap-4">
                      <span className="text-sm text-gray-400">
                        Cost: ${item.cost.toFixed(4)}
                      </span>
                      <span className="text-sm text-gray-400">
                        Latency: {item.latency_ms}ms
                      </span>
                    </div>
                  </div>
                  <pre className="bg-gray-900 p-4 rounded-lg overflow-x-auto text-sm">
                    {item.content}
                  </pre>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}
