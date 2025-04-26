const BASE_URL = '/api';

class HttpError extends Error {
  constructor(status, message, details) {
    super(message);
    this.status = status;
    this.details = details;
  }
}

class Api {
  #url;

  constructor(url) {
    this.#url = url;
  }

  apiFetch = async (endpoint, method, body = null, retry = true) => {
    if (endpoint.startsWith('/')) {
      endpoint = endpoint.slice(1);
    }

    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
    };

    if (body) {
      options.body = JSON.stringify(body);
    }

    try {
      const response = await fetch(`${this.#url}/${endpoint}`, options);

      if (!response.ok) {
        if (response.status === 401 && retry) {
          try {
            await this.apiFetch('accounts/refresh', 'POST', {}, false);
            return await this.apiFetch(endpoint, method, body, false);
          } catch (e) {
            window.location.replace('/login');
            return;
          }
        }

        let errorData = { message: 'Unknown error', details: {} };
        try {
          errorData = await response.json();
        } catch (_) { }
        throw new HttpError(response.status, errorData.message, errorData.details);
      }

      if (response.status === 204) {
        return {};
      }

      try {
        return await response.json();
      } catch (e) {
        throw new HttpError(response.status, 'Error parsing JSON', { cause: e });
      }
    } catch (error) {
      if (error instanceof HttpError) {
        this.handleHttpError(error);
      } else {
        console.error('Unexpected error:', error);
        throw error;
      }
    }
  };

  get = (endpoint, retry = true) => this.apiFetch(endpoint, 'GET', null, retry);
  post = (endpoint, body, retry = true) => this.apiFetch(endpoint, 'POST', body, retry);
  put = (endpoint, body, retry = true) => this.apiFetch(endpoint, 'PUT', body, retry);
  delete = (endpoint, body, retry = true) => this.apiFetch(endpoint, 'DELETE', body, retry);

  handleHttpError = (error) => {
    const status = error.status;

    if (status === 403) {
      alert('権限がありません。');
    } else if (status >= 500) {
      alert('予期せぬエラーが発生しました。');
    }

    throw error;
  };
}

const api = new Api(BASE_URL);

export { HttpError, Api, BASE_URL, api };
