const userBody = document.querySelector('.user-body');
const userFace = document.querySelector('.user-face');
const userHat = document.querySelector('.user-hat');
const element = document.querySelectorAll('.element');
const item = document.querySelectorAll('.type');
const userLvl = parseInt(document.getElementById('level').innerText);
const saveButton = document.getElementById('saveButton');

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
    if (hat.classList.contains('unavailable')){} else
    userHat.src = hat.querySelector(".element__hat").getAttribute("src");
}

function changeFace(face){
    if (face.classList.contains('unavailable')){} else
    userFace.src = face.querySelector(".element__face").getAttribute("src");
}

function changeBody(body){
    if (body.classList.contains('unavailable')){} else
    userBody.src = body.querySelector(".element__body").getAttribute("src");
}
function verificationLvl() {
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
            console.log('Кастомизация принята');
        } else {
            console.log('Не удалось применить кастомизацию');
        }
    }
}

saveButton.addEventListener('click', changeAvatar)