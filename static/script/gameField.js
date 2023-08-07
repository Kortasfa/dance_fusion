let isBtnClicked = false;
let scoreGood = 20;
let scoreOk = 3;
let scorePerfect = 31;
let bossHp;
let bossHealth;

const danceVideo = document.getElementById("video-dance");
const modalElem = document.getElementById("pop-up");
const btnContinue = document.getElementById("btn-continue");
const starOneScale = document.getElementById("star-1");
const starTwoScale = document.getElementById("star-2");
const starThreeScale = document.getElementById("star-3");
const starFourScale = document.getElementById("star-4");
const starFiveScale = document.getElementById("star-5");
const megaStarScale = document.getElementById("mega-star");
const bossHpBar = document.querySelector('.hp-bar');

function getUsersByCookie() {
    for (let i = 0; i < connectedUsers.length; i++) {
        let userID = connectedUsers[i]["userID"];
        let userName = connectedUsers[i]["userName"];
        let bodyImgSrc = connectedUsers[i]["bodyImgSrc"];
        let faceImgSrc = connectedUsers[i]["faceImgSrc"];
        let hatImgSrc = connectedUsers[i]["hatImgSrc"];
        let indexUser = document.getElementById('hero' + (i + 1));
        let indexUserName = indexUser.querySelector(".hero__name");
        let indexUserBodyImg = indexUser.querySelector(".body");
        let indexUserFaceImg = indexUser.querySelector(".face");
        let indexUserHatImg = indexUser.querySelector(".hat");
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
}

function addStats(){
    getBestPlayer(songId);

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
            break;
        }
    }
    if (valueScore === undefined) {
        return;
    }
    valueScore += score;
    if (valueScore <= maxScore) {
        percentage = (valueScore / maxPractice);
        pix = 250 * percentage;
    }
    if (valueScore > maxScore) {
        percentage = (valueScore / maxScore);
        pix = 50 * percentage;
    }
    scale.style.height = pix + 'px';
    if (score > scorePerfect){
        let effect = user.querySelector(".hero__rating-perfect");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 2000);
    } else if (score > scoreGood) {
        let effect = user.querySelector(".hero__rating-good");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 2000);
    } else if (score > scoreOk){
        let effect = user.querySelector(".hero__rating-ok");
        effect.classList.remove("hidden");
        effect.classList.add("hero__rating_visible");
        setTimeout(function() {
            effect.classList.remove("hero__rating_visible");
            effect.classList.add("hidden")
        }, 2000);
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
    if (mode == 'Boss') {
        playerDamage(score);
    }
}

let hpBar;
let bossStage = 1;
function playerDamage(score) {
    const hpPercentage = (bossHp / bossHealth) * 100;
    if (hpPercentage > 66) {
        bossStage = 1;
    } else if (hpPercentage > 33) {
        bossStage = 2;
    } else {
        bossStage = 3;
    }

    bossHp = bossHp - score;
    hpBar =  192 * (bossHp / bossHealth);
    bossHpBar.style.width = hpBar + 'px';

    if (bossHp < 0) {
        mode = 'Classic';
        bossHpBar.classList.add('none');
        document.querySelector('.boss__head').classList.add('defeated')
    }
}

let btnExit = document.querySelector(".btn-exit");
btnExit.addEventListener("click", expelUsers)
window.onbeforeunload = function (){
    sendGameEndInfoToServer()
        .then(() => {expelUsers();})
}

function expelUsers() {
    for (let user of connectedUsers) {
        if (parseInt(user["userID"]) > 0) {
            expelUser(user["userID"]).then(() => {});
        }
    }
}

if (mode == 'Boss') {
    document.querySelector(".boss-container").classList.remove("hidden");
    document.querySelector(".boss__name").innerText = bossInfo.name;
    document.querySelector(".boss__body-img").src = bossInfo.bossBody;
    document.querySelector(".boss__face-img").src = bossInfo.bossFace;
    document.querySelector(".boss__hat-img").src = bossInfo.bossHat;
    bossHealth = parseInt(bossInfo.healthPoint) * connectedUsers.length;
    bossHp = bossHealth;
}

async function expelUser(userID) {
    let response = await fetch("/api/deletePlayerFromGame", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `user_id=${userID}`,
    });
    if (!response.ok) {
        console.log('Не получилось отправить сообщение о выходе ', userID);
    }
    if (socket) {
        socket.close();
    }
}

async function getAchievements() {
    for (let i = 0; i < connectedUsers.length; i++) {
        if (connectedUsers[i]["userID"] > 0) {
            let botsID = [];
            for (let j = i + 1; j < connectedUsers.length; j++) {
                if (connectedUsers[j]["userID"] < 0) {
                    botsID.push(parseInt(connectedUsers[j]["userID"]) * (-1));
                }
            }
            let bossID = 0;
            if (bossInfo) {
                if (bossHp <= 0) {
                    bossID = parseInt(bossInfo.bossId);
                }
            }
            let jsonData = JSON.stringify({
                "user_id": parseInt(connectedUsers[i]["userID"]), // Integer
                "song_id": parseInt(songId), // Integer
                "bot_ids": botsID,
                "boss_id": bossID
            })
            const response = await fetch("/api/checkForAchievements", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: jsonData,
            });
            if (!response.ok) {
                console.log('Не удалось отправить данные для проверки на достижение');
            }
        }
    }
}

let bossContainer;
let bossThought;
let bossSticker;
let bossEmotionJSON;
if (bossInfo) {
    bossContainer = document.querySelector('.boss-container');
    bossThought = document.querySelector('.boss-container__boss-thought');
    bossSticker = bossThought.querySelector('.boss-thought__sticker');
    readJSONFromURL("../static/boss_emotions/boss_emotion_list.json").then(jsonData => {
        bossEmotionJSON = jsonData;
        for (let bossEmotionDict of bossEmotionJSON) {
            if (bossEmotionDict["stage"] === bossStage) {
                bossSticker.src = "../" + bossEmotionDict["emotions"][Math.floor(Math.random() * bossEmotionDict["emotions"].length)];
                break;
            }
        }
        showBossThought();
    });
}

function showBossThought() {
    if (bossHp < 0) {
        return;
    }
    setTimeout(function() {
        bossContainer.classList.add('boss-show-thought');
        setTimeout(function() {
            bossContainer.classList.remove('boss-show-thought');
            setTimeout(function() {
                for (let bossEmotionDict of bossEmotionJSON) {
                    if (bossEmotionDict["stage"] === bossStage) {
                        bossSticker.src = "../" + bossEmotionDict["emotions"][Math.floor(Math.random() * bossEmotionDict["emotions"].length)];
                        break;
                    }
                }
                setTimeout(showBossThought, 1000);
            }, 1000);
        }, 2000);
    }, Math.floor(Math.random() * (8000 - 2000 + 1)) + 2000);
}