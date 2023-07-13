// const danceVideo = document.getElementById("video-dance");
// const modalElem = document.getElementById("pop-up");
// const btnGo = document.getElementById("btn-go");
let isBtnClicked = false;

function getUsersByCookie() {
    console.log('lol')
    let numberOfUser = 0;
    let allCookies = document.cookie;
    let cookiesArray = allCookies.split(';');
    for (let i = 0; i < cookiesArray.length; i++) {
        let cookie = cookiesArray[i].split('=')
        let name = cookie[0];
        let findUser = name.indexOf('User');
        if (findUser === 1) {
            let parts = cookie[1].split(',');
            let userName = parts[1];
            let imgSrc = parts[2];
            numberOfUser = numberOfUser + 1;
            let indexUser = document.getElementById('hero' + numberOfUser);
            let indexUserName = document.getElementById('heroName' + numberOfUser);
            let indexUserImg = document.getElementById('heroImg' + numberOfUser);
            indexUser.classList.remove('hidden');
            indexUserImg.src =  '../' + imgSrc;
            indexUserName.innerText = userName;
            redyPlayer = true;
            ChengeBtn();
        }
    }
}

getUsersByCookie();
// function openModalElem() {
//     modalElem.classList.add("open");
// }

// x

// function emulateClick(btn) {
//     let click = new CustomEvent("mousemove");
//     btn.dispatchEvent(click);
//     console.log("click!")
// }

// function playVideo() {
//     modalElem.classList.remove("open");
//     setTimeout(() => {
//         danceVideo.play();
//     }, 500);
// }

window.onload = openModalElem();
/*window.onload = emulateClick(btnGo);

/*function test() {
    console.log("test click")
}

btnGo.addEventListener("mousemove", test());*/

/*let socket = new WebSocket(WssURL);

socket.onopen = function(event) {
    console.log("WebSocket connection established.");
}

socket.onmessage = function() {
    //let message = event.data;
    /*x.removeAttribute('disabled');
    btnGo.click();*/
/*  playVideo();

}

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
}*/

// btnGo.addEventListener("click", playVideo);