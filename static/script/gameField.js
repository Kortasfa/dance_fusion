let isBtnClicked = false;
let scoreGood = 50;
let scoreOk = 20;
let scorePerfect = 90;

const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnContinue = document.getElementById("btn-continue");
function getUsersByCookie() {
    let numberOfUser = 0;
    let allCookies = document.cookie;
    let cookiesArray = allCookies.split(';');
    for (let i = 0; i < cookiesArray.length; i++) {
        let cookie = cookiesArray[i].split('=')
        let name = cookie[0];
        let findUser = name.indexOf('User');
        if (findUser === 1) {
            let parts = cookie[1].split(',');
            let userID = parts[0];
            let userName = parts[1];
            let imgSrc = parts[2];
            numberOfUser = numberOfUser + 1;
            let userScore = document.getElementById('user-score' + numberOfUser);
            let indexUser = document.getElementById('hero' + numberOfUser);
            let indexUserName = document.getElementById('heroName' + numberOfUser);
            let indexUserImg = document.getElementById('heroImg' + numberOfUser);
            userScore.innerText = userName + ":";
            userScore.classList.remove('hidden');
            indexUser.classList.remove('hidden');
            indexUser.id = userID;
            indexUserImg.src =  '../' + imgSrc;
            indexUserName.innerText = userName;
        }
    }
}

function showStats() {
    addStats();
    modalElem.classList.remove("hidden");
    modalElem.classList.add("open");
    console.log("end video")
}

function addStats(){
    let score = document.querySelectorAll('.hero__score');
    let info = document.querySelectorAll('.pop-up-box__user-score');
    for (let i = 0; i < 4; i++){
        info[i].innerText = info[i].innerText + ' ' + score[i].innerText;
    }
}
getUsersByCookie();

danceVideo.addEventListener('ended', showStats);
btnContinue.addEventListener("click", function () {
    window.location.href = "/room";
})


// function openModalElem() {
//     modalElem.classList.add("open");
// }

// x

function AddScore(userID, Score){
    let user = document.getElementById(userID);
    let userScore = user.querySelector(".hero__score");
    userScore.innerText = parseInt(userScore.innerText) + Score;
    if (Score > scorePerfect){
        let effect = user.querySelector(".hero__rating-perfect");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    } else if (Score > scoreGood) {
        let effect = user.querySelector(".hero__rating-good");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    }else if (Score > scoreOk){
        let effect = user.querySelector(".hero__rating-ok");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    }else {
        let effect = user.querySelector(".hero__rating-x");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_bad");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    }
}
// function emulateClick(btn) {
//     let click = new CustomEvent("mousemove");
//     btn.dispatchEvent(click);
//     console.log("click!")
// }

// function playVideo() {
//     modalElem.classList.remove("open");
//     setTimeout(() => {
//         danceVideo.play();
//     }, 500);
// }
/*window.onload = emulateClick(btnGo);

/*function test() {
    console.log("test click")
}

btnGo.addEventListener("mousemove", test());*/

/*let socket = new WebSocket(WssURL);

socket.onopen = function(event) {
    console.log("WebSocket connection established.");
}

socket.onmessage = function() {
    //let message = event.data;
    /*x.removeAttribute('disabled');
    btnGo.click();*/
/*  playVideo();

}

socket.onclose = function(event) {
  console.log("WebSocket connection closed.");
}*/

// btnGo.addEventListener("click", playVideo);