class HttpError extends Error {
  constructor(status, message, details) {
    super(message);
    this.status = status;
    this.details = details;
  }
}

class Api {
  constructor(url) {
    this.url = url;
  }

  async createFetchOptions(method, body) {
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    };

    if (body) {
      options.body = JSON.stringify(body);
    }
    return options;
  }

  async apiFetch(endpoint, method, body, retry = true) {
    if (endpoint.startsWith('/')) {
      endpoint = endpoint.slice(1);
    }

    const options = await this.createFetchOptions(method, body);
    const response = await fetch(`${this.url}/${endpoint}`, options);

    if (!response.ok) {
      if (response.status === 401 && retry && window.location.pathname !== '/login') {
        const refreshed = await this.tryRefreshToken();
        if (refreshed) {
          return this.apiFetch(endpoint, method, body, false);
        }
      }

      let errorData = { error: 'Unknown error', details: {} };
      try {
        errorData = await response.json();
      } catch {
        // ignore parse error
      }

      const error = new HttpError(response.status, errorData.error, errorData.details);
      this.handleHttpError(error);
      throw error;
    }

    if (response.status === 204) {
      return {};
    }

    return await response.json();
  }

  async tryRefreshToken() {
    try {
      const response = await fetch(`${this.url}/accounts/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      return response.ok;
    } catch {
      return false;
    }
  }

  async get(endpoint) {
    return this.apiFetch(endpoint, 'GET');
  }

  async post(endpoint, body) {
    return this.apiFetch(endpoint, 'POST', body);
  }

  async put(endpoint, body) {
    return this.apiFetch(endpoint, 'PUT', body);
  }

  async delete(endpoint, body) {
    return this.apiFetch(endpoint, 'DELETE', body);
  }

  async patch(endpoint, body) {
    return this.apiFetch(endpoint, 'PATCH', body);
  }

  handleHttpError(error) {
    console.error(error);
    const status = error.status;
    if (status === 403) {
      alert('アクセスが拒否されました');
    } else if (status === 401 && window.location.pathname !== '/login') {
      window.location.replace('/login');
    }
  }
}

const api = new Api('/api');

export { Api, api, HttpError };
