const listSong = document.getElementById('listSong');
const listGenre = document.getElementById('listGenre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('returnButton');
const playBtn = document.getElementById('play');
const songs = document.querySelector('.songs');
const gameMode = document.getElementById('gameMode');
const bots = document.getElementById('bots');
const botsMenu = document.getElementById('botMenu');
const boss = document.getElementById('boss');

let readyGame = false;
let numberOfUser = 0;
let readyPlayer = false;
let mode = ''
let connectedUsers = [];
let connectedBots = [];
let startedGame = false;

btnOpenInfo.addEventListener('click', openGuide);
window.onbeforeunload = sendGameEndInfoToServer;

const audio = document.querySelector("audio");
audio.volume = 0.5;
let notClicked = 1;


window.addEventListener("click", event => {
    if (notClicked){
        audio.play();
        notClicked = 0;
    }
});

function changeButton() {
    readyGame = readyPlayer && readySong;
    if (readyGame) {
        playBtn.classList.add('button_ready');
    }
}

function openSong(styleButtonBlock) {
    listGenre.classList.add('none');
    songs.classList.remove('none');
    const songBlocks = document.getElementsByClassName('song__section');
    for (let i = 0; i < songBlocks.length; i++) {
        if (songBlocks[i].id === styleButtonBlock.id && songBlocks[i].classList.contains('none')) {
            songBlocks[i].classList.remove('none');
        }
    }
}

function openStyles(regime) {
    listGenre.classList.remove('none');
    returnBtn.classList.toggle('hide');
    gameMode.classList.add('none');
    mode = regime.getAttribute('mode');
    if (mode == 'Boss'){
        boss.classList.remove('none');
        playBtn.classList.add('none');
    } else if (mode == 'Bots'){
        bots.classList.remove('none');
        setTimeout(function() {
            bots.classList.remove('none');
        }, 800);
    }
}

returnBtn.addEventListener('click', closeList);

function toggleBots(){
    if (botsMenu.classList.contains('bots_open')){
        botsMenu.classList.add('bots_close');
        botsMenu.classList.remove('bots_open');
    }else{
        botsMenu.classList.add('bots_open');
        botsMenu.classList.remove('bots_close');
    }
}
function closeList() {
    if (listSong.classList.contains('none')){
        audio.play();
        returnBtn.classList.toggle('hide');
        listGenre.classList.add('none');
        gameMode.classList.remove('none');
        if(mode === 'Bots'){
            botsMenu.classList.add('bots_close');
            botsMenu.classList.add('bots_open');
            setTimeout(function() {
                bots.classList.add('none');
                botsMenu.classList.remove('bots_close');
            }, 1001);
        }
        if (mode === 'Boss') {
            boss.classList.add('none');
            playBtn.classList.remove('none');
        }
        for (let i = 0; i < connectedUsers.length; i++) {
            if (parseInt(connectedUsers[i]["userID"]) < 0 ) {
                removeUser(connectedUsers[i]["userID"]);
                if (parseInt(connectedUsers[i]["userID"]) < 0 ) {
                    removeUser(connectedUsers[i]["userID"]);
                }
                if (parseInt(connectedUsers[i]["userID"]) < 0 ) {
                    removeUser(connectedUsers[i]["userID"]);
                }
                if (parseInt(connectedUsers[i]["userID"]) < 0 ) {
                    removeUser(connectedUsers[i]["userID"]);
                }
            }
        }
    } else {
        const songBlocks = document.getElementsByClassName('song__section');
        for (let i = 0; i < songBlocks.length; i++) {
            if (!songBlocks[i].classList.contains('none')) {
                songBlocks[i].classList.add('none');
            }
        }
        listSong.classList.add('none');
        listGenre.classList.remove('none');
    }
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
let songNeuro = '';
let fullSongName = '';
let songId = 0;
let difficulty = 0;
let readySong = false;
const videoPlayer = document.getElementById('videoPlayer');

function onImageClick(element) {
    const videoSrcID = 'song' + element.id;
    const video = document.getElementById(videoSrcID);
    const fullVideo = document.getElementById('full' + videoSrcID);

    audio.pause();
    let difficultyList = document.querySelector('.game-menu__difficulty');
    let difficultySegment = difficultyList.querySelectorAll('.segment')

    difficulty  =  document.getElementById('difficulty' + videoSrcID).innerText;
    difficultyList.classList.remove('none')
    for (let i = 0; i < 4; i++) {
        if (i < difficulty) {
            difficultySegment[i].classList.add('segment_on')
        } else {
            difficultySegment[i].classList.remove('segment_on')
        }
    }

    songId = element.id;

    songName = document.querySelector('.song' + element.id).innerHTML;

    songNeuro = camelCase(songName);

    function camelCase(value) {
        return value.toLowerCase().replace(/\s+(.)/g, function(match, group1) {
            return group1.toUpperCase();
        });
    }
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

var numb;
var classifier;

$(document).ready(function() {
    const playButton = $('#play');
});
let playButton = document.getElementById('play');
playButton.addEventListener('click', gameStart);
function gameStart() {
    const contentContainer = $('#content');
    if (readyGame) {
        startedGame = true;
        sendGameStartInfoToServer().then(() => {})
        numb = (Math.round(Math.random()*1000)).toString();
        let firstComponent = '/static/test/edge-impulse-standalone.js?version=' + numb;
        let secondComponent = '/static/test/run-impulse.js?version=' + numb;
        let thirdComponent = '/static/html/game.html?version=' + numb;
        $.getScript(firstComponent, function() {
            $.getScript(secondComponent, function() {
                (async () => {
                    document.getElementsByClassName('loading')[0].style.display= 'flex';
                    videoPlayer.src = '';
                    classifier = new EdgeImpulseClassifier();
                    await classifier.init();
                    let project = classifier.getProjectInfo();
                    console.log(project.owner + ' / ' + project.name + ' (version ' + project.deploy_version + ')');

                    contentContainer.load(thirdComponent, function() {
                        const video = $('#video-dance')[0];
                        const src = $('#video-src')[0];
                        src.setAttribute('src', fullSongName);

                        video.addEventListener('loadeddata', function() {
                            video.play();
                            socket.send(songName);
                            console.log(songName);
                            socket.close(); // Закрываем вебсокет mainRoom
                        });
                    });
                })();
            });
        });
        fetch('../static/script/scoreGrade.js')
            .then(response => response.text())
            .then(scriptText => {
                // Создаем элемент <script>
                const script = document.createElement('script');
                script.textContent = scriptText;
                document.head.appendChild(script);
            })
            .catch(error => {
                console.error('Ошибка загрузки скрипта:', error);
            });
    }
    return readyGame;
}

function showVideo(videoID) {
    let videoSrcID = 'song' + videoID.id;
    let video = document.getElementById(videoSrcID);
    let videoPlayer = document.getElementById('videoPlayer');
    videoPlayer.src = video.innerText;
}

function readJSONFromURL(url) {
    return fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка при получении файла');
            }
            return response.json();
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}

function addBot(botName) {
    let botInfo = {};
    fetch("../api/getBotPath", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `bot_name=${botName}`,
    })
        .then(response => {
            if (!response.ok) {
                if (response.status === 409) {
                    throw new Error('No bot with such name');
                }
                else {
                    throw new Error('Server Error');
                }
            }
            return response.json();
        })
        .then(data => {
            if (data.BotScoresPath) {
                botInfo = data;
                if (!addUser( "-" + data.BotId,  "" + botName, "../" + data.BotImgHat, "../" + data.BotImgFace, "../" + data.BotImgBody)) {
                    return;
                }
                readJSONFromURL("../" + data.BotScoresPath).then(jsonData => {
                    connectedBots.push({"botID": "-" + data.BotId,  "botScores":  jsonData});
                    console.log(connectedBots);
                });
                fetch("/api/addBot", { // Добавление бота на бэке
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        "room_id": window.location.pathname.split('/').pop(),
                        "bot_id": "-" + botInfo.BotId
                    })
                })
                    .then(function (response) {
                        if (response.ok) {
                            console.log('Бот добавлен на бэке');
                        } else {
                            console.log('Ошибка при добавлении бота на бэке. Статус:', response.status);
                        }
                    })
                    .catch(function (error) {
                        console.log('Ошибка при отправке данных: ', error);
                    });
            } else {
                console.log('Fail');
            }
        })
        .catch(error => {
            console.log('Error:', error);
            return
        });
}

