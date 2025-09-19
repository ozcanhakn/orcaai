import { Activity, Database, DollarSign, Settings } from 'lucide-react';
import { useEffect, useState } from 'react';
import Layout from '../components/Layout';

export default function Providers() {
  const [providers, setProviders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState(null);

  useEffect(() => {
    fetchProviders();
  }, []);

  const fetchProviders = async () => {
    try {
      const response = await fetch('/api/v1/providers');
      const data = await response.json();
      setProviders(data);
    } catch (error) {
      console.error('Error fetching providers:', error);
    } finally {
      setLoading(false);
    }
  };

  const updateProvider = async (id, config) => {
    try {
      await fetch(`/api/v1/providers/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(config),
      });
      fetchProviders();
      setShowConfigModal(false);
    } catch (error) {
      console.error('Error updating provider:', error);
    }
  };

  return (
    <Layout>
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">AI Providers</h1>
          <p className="text-gray-400">Configure and monitor your AI providers.</p>
        </div>

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
            <p className="mt-4 text-gray-400">Loading providers...</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {providers.map((provider) => (
              <div key={provider.id} className="bg-gray-800 rounded-lg overflow-hidden">
                <div className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center">
                      <Database className="w-6 h-6 text-blue-400 mr-2" />
                      <h2 className="text-xl font-bold">{provider.name}</h2>
                    </div>
                    <button
                      onClick={() => {
                        setSelectedProvider(provider);
                        setShowConfigModal(true);
                      }}
                      className="text-gray-400 hover:text-white"
                    >
                      <Settings className="w-5 h-5" />
                    </button>
                  </div>

                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Status</span>
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        provider.healthy ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'
                      }`}>
                        {provider.healthy ? 'Healthy' : 'Issues Detected'}
                      </span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Success Rate</span>
                      <span className="text-white">{(provider.success_rate * 100).toFixed(1)}%</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Avg. Latency</span>
                      <span className="text-white">{provider.avg_latency}ms</span>
                    </div>

                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Cost/1K tokens</span>
                      <span className="text-white">${provider.cost_per_1k.toFixed(4)}</span>
                    </div>
                  </div>

                  <div className="mt-6 pt-6 border-t border-gray-700">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-gray-400">24h Usage</span>
                      <Activity className="w-4 h-4 text-blue-400" />
                    </div>
                    <div className="h-12 bg-gray-900 rounded-lg overflow-hidden">
                      {/* Add mini usage chart here */}
                    </div>
                  </div>
                </div>

                <div className="bg-gray-900 px-6 py-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center">
                      <DollarSign className="w-4 h-4 text-green-400 mr-1" />
                      <span className="text-gray-400">Today's Cost</span>
                    </div>
                    <span className="text-white font-bold">${provider.today_cost.toFixed(2)}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Provider Configuration Modal */}
        {showConfigModal && selectedProvider && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
            <div className="bg-gray-800 rounded-lg p-6 w-96">
              <h2 className="text-xl font-bold mb-4">Configure {selectedProvider.name}</h2>
              <form onSubmit={(e) => {
                e.preventDefault();
                updateProvider(selectedProvider.id, {
                  api_key: e.target.api_key.value,
                  max_tokens: parseInt(e.target.max_tokens.value),
                  priority: parseInt(e.target.priority.value),
                });
              }}>
                <div className="space-y-4">
                  <div>
                    <label className="block text-gray-400 mb-2">API Key</label>
                    <input
                      type="password"
                      name="api_key"
                      className="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2"
                      defaultValue={selectedProvider.api_key}
                    />
                  </div>
                  <div>
                    <label className="block text-gray-400 mb-2">Max Tokens</label>
                    <input
                      type="number"
                      name="max_tokens"
                      className="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2"
                      defaultValue={selectedProvider.max_tokens}
                    />
                  </div>
                  <div>
                    <label className="block text-gray-400 mb-2">Priority (1-5)</label>
                    <input
                      type="number"
                      name="priority"
                      min="1"
                      max="5"
                      className="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2"
                      defaultValue={selectedProvider.priority}
                    />
                  </div>
                </div>
                <div className="flex justify-end gap-2 mt-6">
                  <button
                    type="button"
                    onClick={() => setShowConfigModal(false)}
                    className="px-4 py-2 text-gray-400 hover:text-white"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-700"
                  >
                    Save Changes
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}
