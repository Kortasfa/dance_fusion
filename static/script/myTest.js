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


let maxTheory = 5600;
let maxPractice = maxTheory - maxTheory * 0.2; //4480
let value = 0;

function startAnimation() {
    if (value >= 30) return
    value +=10;
    console.log(value);
    const percentage = (value / maxPractice) * 100;
    const block = document.querySelector('.block');
    block.style.height = percentage + '%';
    requestAnimationFrame(startAnimation); // Рекурсивно вызываем функцию для создания плавной анимации
}
// let value = 0
function anim(score) {
    console.log("score: " + score);
    if (value > 5600) return
    value += score;
    console.log("value: " + value);
    const percentage = (value / maxPractice);
    console.log("percentage: " + percentage);
    let pix = 250 * percentage;
    console.log("pix: " + pix);
    const block = document.querySelector('.block');
    block.style.height = pix + 'px';// Рекурсивно вызываем функцию для создания плавной анимации
}
