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

  apiFetch = async (endpoint, method, body = null) => {
    if (endpoint.startsWith('/')) {
      endpoint = endpoint.slice(1);
    }

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

    try {
      const response = await fetch(`${this.#url}/${endpoint}`, options);

      if (!response.ok) {
        let errorData = { message: 'Unknown error', details: {} };
        try {
          errorData = await response.json();
        } catch (e) {
          // レスポンスがJSONでない場合（HTMLなど）、そのまま fallback メッセージ
        }
        throw new HttpError(response.status, errorData.message, errorData.details);
      }

      if (response.status === 204) {
        return {}; // No Content
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

  get = (endpoint) => this.apiFetch(endpoint, 'GET');
  post = (endpoint, body) => this.apiFetch(endpoint, 'POST', body);
  put = (endpoint, body) => this.apiFetch(endpoint, 'PUT', body);
  delete = (endpoint, body) => this.apiFetch(endpoint, 'DELETE', body);

  handleHttpError = (error) => {
    console.error('HTTP Error:', error);
    throw error;
  };
}

const api = new Api(BASE_URL);

api.handleHttpError = (error) => {
  const status = error.status;

  if (status === 401 && window.location.pathname !== '/login') {
    window.location.replace('/login');
  } else if (status === 403) {
    alert("権限がありません。");
  } else if (status >= 500) {
    alert("予期せぬエラーが発生しました。");
  }

  throw error;
};

export { HttpError, Api, BASE_URL, api };
