const listSong = document.getElementById('list-song');
const listGenre = document.getElementById('list-genre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('return-button');
const PlayBtn = document.getElementById('Play');

let readyGame = false;
let readyPlayer = false;
let readySong = false;

let connectedUsers = [];

btnOpenInfo.addEventListener('click', openGuide);

function changeButton() {
    if (readyPlayer && readySong) {
        readyGame = true;
    }
    if (readyGame) {
        PlayBtn.classList.add('button_ready');
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

let songName = '';
let test = document.getElementsByClassName("section__img");
let testing = '';
// Iterate over each element in the collection
Array.from(test).forEach(function (element) {
    element.addEventListener('click', function () {
        let videoSrcID = '9' + element.id;
        let video = document.getElementById(videoSrcID);
        let fullVideo = document.getElementById('full' + videoSrcID);
        let videoPlayer = document.getElementById('videoPlayer');
        songName = document.querySelector('.song' + element.id).innerHTML;
        // let menuItem = parent.querySelectorAll('.button_yellow');
        // // Отлавливаем элемент в родители на который мы нажали
        // for(let i = 0; i < menuItem.length; i++) {
        //   menuItem[i].classList.remove('button_yellow');
        // }
        readySong = true;
        changeButton();
        videoPlayer.src = video.innerText;
        testing = fullVideo.innerText;
    });
});

$(document).ready(function () {
    let trigger = $('#Play');
    let container = $('#content');

    // Fire on click
    trigger.on('click', function() {
        if (readyGame) {

            socket.send(songName);
            socket.close();
            console.log(songName)// Можно отправить pause или resume


            container.load("/static/html/game.html", function () {
                let video = $('#video-dance').get(0);
                let src = $('#video-src').get(0);
                src.setAttribute('src', testing);

                video.addEventListener('loadeddata', function () {
                    video.play();
                });
            });
            return false;
        }
    });
});


function showVideo(videoID) {
    let videoSrcID = '9' + videoID.id;
    let video = document.getElementById(videoSrcID);
    let videoPlayer = document.getElementById('videoPlayer');
    videoPlayer.src = video.innerText;
}

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

let socket = new WebSocket(WssURL);

socket.onopen = function (event) {
    console.log("WebSocket connection established.");
};

socket.onmessage = function (event) {
    let message = event.data;
    let parts = message.split('|');
    let action = parts[0];
    let userID = parts[1];
    if (action === "add") {
        let userName = parts[2];
        let imgSrc = parts[3];
        addUser(userID, userName, imgSrc);
    } else if (action === "remove") {
        removeUser(userID)
    }

};

socket.onclose = function (event) {
    console.log("WebSocket connection closed.");
};

function addUser(userID, userName, imgSrc) {
    console.log('Пользователь присоединился: ' + userID);
    connectedUsers.push({"userID": userID, "userName": userName, "imgSrc": imgSrc});

    let userMessage = document.getElementById('needUser');
    userMessage.classList.add('none');
    let indexUser = document.getElementById('user' + connectedUsers.length);
    let indexUserName = document.getElementById('userName' + connectedUsers.length);
    let indexUserImg = indexUser.querySelector(".user__avatar");
    indexUser.classList.remove('none');
    indexUserImg.src = '../' + imgSrc;
    indexUserName.innerText = userName;

    readyPlayer = true;
    changeButton();
}

function removeUser(userID) {
    console.log('Пользователь вышел: ' + userID);

    let removedUserIndex = 0;
    for (let i = 0; i < connectedUsers.length; i++) {
        if (connectedUsers[i]["userID"] === userID) {
            connectedUsers.splice(i, 1);
            removedUserIndex = i;
            break;
        }
    }

    if (connectedUsers.length === 0) {
        let userMessage = document.getElementById('needUser');
        userMessage.classList.remove('none');
        readyPlayer = false;
        changeButton();
    }

    let indexUser = document.getElementById('user' + (removedUserIndex + 1));
    let indexUserName = document.getElementById('userName' + (removedUserIndex + 1));
    indexUser.classList.add('none');
    indexUser.id = 'user4';
    indexUserName.id = 'userName4';
    const parent = indexUser.parentElement;
    parent.removeChild(indexUser);
    parent.appendChild(indexUser);

    for (let i = removedUserIndex + 2; i <= 4; i++) {
        let indexUser = document.getElementById('user' + i);
        let indexUserName = document.getElementById('userName' + i);
        indexUser.id = 'user' + (i - 1);
        indexUserName.id = 'userName' + (i - 1);
    }
}
