const listSong = document.getElementById('listSong');
const listGenre = document.getElementById('listGenre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('returnButton');
const PlayBtn = document.getElementById('play');

let readyGame = false;
let numberOfUser = 0;
let readyPlayer = false;

btnOpenInfo.addEventListener('click', openGuide);

function changeButton() {
    readyGame = readyPlayer && readySong;
    if (readyGame) {
        PlayBtn.classList.add('button_ready');
    }
}

function openSong(styleButtonBlock) {
    returnBtn.classList.toggle('hide');
    listGenre.classList.add('none');
    document.querySelector('.songs').classList.remove('none');
    const songBlocks = document.getElementsByClassName('song__section');
    for (let i = 0; i < songBlocks.length; i++) {
        if (songBlocks[i].id === styleButtonBlock.id && songBlocks[i].classList.contains('none')) {
            songBlocks[i].classList.remove('none');
        }
    }
}

returnBtn.addEventListener('click', closeSong);

function closeSong() {
    returnBtn.classList.toggle('hide');
    const songBlocks = document.getElementsByClassName('song__section');
    for (let i = 0; i < songBlocks.length; i++) {
        if (!songBlocks[i].classList.contains('none')) {
            songBlocks[i].classList.add('none');
        }
    }
    listSong.classList.add('none');
    listGenre.classList.remove('none');
}

function openGuide() {
    guide.classList.add('play');
    guide.classList.remove('unPlay');
}

function closeGuide() {
    guide.classList.add('unPlay');
    guide.classList.remove('play');
}

let songName = '';
let fullSongName = '';
let readySong = false;

function onImageClick(element) {
    const videoSrcID = 'song' + element.id;
    const video = document.getElementById(videoSrcID);
    const fullVideo = document.getElementById('full' + videoSrcID);
    const videoPlayer = document.getElementById('videoPlayer');

    songName = document.querySelector('.song' + element.id).innerHTML;
    readySong = true;
    changeButton();

    videoPlayer.src = video.innerText;
    fullSongName = fullVideo.innerText;
}

const test = document.getElementsByClassName('section__img');

Array.from(test).forEach(function (element) {
    element.addEventListener('click', function () {
        onImageClick(element);
    });
});

$(document).ready(function () {
    const playButton = $('#play');
    const contentContainer = $('#content');

    playButton.on('click', function () {
        if (readyGame) {
            socket.send(songName);

            contentContainer.load("/static/html/game.html", function () {
                const video = $('#video-dance')[0];
                const src = $('#video-src')[0];
                src.setAttribute('src', fullSongName);

                video.addEventListener('loadeddata', function () {
                    video.play();
                });
            });

            return false;
        }
    });
});


function showVideo(videoID) {
    let videoSrcID = 'song' + videoID.id;
    let video = document.getElementById(videoSrcID);
    let videoPlayer = document.getElementById('videoPlayer');
    videoPlayer.src = video.innerText;
}

const parent = document.querySelector('.songs');

function changeColor(song) {
    song.classList.add('button_yellow');
}

function addColor(song) {
    const menuItem = parent.querySelectorAll('.button_yellow');
    for (let i = 0; i < menuItem.length; i++) {
        menuItem[i].classList.remove('button_yellow');
    }
    changeColor(song);
}

let socket = new WebSocket(WssURL);

socket.onopen = function (event) {
    console.log("WebSocket connection established.");
};

socket.onmessage = function (event) {
    const userMessage = document.getElementById('needUser');
    const message = event.data;
    const parts = message.split('|');
    const userID = parts[0];
    const userName = parts[1];
    const imgSrc = parts[2];

    numberOfUser = numberOfUser + 1;
    document.cookie = `User${numberOfUser}=${parts};path=/`;

    userMessage.classList.add('none');
    const indexUser = document.getElementById(`user${numberOfUser}`);
    const indexUserName = document.getElementById(`userName${numberOfUser}`);
    const indexUserImg = indexUser.querySelector(".user__avatar");
    indexUser.classList.remove('none');
    indexUserImg.src = `../${imgSrc}`;

    readyPlayer = true;
    changeButton();

    indexUserName.innerText = userName;

    console.log(`Пользователь присоединился: ${userID}`);
    console.log(`Его имя: ${userName}`);
    console.log(`Его фотка: ${imgSrc}`);
};


socket.onclose = function (event) {
    console.log("WebSocket connection closed.");
};
