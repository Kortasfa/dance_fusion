let isBtnClicked = false;
let scoreGood = 25;
let scoreOk = 13;
let scorePerfect = 33;

const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnContinue = document.getElementById("btn-continue");
function getUsersByCookie() {
    for (let i = 0; i < connectedUsers.length; i++) {
        let userID = connectedUsers[i]["userID"];
        let userName = connectedUsers[i]["userName"];
        let imgSrc = connectedUsers[i]["imgSrc"];
        let userScore = document.getElementById('user-score' + (i + 1));
        let indexUser = document.getElementById('hero' + (i + 1));
        let indexUserName = document.getElementById('heroName' + (i + 1));
        let indexUserImg = document.getElementById('heroImg' + (i + 1));
        userScore.innerText = userName + ":";
        userScore.classList.remove('hidden');
        indexUser.classList.remove('hidden');
        indexUser.id = userID;
        indexUserImg.src =  '../' + imgSrc;
        indexUserName.innerText = userName;
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

function addScore(userID, Score){
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
    } else if (Score > scoreOk){
        let effect = user.querySelector(".hero__rating-ok");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    } else {
        let effect = user.querySelector(".hero__rating-x");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_bad");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    }
}
