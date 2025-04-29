import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function () {
  getAccount();
  document.getElementById("logout").addEventListener("click", logout);
});

const getAccount = async () => {
  try {
    const response = await api.get('accounts/me');
    console.log(response)
  } catch (e) {
    if (e.status === 401) {
      window.location.replace('/login');
    }
  }
}

const logout = async () => {
  await api.post('accounts/logout');
  window.location.replace('/login');
}