let bossInfo;

function bossGame(bossBlock) {
    let name = bossBlock.querySelector(".boss__name").innerText;
    let healthPoint = bossBlock.querySelector(".boss__health-point").innerText.replace('Hp: ', '');
    let bossBody= bossBlock.querySelector(".body").src;
    let bossFace = bossBlock.querySelector(".face").src;
    let bossHat = bossBlock.querySelector(".hat").src;
    let bossId = bossBlock.id;
    if (!gameStart()) {
        console.log("Игра не готова")
        return;
    }
    bossInfo = {"bossId": bossId, "name": name, "healthPoint": healthPoint, "bossBody": bossBody, "bossFace": bossFace, "bossHat": bossHat};

}


function changeColor(song) {
    song.classList.add('button_yellow');
}

function addColor(song) {
    const menuItem = songs.querySelectorAll('.button_yellow');
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
        let hatImgSrc = "../" + parts[3];
        let faceImgSrc = "../" + parts[4];
        let bodyImgSrc = "../" + parts[5];
        addUser(userID, userName, hatImgSrc, faceImgSrc, bodyImgSrc);
    } else if (action === "remove") {
        removeUser(userID)
    }

};

socket.onclose = function (event) {
    if (!startedGame) {
        window.location.href = "/room";
    }
    console.log("WebSocket mainRoom connection closed.");
};

