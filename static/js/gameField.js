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

function playVideo() {
    modalElem.classList.remove("open");
    setTimeout(() => {
        danceVideo.play();
    }, 2000);
}

window.onload = openModalElem();
btnGo.addEventListener("click", playVideo);
