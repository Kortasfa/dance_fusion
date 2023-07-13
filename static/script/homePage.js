function redirect() {
    window.location.href = '/room';
}
let btnGo = document.querySelector('.button');
btnGo.addEventListener("click", redirect);

document.addEventListener("DOMContentLoaded", function  ClearCookie(){{document.cookie = `User1=;expires=${new Date(0)}`;}
    {document.cookie = `User2=;expires=${new Date(0)}`;}
    {document.cookie = `User3=;expires=${new Date(0)}`;}
    {document.cookie = `User4=;expires=${new Date(0)}`;}});

