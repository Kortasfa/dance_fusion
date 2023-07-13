function redirect() {
    if (/Android|webOS|iPhone|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)) { window.location.href = '/logIn'} else{ window.location.href = '/room';}
}
let btnGo = document.querySelector('.button');
btnGo.addEventListener("click", redirect);

document.addEventListener("DOMContentLoaded", function  ClearCookie(){{document.cookie = `User1=;expires=${new Date(0)}`;}
    {document.cookie = `User2=;expires=${new Date(0)}`;}
    {document.cookie = `User3=;expires=${new Date(0)}`;}
    {document.cookie = `User4=;expires=${new Date(0)}`;}});

