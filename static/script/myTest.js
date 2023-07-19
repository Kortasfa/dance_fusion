// const avatarInput = document.getElementById("user-photo");
// const userPhoto = document.querySelector(".block");
//
// function addAvatar() {
//
//     const file = avatarInput.files[0];
//
//     if (file) {
//         const reader  = new FileReader();
//         reader.onload = function () {
//             userPhoto.style.backgroundImage = "url(" + reader.result + ")";
//         }
//         reader.readAsDataURL(file);
//     }
// }
//
// avatarInput.addEventListener("input", addAvatar)


const maxValue = 1000; value / MaxValue
let lenArray = (forgetYou); // длина массива - кол-во движений
let parametr=  Math.floor((lenArray / MaxValue) * 1000) / 1000);
let value = ;

function startAnimation() {
    if (value >= maxValue) return;
x
    value = 1;
    const percentage = (value / maxValue) * 100;
    const block = document.querySelector('.block');
    block.style.height = percentage + '%';

    requestAnimationFrame(startAnimation); // Рекурсивно вызываем функцию для создания плавной анимации
}