const userBody = document.querySelector('.user-body');
const userFace = document.querySelector('.user-face');
const userHat = document.querySelector('.user-hat')

function openElement(type){
    let nameType = '.' + type.id;
    let itemSelect = document.querySelectorAll(nameType);
    let element = document.querySelectorAll('.element');
    let item = document.querySelectorAll('.type');
    for(let i = 0;  i < 3; i++){
        item[i].classList.remove('type-selected')
    }
    type.classList.add('type-selected');

    for(let i = 0;  i < element.length; i++){
        element[i].classList.add('none')
    }
    for(let i = 0;  i < itemSelect.length; i++){
        itemSelect[i].classList.remove('none');
    }
    type.classList.add('type-selected')
}

function changeHat(hat){
    userHat.src = hat.querySelector(".element__hat").src;
}

function changeFace(hat){
    userFace.src = hat.querySelector(".element__face").src;
}

function changeBody(hat){
    userBody.src = hat.querySelector(".element__body").src;
}