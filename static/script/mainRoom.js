const listSong = document.getElementById('listSong');
const listGenre = document.getElementById('listGenre');
const btnOpenInfo = document.getElementById('openGuide');
const guide = document.getElementById('guide');
const returnBtn = document.getElementById('returnButton');
const PlayBtn = document.getElementById('play');
const songs = document.querySelector('.songs')

let readyGame = false;
let readyPlayer = false;
let connectedUsers = [];
let connectedBots = [];

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
    songs.classList.remove('none');
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
let songNeuro = '';
let fullSongName = '';
let songId = 0;
let difficulty = 0;
let readySong = false;

function onImageClick(element) {
    const videoSrcID = 'song' + element.id;
    const video = document.getElementById(videoSrcID);
    const fullVideo = document.getElementById('full' + videoSrcID);
    const videoPlayer = document.getElementById('videoPlayer');
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
    playButton.on('click', gameStart());
});

function gameStart() {
    const contentContainer = $('#content');
    if (readyGame) {
        numb = (Math.round(Math.random()*1000)).toString();
        let firstComponent = '/static/test/edge-impulse-standalone.js?version=' + numb;
        let secondComponent = '/static/test/run-impulse.js?version=' + numb;
        let thirdComponent = '/static/html/game.html?version=' + numb;
        $.getScript(firstComponent, function() {
            $.getScript(secondComponent, function() {
                (async () => {
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
    }
    return readyGame;
}

function showVideo(videoID) {
    let videoSrcID = 'song' + videoID.id;
    let video = document.getElementById(videoSrcID);
    let videoPlayer = document.getElementById('videoPlayer');
    videoPlayer.src = video.innerText;
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
                console.log('Bot Info:');
                console.log('Bot Id: -', data.BotId);
                console.log('Bot Scores Src:', data.BotScoresPath);
                console.log('Bot Hat Src:', data.BotImgHat);
                console.log('Bot Face Src:', data.BotImgFace);
                console.log('Bot Body Src:', data.BotImgBody);
                botInfo = data;
                addUser( "-" + data.BotId, "bot_1", "../" + data.BotImgHat, "../" + data.BotImgFace, "../" + data.BotImgBody);
                function readJSONFromURL(url) {
                    fetch(url)
                        .then(response => {
                            if (!response.ok) {
                                throw new Error('Ошибка при получении файла');
                            }
                            return response.json();
                        })
                        .then(jsonData => {
                            connectedBots.push({"botID": "-1",  "botScores":  jsonData});
                            console.log(connectedBots);
                        })
                        .catch(error => {
                            console.error('Ошибка:', error);
                        });
                }
                readJSONFromURL("../" + data.BotScoresPath);
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
    let healthPoint = bossBlock.querySelector(".boss__health-point").innerText;
    let bossBody= bossBlock.querySelector(".boss__body-img").src;
    let bossFace = bossBlock.querySelector(".boss__face-img").src;
    let bossHat = bossBlock.querySelector(".boss__hat-img").src;
    if (!gameStart()) {
        console.log("Игра не готова")
        return;
    }

    bossInfo = {"name": name, "healthPoint": healthPoint, "bossBody": bossBody, "bossFace": bossFace, "bossHat": bossHat};

}


const parent = document.querySelector('.songs');

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
    console.log("WebSocket mainRoom connection closed.");
};

function addUser(userID, userName, hatImgSrc, faceImgSrc, bodyImgSrc) {
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


const domain = window.location.protocol + "//" + window.location.hostname + "/join";
const qr = new QRCode(document.getElementById("qrcode"), {
    text: domain,
    width: 125,
    height: 125,
});

window.addEventListener('load', () => {
    let difficulty = songs.querySelectorAll('.song__difficulty');

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