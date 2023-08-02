let isBtnClicked = false;
let scoreGood = 25;
let scoreOk = 13;
let scorePerfect = 33;

const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnContinue = document.getElementById("btn-continue");
const starOneScale = document.getElementById("star-1");
const starTwoScale = document.getElementById("star-2");
const starThreeScale = document.getElementById("star-3");
const starFourScale = document.getElementById("star-4");
const starFiveScale = document.getElementById("star-5");
const megaStarScale = document.getElementById("mega-star");

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
        let indexUserScale = document.getElementById("scale-" + (i + 1));
        indexUserScale.id = "scale-" + userID + "-for-user";
        indexUserScale.classList.remove("hidden");
        indexUserBodyImg.src = bodyImgSrc;
        indexUserFaceImg.src = faceImgSrc;
        indexUserHatImg.src = hatImgSrc;
        indexUserName.innerText = userName;
        let heroStars = document.getElementById("hero-score-" + (i + 1));
        let megaStar = heroStars.querySelectorAll(".score__star")[5];
        megaStar.id = "mega-star-" + userID;
        let starOne = heroStars.querySelectorAll(".score__star")[0];
        starOne.id = "star-1-" + userID;
        let starTwo = heroStars.querySelectorAll(".score__star")[1];
        starTwo.id = "star-2-" + userID;
        let starThree = heroStars.querySelectorAll(".score__star")[2];
        starThree.id = "star-3-" + userID;
        let starFour = heroStars.querySelectorAll(".score__star")[3];
        starFour.id = "star-4-" + userID;
        let starFive = heroStars.querySelectorAll(".score__star")[4];
        starFive.id = "star-5-" + userID;
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

let pix;
let percentage;

function addScore(userID, score, maxScore) {
    let valueScore;
    let starComplete = 0;
    let user = document.getElementById(userID);
    let scale = document.getElementById("scale-" + userID + "-for-user")
    let maxPractice = 0.8 * maxScore;
    let starOne = document.getElementById("star-1-" + userID);
    let starTwo = document.getElementById("star-2-" + userID);
    let starThree= document.getElementById("star-3-" + userID);
    let starFour = document.getElementById("star-4-" + userID);
    let starFive = document.getElementById("star-5-" + userID);
    let megaStar = document.getElementById("mega-star-" + userID);

    let userIndex;
    for (userIndex = 0; userIndex < connectedUsers.length; userIndex++) {
        let userInfo = connectedUsers[userIndex];
        if (userInfo["userID"] === userID) {
            valueScore = userInfo["valueScore"];
            console.log("test valueScore", valueScore)
            break;
        }
    }
    if (valueScore === undefined) {
        console.log("пользователь не найден с таким id")
        return;
    }
    valueScore += score;
    console.log("valueScore: " + valueScore);
    if (valueScore <= maxScore) {
        console.log("valueScore: " + valueScore);
        percentage = (valueScore / maxPractice);
        console.log("percentage: " + percentage);
        pix = 250 * percentage;
        console.log("pix: " + pix);
    }
    if (valueScore > maxScore) {
        percentage = (valueScore / maxScore);
        console.log("percentage: " + percentage);
        pix = 50 * percentage;
        console.log("pix: " + pix);
    }
    scale.style.height = pix + 'px';
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
        starOne.src = "/static/img/star_blue.svg";
        starComplete = 1;
    }
    if (valueScore >= 0.4 * maxPractice) {
        starTwo.src = "/static/img/star_blue.svg"
        starComplete = 2;
    }
    if (valueScore >= 0.6 * maxPractice) {
        starThree.src = "/static/img/star_blue.svg"
        starComplete = 3;
    }
    if (valueScore >= 0.8 * maxPractice) {
        starFour.src = "/static/img/star_blue.svg"
        starComplete = 4;
    }
    if (valueScore >= maxPractice) {
        starFive.src = "/static/img/star_blue.svg"
        starComplete = 5;
    }
    if (valueScore >= 0.9 * maxScore) {
        megaStar.classList.remove("hidden");
        starComplete = 6;
    }
    switch(starComplete) {
        case 1:
            starOneScale.src = "/static/img/star_blue.svg"
            break;
        case 2:
            starTwoScale.src = "/static/img/star_blue.svg"
            break;
        case 3:
            starThreeScale.src = "/static/img/star_blue.svg"
            break;
        case 4:
            starFourScale.src = "/static/img/star_blue.svg"
            break;
        case 5:
            starFiveScale.src = "/static/img/star_blue.svg"
            break;
        case 6:
            megaStarScale.classList.remove("hidden");
            break;
    }

    connectedUsers[userIndex]["valueScore"] = valueScore;
}
