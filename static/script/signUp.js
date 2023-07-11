const passwordEye = document.querySelector(".field__eye");
const nameField = document.getElementById("name");
const passwordField = document.getElementById("password");
const signUpField = document.querySelector(".entrance");
const signUpBtn = document.getElementById("btn-sign-up");
const warningMessage = document.querySelector(".field__warning")
const messageToLogIn = document.querySelector(".message-log-in")
function signUp() {
    let userInfo = {
        "userName": nameField.value,
        "password": passwordField.value
    }
    let messageContent = JSON.stringify(userInfo);
    let XHR = new XMLHttpRequest();
    XHR.open("POST", "/api/signUp");
    XHR.onload = function () {
        if (XHR.status === 200) {
            signUpField.classList.add("hidden");
            messageToLogIn.classList.remove("hidden");
            console.log("Successfully registered!");
        } else if (XHR.status === 409) {
            nameField.classList.add("warning-border");
            warningMessage.classList.remove("hidden")
            console.log("Username is taken!");
        } else {
            alert("Failed to register!");
        }
    }
    XHR.send(messageContent);
}

function changeEye() {
    console.log(passwordField.type);
    if (passwordField.type === "password") {
        passwordField.type = "text";
        passwordEye.src = "/static/img/eye.svg";
    }
    else if (passwordField.type === "text") {
        passwordField.type = "password";
        passwordEye.src = "/static/img/eye_off.svg";
    }
}

function eye() {
    if (passwordField.value === "") {
        passwordEye.classList.add("hidden")
    }
    else {
        passwordEye.classList.remove("hidden")
    }
}

passwordField.addEventListener("input", eye);
passwordEye.addEventListener("click", changeEye);
signUpBtn.addEventListener("click", signUp)
