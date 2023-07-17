const listSong = document.getElementById('list-song');
const listGenre = document.getElementById('list-genre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('return-button');
const PlayBtn = document.getElementById('Play');

let readyGame = false;
let numberOfUser = 0;
let readyPlayer = false;
let readySong = false;

btnOpenInfo.addEventListener('click', openGuide);

document.addEventListener("DOMContentLoaded", getUsersByCookie);

function getUsersByCookie() {
    let allCookies = document.cookie;
    let cookiesArray = allCookies.split(';');
    for (let i = 0; i < cookiesArray.length; i++) {
        let name = cookiesArray[i].split('=');
        let mane = name[0];
        let n = mane.indexOf('User');
        if (n === 1) {
            let parts = name[1].split(',');
            let userID = parts[0];
            let userName = parts[1];
            let imgSrc = parts[2];
            numberOfUser = numberOfUser + 1;
            let userMSG = document.getElementById('needUser');
            userMSG.classList.add('none');
            let indexUser = document.getElementById('user' + numberOfUser);
            let indexUserName = document.getElementById('userName' + numberOfUser);
            let indexUserImg = indexUser.querySelector(".user__avatar");
            indexUser.classList.remove('none');
            indexUserImg.src = '../' + imgSrc;
            indexUserName.innerText = userName;
            readyPlayer = true;
                changeButton();
        }
    }
}

function changeButton() {
    if (readyPlayer && readySong) {
        readyGame = true;
    }
    if (readyGame) {
        PlayBtn.classList.add('button_yellow');
    }
}

function openSong(styleButtonBlock) {
    returnBtn.classList.toggle('hide')
    listGenre.classList.add('none');
    document.querySelector('.songs').classList.remove('none');
    const songBlock = document.getElementsByClassName('song__section');
    for (let i = 0; i < songBlock.length; i++) {
        if ((songBlock[i].id === styleButtonBlock.id) && songBlock[i].classList.contains('none')) {
            songBlock[i].classList.remove('none');
        }
    }
}

function closeSong() {
    returnBtn.classList.toggle('hide')
    const songBlock = document.getElementsByClassName('song__section');
    for (let i = 0; i < songBlock.length; i++) {
        if (!songBlock[i].classList.contains('none')) {
            songBlock[i].classList.add('none');
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

let test = document.getElementsByClassName("section__img");
let testing = '';
// Iterate over each element in the collection
Array.from(test).forEach(function (element) {
    element.addEventListener('click', function () {
        let videoSrcID = '9' + element.id;
        let video = document.getElementById(videoSrcID);
        let videoPlayer = document.getElementById('videoPlayer');

        // let menuItem = parent.querySelectorAll('.button_yellow');
        // // Отлавливаем элемент в родители на который мы нажали
        // for(let i = 0; i < menuItem.length; i++) {
        //   menuItem[i].classList.remove('button_yellow');
        // }
        readySong = true;
        changeButton();
        videoPlayer.src = video.innerText;
        testing = video.innerText;
    });
});

$(document).ready(function () {
    let trigger = $('#Play');
    let container = $('#content');

    // Fire on click
    trigger.on('click', function () {
        if (PlayBtn.classList.contains('button_yellow')) {
            container.load("/static/html/game.html", function () {
                let video = $('#video-dance').get(0);
                let src = $('#video-src').get(0);
                src.setAttribute('src', testing);

                video.addEventListener('loadeddata', function () {
                    video.play();
                });

                video.addEventListener('ended', function () {
                    showStats();
                });
            });

            return false;
        }
    });
});

function showStats() {
    document.getElementById('video-dance').style.display = "none";
}

function showVideo(videoID) {
    let videoSrcID = '9' + videoID.id;
    let video = document.getElementById(videoSrcID);
    let videoPlayer = document.getElementById('videoPlayer');
    videoPlayer.src = video.innerText;
}

let socket = new WebSocket(WssURL);

socket.onopen = function (event) {
    console.log("WebSocket connection established.");
};

socket.onmessage = function (event) {
    let userMessage = document.getElementById('needUser');
    let message = event.data;
    let parts = message.split('|');
    let userID = parts[0];
    let userName = parts[1];
    let imgSrc = parts[2];
    numberOfUser = numberOfUser + 1;
    document.cookie = "User" + numberOfUser + '=' + parts + ';path=/';
    userMessage.classList.add('none');
    let indexUser = document.getElementById('user' + numberOfUser);
    let indexUserName = document.getElementById('userName' + numberOfUser);
    let indexUserImg = indexUser.querySelector(".user__avatar");
    indexUser.classList.remove('none');
    indexUserImg.src = '../' + imgSrc;
    readyPlayer = true;
    changeButton();
    indexUserName.innerText = userName;
    console.log('Пользователь присоединился: ' + userID);
    console.log('Его имя: ' + userName);
    console.log('Его фотка: ' + imgSrc);
};

socket.onclose = function (event) {
    console.log("WebSocket connection closed.");
};

let parent = document.querySelector('.songs');

function addColor(song) {
    let menuItem = parent.querySelectorAll('.button_yellow');
    for (let i = 0; i < menuItem.length; i++) {
        // Убираем у других
        menuItem[i].classList.remove('button_yellow');
    }
    setTimeout(changeColor(song), 1000);
}

function changeColor(song) {
    song.classList.add('button_yellow');
}