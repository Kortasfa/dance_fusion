function redirect() {
    window.location.href = '/room';
}
let btnGo = document.querySelector('.button');
btnGo.addEventListener("click", redirect);
