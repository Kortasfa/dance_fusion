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

let value = 0;
const maxValue = 1000;
const percentageStep = 0.1;

function startAnimation() {
    if (value >= maxValue) return;

    value += 1;
    const percentage = (value / maxValue) * 100;
    const block = document.querySelector('.block');
    block.style.height = percentage + '%';

    requestAnimationFrame(startAnimation); // Рекурсивно вызываем функцию для создания плавной анимации
}