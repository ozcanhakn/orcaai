/**
 * OrcaAI JavaScript SDK
 * 
 * A JavaScript SDK for the OrcaAI platform.
 * 
 * Basic usage:
 * 
 * const { OrcaClient } = require('@orcaai/sdk');
 * 
 * const client = new OrcaClient({ apiKey: 'your-api-key' });
 * const result = await client.query('Summarize this report');
 * console.log(result);
 */

const axios = require('axios');

class OrcaAIError extends Error {
  constructor(message) {
    super(message);
    this.name = 'OrcaAIError';
  }
}

class AuthenticationError extends OrcaAIError {
  constructor(message) {
    super(message);
    this.name = 'AuthenticationError';
  }
}

class APIError extends OrcaAIError {
  constructor(message) {
    super(message);
    this.name = 'APIError';
  }
}

class OrcaClient {
  /**
   * Initialize the OrcaClient
   * @param {Object} options - Client options
   * @param {string} options.apiKey - Your OrcaAI API key
   * @param {string} [options.baseUrl='http://localhost:8080'] - Base URL for the OrcaAI API
   */
  constructor(options = {}) {
    if (!options.apiKey) {
      throw new Error('API key is required');
    }

    this.apiKey = options.apiKey;
    this.baseUrl = options.baseUrl || 'http://localhost:8080';
    this.axios = axios.create({
      baseURL: this.baseUrl,
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
        'Content-Type': 'application/json'
      }
    });
  }

  /**
   * Send a query to the OrcaAI platform
   * @param {string} prompt - The prompt to send to the AI
   * @param {Object} [options] - Query options
   * @param {string} [options.taskType='text-generation'] - The type of task
   * @param {string} [options.provider] - Specific provider to use
   * @param {string} [options.model] - Specific model to use
   * @returns {Promise<Object>} The AI response
   */
  async query(prompt, options = {}) {
    const payload = {
      prompt,
      task_type: options.taskType || 'text-generation'
    };

    if (options.provider) {
      payload.provider = options.provider;
    }

    if (options.model) {
      payload.model = options.model;
    }

    try {
      const response = await this.axios.post('/api/v1/ai/query', payload);
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Get available AI providers
   * @returns {Promise<Object>} List of available providers
   */
  async getProviders() {
    try {
      const response = await this.axios.get('/api/v1/ai/providers');
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Get usage metrics
   * @returns {Promise<Object>} Usage metrics
   */
  async getMetrics() {
    try {
      const response = await this.axios.get('/api/v1/metrics');
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Get all API keys for the authenticated user
   * @returns {Promise<Object>} List of API keys
   */
  async getApiKeys() {
    try {
      const response = await this.axios.get('/api/v1/keys');
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Create a new API key
   * @param {string} name - Name for the new API key
   * @returns {Promise<Object>} The created API key
   */
  async createApiKey(name) {
    try {
      const response = await this.axios.post('/api/v1/keys', { name });
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Delete an API key
   * @param {string} keyId - ID of the API key to delete
   * @returns {Promise<Object>} Confirmation of deletion
   */
  async deleteApiKey(keyId) {
    try {
      const response = await this.axios.delete(`/api/v1/keys/${keyId}`);
      return response.data;
    } catch (error) {
      this._handleError(error);
    }
  }

  /**
   * Handle API errors
   * @param {Object} error - The error object
   * @private
   */
  _handleError(error) {
    if (error.response) {
      const { status, data } = error.response;
      
      if (status === 401) {
        throw new AuthenticationError('Invalid API key');
      } else {
        throw new APIError(`API request failed: ${data.error || data.message || 'Unknown error'}`);
      }
    } else if (error.request) {
      throw new OrcaAIError('No response received from server');
    } else {
      throw new OrcaAIError(`Request failed: ${error.message}`);
    }
  }
}

module.exports = {
  OrcaClient,
  OrcaAIError,
  AuthenticationError,
  APIError
};