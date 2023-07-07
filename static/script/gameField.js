const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnGo = document.getElementById("btn-go")
let isBtnClicked = false

function openModalElem() {
    modalElem.classList.add("open");
}

function closeModalElem() {
    isBtnClicked = true;
    console.log(isBtnClicked);
    modalElem.classList.remove("open");
}

function emulateClick(btn) {
    let click = new CustomEvent("mousemove");
    btn.dispatchEvent(click);
    console.log("click!")
}

function playVideo() {
    modalElem.classList.remove("open");
    setTimeout(() => {
        danceVideo.play();
    }, 500);
}

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

btnGo.addEventListener("click", playVideo);