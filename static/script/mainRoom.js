const listSong = document.getElementById('listSong');
const listGenre = document.getElementById('listGenre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('returnButton');
const PlayBtn = document.getElementById('play');

let readyGame = false;
let numberOfUser = 0;
let readyPlayer = false;
let connectedUsers = [];

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
            socket.close(); // Закрываем вебсокет mainRoom

            function loadScript(url) {
                return new Promise(function (resolve, reject) {
                    const script = document.createElement('script');
                    script.src = url;
                    script.onload = resolve;
                    script.onerror = reject;
                    document.head.appendChild(script);
                });
            }

            async function loadEdgeImpulseAndRunImpulse() {
                try {
                    await loadScript('/static/test/edge-impulse-standalone.js');
                    await loadScript('/static/test/run-impulse.js');
                    console.log('Both scripts loaded and executed successfully.');
                    // You can now call any functions or perform actions from the loaded scripts.
                    // For example: Module.onRuntimeInitialized();
                } catch (error) {
                    console.error('Error loading scripts:', error);
                }
            }

            loadEdgeImpulseAndRunImpulse().then(r => {
                contentContainer.load("/static/html/game.html", function () {
                    // This function will be executed after the new content is loaded.
                    const video = $('#video-dance')[0];
                    const src = $('#video-src')[0];
                    src.setAttribute('src', fullSongName);

                    video.addEventListener('loadeddata', function () {
                        video.play();
                    });

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
    console.log("WebSocket mainRoom connection established.");
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
    console.log("WebSocket mainRoom connection closed.");
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