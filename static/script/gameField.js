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
    }, 500);
}

window.onload = openModalElem();

async function play() {
    const response = await fetch("/api/start");
    if (response.ok) {
        playVideo();
    }
}

//btnGo.addEventListener("click", playVideo);

