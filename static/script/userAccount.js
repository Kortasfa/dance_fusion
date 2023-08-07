const userBody = document.querySelector('.user-body');
const userFace = document.querySelector('.user-face');
const userHat = document.querySelector('.user-hat');
const element = document.querySelectorAll('.element');
const item = document.querySelectorAll('.type');
let userLvl = 0;
const saveButton = document.getElementById('saveButton');
parseInt(document.getElementById('level').innerText)
function getJsonCookie(name){
    const cookies = document.cookie.split(';');
    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + '=')) {
            const encodedValue = cookie.substring(name.length + 1);
            const decodedValue = decodeURIComponent(encodedValue);
            return JSON.parse(decodedValue);
        }
    }
    return null;
}

const userInfo = getJsonCookie("userInfoCookie");
userHat.src = userInfo.HatSrc;
userFace.src = userInfo.FaceSrc;
userBody.src = userInfo.BodySrc;
document.querySelector('.avatar__name').innerText = userInfo.UserName;
function openElement(type){
    let nameType = '.' + type.id;
    let itemSelect = document.querySelectorAll(nameType);
    for(let i = 0;  i < 3; i++){                  //выделение выбранного окна с изменением
        item[i].classList.remove('type-selected')
    }
    type.classList.add('type-selected');

    for(let i = 0;  i < element.length; i++){   //скрываем все элементы
        element[i].classList.add('none')
    }
    for(let i = 0;  i < itemSelect.length; i++){  //открываем выбранные элементы
        itemSelect[i].classList.remove('none');
    }
}

function changeHat(hat){
    if (!hat.classList.contains('unavailable')) {
        userHat.src = hat.querySelector(".element__hat").getAttribute("src");
    }
}

function changeFace(face){
    if (!face.classList.contains('unavailable')) {
        userFace.src = face.querySelector(".element__face").getAttribute("src");
    }
}

function changeBody(body){
    if (!body.classList.contains('unavailable')){
        userBody.src = body.querySelector(".element__body").getAttribute("src");
    }
}
function verificationLvl() {
    let lvlScore = 3000;
    while (userScore > lvlScore) {
        userLvl++;
        userScore = userScore - lvlScore;
        lvlScore = Math.round(lvlScore * 1.1);
    }
    document.getElementById('level').innerText = userLvl;
    for(let i = 0;  i < element.length; i++){
        if (userLvl < element[i].getAttribute('lvl'))
        element[i].classList.add('unavailable')
    }
}
verificationLvl()


async function changeAvatar() {

    if ((userHat.getAttribute("src") !== null) && (userFace.getAttribute("src") !== null) && (userBody.getAttribute("src") !== null)){
        const dataToSend = {
            hatSrc: userHat.getAttribute("src"),
            faceSrc: userFace.getAttribute("src"),
            bodySrc: userBody.getAttribute("src"),
        };
        const response = await fetch("/api/custom", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(dataToSend),
        });
        if (response.ok) {
            window.location.href = 'join';
        } else {
            console.log('Не удалось применить кастомизацию');
        }
    }
}

saveButton.addEventListener('click', changeAvatar)

async function changeUserName(userID, newUserName) {
    const response = await fetch("/api/changeUserName", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"user_id": userID,
            "new_user_name": newUserName}),
    });
    if (!response.ok) {
        console.log('Username updated successfully');
    } else if (response.status === 409) {
        console.log("Username already taken!!!");
    } else {
        console.log("Name change error: " + response.status);
    }
}

async function changeUserPassword(userID, newUserPassword) {
    const response = await fetch("/api/changeUserPassword", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"user_id": userID,
            "new_user_password": newUserPassword}),
    });
    if (response.ok) {
        console.log('Password updated successfully');
    } else {
        console.log("Password change error: " + response.status);
    }
}