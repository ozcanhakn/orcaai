import { Copy, Key, Plus, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import { apiFetch } from '../components/apiClient';

export default function APIKeys() {
  const [keys, setKeys] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showNewKeyModal, setShowNewKeyModal] = useState(false);

  useEffect(() => {
    fetchAPIKeys();
  }, []);

  const fetchAPIKeys = async () => {
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
      const data = await apiFetch('/keys', { token });
      setKeys(data.api_keys || []);
    } catch (error) {
      console.error('Error fetching API keys:', error);
    } finally {
      setLoading(false);
    }
  };

  const createNewKey = async (name) => {
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
      const data = await apiFetch('/keys', { method: 'POST', body: { name }, token });
      setKeys([data.api_key, ...keys]);
      setShowNewKeyModal(false);
    } catch (error) {
      console.error('Error creating API key:', error);
    }
  };

  const deleteKey = async (id) => {
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
      await apiFetch(`/keys/${id}`, { method: 'DELETE', token });
      setKeys(keys.filter(key => key.id !== id));
    } catch (error) {
      console.error('Error deleting API key:', error);
    }
  };

  const rotateKey = async (id) => {
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : undefined;
      const data = await apiFetch(`/keys/${id}/rotate`, { method: 'POST', token });
      if (data?.api_key) setKeys([data.api_key, ...keys.filter(k => k.id !== id)]);
    } catch (error) {
      console.error('Error rotating API key:', error);
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    // Show toast notification
  };

  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold mb-2">API Keys</h1>
            <p className="text-gray-400">Manage your API keys for accessing the OrcaAI API.</p>
          </div>
          <button
            onClick={() => setShowNewKeyModal(true)}
            className="px-4 py-2 bg-blue-600 rounded-lg hover:bg-blue-700 flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            New API Key
          </button>
        </div>

        {loading ? (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
            <p className="mt-4 text-gray-400">Loading API keys...</p>
          </div>
        ) : (
          <div className="bg-gray-800 rounded-lg overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="bg-gray-900">
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Name</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Key</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Created</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Last Used</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-400 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-700">
                {keys.map((key) => (
                  <tr key={key.id}>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <Key className="w-4 h-4 text-gray-400 mr-2" />
                        {key.name}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <code className="bg-gray-900 px-2 py-1 rounded">
                          {key.key ? `${key.key.slice(0,4)}••••••••••••${key.key.slice(-4)}` : `${key.prefix || ''}...${key.suffix || ''}`}
                        </code>
                        <button
                          onClick={() => copyToClipboard(key.key)}
                          className="ml-2 text-gray-400 hover:text-white"
                        >
                          <Copy className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-400">
                      {new Date(key.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-gray-400">
                      {key.last_used ? new Date(key.last_used).toLocaleDateString() : 'Never'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right">
                      <button
                        onClick={() => deleteKey(key.id)}
                        className="text-red-400 hover:text-red-300"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => rotateKey(key.id)}
                        className="ml-3 text-blue-400 hover:text-blue-300"
                      >
                        Rotate
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* New Key Modal */}
        {showNewKeyModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
            <div className="bg-gray-800 rounded-lg p-6 w-96">
              <h2 className="text-xl font-bold mb-4">Create New API Key</h2>
              <form onSubmit={(e) => {
                e.preventDefault();
                createNewKey(e.target.name.value);
              }}>
                <div className="mb-4">
                  <label className="block text-gray-400 mb-2">Key Name</label>
                  <input
                    type="text"
                    name="name"
                    className="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2"
                    placeholder="e.g., Production Key"
                  />
                </div>
                <div className="flex justify-end gap-2">
                  <button
                    type="button"
                    onClick={() => setShowNewKeyModal(false)}
                    className="px-4 py-2 text-gray-400 hover:text-white"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-700"
                  >
                    Create Key
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
