import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function () {
  document.getElementById("signup").addEventListener("click", signup);
});


const signup = async () => {
  const form = document.getElementById("signup-form");
  if (!validate(form)) return;

  const name = form.elements['name'].value;
  const password = form.elements['password'].value;

  const body = {
    name: name,
    password: password
  };

  try {
    await api.post('accounts/signup', body);
    window.location.replace('/login');
  } catch (e) {
    document.getElementById("error").innerHTML = (e.status === 409)
      ? "ユーザ名が既に使われています。"
      : "登録に失敗しました。";
  }
}

const validate = (form) => {
  const name = form.elements['name'].value;
  const password = form.elements['password'].value;
  const password_confirm = form.elements['password_confirm'].value;

  let error = "";
  if (name === "") {
    error = "ユーザ名を入力して下さい。";
  } else if (password === "") {
    error = "パスワードを入力して下さい。";
  } else if (password.length < 8) {
    error = "パスワードは8文字以上で入力してください。";
  } else if (password !== password_confirm) {
    error = "パスワードが一致していません。";
  }

  document.getElementById("error").innerHTML = error;
  return error === "";
}