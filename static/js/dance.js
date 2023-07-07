const song__list = document.getElementById('list-song');
const musics = document.getElementById('list-genry');

function openSong(styleButtonBlock) {
  musics.classList.add('none');
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
  song__list.classList.add('none');
  musics.classList.remove('none');
}

function startGame() {
  window.location.href = 'homepage-url';
  //проверка выбрана ли песня
  //отправка в выбранную игру
}