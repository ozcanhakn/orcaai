import { Send, Terminal } from 'lucide-react';
import { useState } from 'react';
import Layout from '../../components/Layout';

export default function Playground() {
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [history, setHistory] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    setLoading(true);
    try {
      const token = localStorage.getItem('token'); // Auth token'ı localStorage'dan al
      const response = await fetch('http://localhost:8080/api/v1/ai/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`, // Auth token'ı ekle
        },
        body: JSON.stringify({
          prompt: input,
          options: {
            cost_weight: 0.7,
            latency_weight: 0.3
          }
        })
      });

      const data = await response.json();
      setHistory(prev => [{
        prompt: input,
        response: data,
        timestamp: new Date(),
      }, ...prev]);
      
      setInput('');
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="p-6 max-w-4xl mx-auto">
        <div className="flex items-center gap-3 mb-6">
          <Terminal className="w-8 h-8 text-blue-500" />
          <h1 className="text-2xl font-bold">AI Playground</h1>
        </div>

        <form onSubmit={handleSubmit} className="mb-8">
          <div className="relative">
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Enter your prompt here..."
              className="w-full h-32 bg-gray-800 border border-gray-700 rounded-lg p-4 pr-12 text-white resize-none focus:outline-none focus:border-blue-500"
            />
            <button
              type="submit"
              disabled={loading || !input.trim()}
              className="absolute bottom-4 right-4 p-2 bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Send className="w-5 h-5" />
            </button>
          </div>
        </form>

        <div className="space-y-6">
          {history.map((item, index) => (
            <div key={index} className="bg-gray-800 rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-400">
                    {item.timestamp.toLocaleTimeString()}
                  </span>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-sm text-gray-400">
                    Cost: ${item.response.cost?.toFixed(4) || '0.00'}
                  </span>
                  <span className="text-sm text-gray-400">
                    Latency: {item.response.latency_ms || 0}ms
                  </span>
                  <span className="text-sm text-gray-400">
                    Provider: {item.response.provider || 'Unknown'}
                  </span>
                </div>
              </div>
              
              <div className="mb-4">
                <div className="text-sm text-gray-400 mb-2">Prompt:</div>
                <div className="bg-gray-900 rounded-lg p-4">
                  {item.prompt}
                </div>
              </div>

              <div>
                <div className="text-sm text-gray-400 mb-2">Response:</div>
                <div className="bg-gray-900 rounded-lg p-4">
                  {item.response.content || 'No response'}
                </div>
              </div>
            </div>
          ))}

          {history.length === 0 && !loading && (
            <div className="text-center py-12 text-gray-400">
              No queries yet. Try sending your first prompt!
            </div>
          )}

          {loading && (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
              <p className="mt-4 text-gray-400">Processing your request...</p>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
}
