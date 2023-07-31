let isBtnClicked = false;
let scoreGood = 25;
let scoreOk = 13;
let scorePerfect = 33;

const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnContinue = document.getElementById("btn-continue");
const megaStar = document.querySelectorAll(".score__star")[5];
const starOne = document.querySelectorAll(".score__star")[0];
const starTwo = document.querySelectorAll(".score__star")[1];
const starThree = document.querySelectorAll(".score__star")[2];
const starFour = document.querySelectorAll(".score__star")[3];
const starFive = document.querySelectorAll(".score__star")[4];
function getUsersByCookie() {
    for (let i = 0; i < connectedUsers.length; i++) {
        let userID = connectedUsers[i]["userID"];
        let userName = connectedUsers[i]["userName"];
        let bodyImgSrc = connectedUsers[i]["bodyImgSrc"];
        let faceImgSrc = connectedUsers[i]["faceImgSrc"];
        let hatImgSrc = connectedUsers[i]["hatImgSrc"];
        let userScore = document.getElementById('user-score' + (i + 1));
        let indexUser = document.getElementById('hero' + (i + 1));
        let indexUserName = indexUser.querySelector(".hero__name");
        let indexUserBodyImg = indexUser.querySelector(".body");
        let indexUserFaceImg = indexUser.querySelector(".face");
        let indexUserHatImg = indexUser.querySelector(".hat");
        userScore.innerText = userName + ":";
        userScore.classList.remove('hidden');
        indexUser.classList.remove('hidden');
        indexUser.id = userID;
        indexUserBodyImg.src = bodyImgSrc;
        indexUserFaceImg.src = faceImgSrc;
        indexUserHatImg.src = hatImgSrc;
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
    // let score = document.querySelectorAll('.hero__score');
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

let valueScore= 0;
function addScore(userID, score, maxScore){
    let user = document.getElementById(userID);
    let maxPractice = maxScore - 0.2 * maxScore;
    valueScore += score;
    console.log("valueScore: ", valueScore);
    if (score > scorePerfect){
        let effect = user.querySelector(".hero__rating-perfect");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    } else if (score > scoreGood) {
        let effect = user.querySelector(".hero__rating-good");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 1000);
    } else if (score > scoreOk){
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
    if (valueScore >= 0.2 * maxPractice) {
        starOne.src = "/static/img/star_blue.svg"
    }
    if (valueScore >= 0.4 * maxPractice) {
        starTwo.src = "/static/img/star_blue.svg"
    }
    if (valueScore >= 0.6 * maxPractice) {
        starThree.src = "/static/img/star_blue.svg"
    }
    if (valueScore >= 0.8 * maxPractice) {
        starFour.src = "/static/img/star_blue.svg"
    }
    if (valueScore >= maxPractice) {
        starFive.src = "/static/img/star_blue.svg"
    }
    if (valueScore >= 0.9 * maxScore) {
        megaStar.classList.remove("hidden");
    }
    if (bossInfo) {
        playerDamage(score);
    }
}

function playerDamage(score) {
    let bossHPCount = document.querySelector(".boss__hp-bar");
    bossHPCount.innerText = (parseInt(bossHPCount.innerText) - score).toString();
}

if (bossInfo) {
    console.log("zahar best");
    document.querySelector(".boss-container").classList.remove("hidden");
    document.querySelector(".boss__name").innerText = bossInfo.name;
    document.querySelector(".boss__hp-bar").innerText = (parseInt(bossInfo.healthPoint) * connectedUsers.length).toString();
    document.querySelector(".boss__body-img").src = bossInfo.bossBody;
    document.querySelector(".boss__face-img").src = bossInfo.bossFace;
    document.querySelector(".boss__hat-img").src = bossInfo.bossHat;
}
