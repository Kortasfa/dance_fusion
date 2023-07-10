const listSong = document.getElementById('list-song');
const listGenre = document.getElementById('list-genre');
const btnOpenInfo =  document.getElementById('openGuide');
const guide = document.getElementById('guide');
let numberOfUser = 0;

btnOpenInfo.addEventListener('click', openGuide);

function openSong(styleButtonBlock) {
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
  guide.classList.remove('unplay');
}

function closeGuide() {
  guide.classList.add('unplay');
  guide.classList.remove('play');
}

function startGame() {
  window.location.href = 'homepage-url';
  //проверка выбрана ли песня
  //отправка в выбранную игру
}

function showVideo(videoID) {
  let videoSrcID = '9' + videoID.id;
  let video = document.getElementById(videoSrcID);
  let videoPlayer = document.getElementById('videoPlayer');
  videoPlayer.src = video.innerText;
}

let socket = new WebSocket(WssURL);

socket.onopen = function(event) {
  console.log("WebSocket connection established.");
};

socket.onmessage = function(event) {
  let userMSG = document.getElementById('needUser')
  let message = event.data;
  let parts = message.split('|');
  let userID = parts[0];
  let userName = parts[1];
  let imgSrc = parts[2];
  numberOfUser = numberOfUser + 1;
  userMSG.classList.add('none');
  let indexUser = document.getElementById('user' + numberOfUser);
  let indexUserName = document.getElementById('userName' + numberOfUser);
  indexUser.classList.remove('none');

  indexUserName.innerText = userName;
  console.log('Пользователь присоединился: ' + userID);
  console.log('Его имя: ' + userName);
  console.log('Его фотка: ' + imgSrc);
};

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
};