function addUser(userID, userName, hatImgSrc, faceImgSrc, bodyImgSrc) {
    for (let userInfo of connectedUsers) {
        if (userInfo["userID"] === userID) {
            if(userID > 0){
                console.log('Пользователь уже присоединён: ' + userID);
            } else {
                removeUser(userID);
                console.log('Бот удалён: ' + userID);
            }
            return false;
        }
    }
    if (connectedUsers.length >= 4) {
        console.log('Пользователь не может присоединиться, так как комната переполнена: ' + userID);
        return false;
    }
    console.log('Пользователь присоединился: ' + userID);
    connectedUsers.push({"userID": userID, "userName": userName, "valueScore": 0, "bodyImgSrc": bodyImgSrc, "faceImgSrc": faceImgSrc, "hatImgSrc": hatImgSrc});

    let userMessage = document.getElementById('needUser');
    userMessage.classList.add('none');
    let indexUser = document.getElementById('user' + connectedUsers.length);
    let indexUserName = document.getElementById('userName' + connectedUsers.length);
    let indexUserBodyImg = indexUser.querySelector(".body");
    let indexUserFaceImg = indexUser.querySelector(".face");
    let indexUserHatImg = indexUser.querySelector(".hat");
    indexUser.classList.remove('none');
    indexUserBodyImg.src = bodyImgSrc;
    indexUserFaceImg.src = faceImgSrc;
    indexUserHatImg.src = hatImgSrc;
    indexUserName.innerText = userName;

    readyPlayer = true;
    changeButton();
    return true;
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
    if (parseInt(userID) < 0) {
        fetch("/api/removeBot", { // Если это бот, то удаляем бота на бэке
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                "room_id": window.location.pathname.split('/').pop(),
                "bot_id": userID
            })
        })
            .then(function (response) {
                if (response.ok) {
                    console.log('Бот удалён на бэке');
                } else {
                    console.log('Ошибка при удалении бота на бэке. Статус:', response.status);
                }
            })
            .catch(function (error) {
                console.log('Ошибка при отправке данных: ', error);
            });
    }
}


const domain = window.location.protocol + "//" + window.location.hostname + "/join";
const qr = new QRCode(document.getElementById("qrcode"), {
    text: domain,
    width: 125,
    height: 125,
});

window.addEventListener('load', () => {
    let difficulty = document.querySelectorAll('.difficulty');

    for (let i = 0; i < difficulty.length; i++) {
        let difficultySegment = difficulty[i].querySelectorAll('.piece')
        let complexity = difficulty[i].getAttribute('difficulty');

        for (let i = 0; i < 4; i++) {
            if (i < complexity) {
                difficultySegment[i].classList.add('segment_on')
            } else {
                difficultySegment[i].classList.remove('segment_on')
            }
        }
    }
});

async function sendGameStartInfoToServer() {
    let response = await fetch("/api/startGame", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `room_id=${window.location.pathname.split('/').pop()}`,
    });
    if (response.ok) {
        console.log('Отправил сообщение о начале игры');
    } else {
        console.log('Не получилось отправить сообщение о начале', response.status);
    }
}

async function sendGameEndInfoToServer() {
    let response = await fetch("/api/endGame", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `room_id=${window.location.pathname.split('/').pop()}`,
    });
    if (response.ok) {
        console.log('Отправил сообщение о конце игры');
    } else {
        console.log('Не получилось отправить сообщение о конче игры', response.status);
    }
}