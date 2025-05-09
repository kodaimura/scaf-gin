import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function () {
  document.getElementById("login").addEventListener("click", login);
});


const login = async () => {
  const form = document.getElementById("login-form");
  const name = form.elements['name'].value;
  const password = form.elements['password'].value;

  const body = {
    name: name,
    password: password
  };

  try {
    await api.post('accounts/login', body, false);
    window.location.replace('/');
  } catch (e) {
    document.getElementById("error").innerHTML = (e.status === 401)
      ? "ユーザ名またはパスワードが異なります。"
      : "ログインに失敗しました。";
  }
}