const passwordEye = document.querySelector(".field__eye");
const nameField = document.getElementById("name");
const passwordField = document.getElementById("password");
const logInBtn = document.getElementById("btn-sign-up");
const warningMessage = document.querySelectorAll(".field__warning-wrong")
const usernameEmpty = document.getElementById("username-empty");
const passwordEmpty = document.getElementById("password-empty");
const btnSignUp = document.getElementById("sign-up-link");
function logIn() {
    let userInfo = {
        "userName": nameField.value,
        "password": passwordField.value
    }
    let messageContent = JSON.stringify(userInfo);
    let XHR = new XMLHttpRequest();
    if ((nameField.value === "") && (passwordField.value !== "")) {
        nameField.classList.add("warning-input");
        usernameEmpty.classList.remove("hidden");
        passwordField.classList.remove("warning-input");
        passwordEmpty.classList.add("hidden");
        warningMessage.forEach(element => element.classList.add("hidden"));
    }
    else if ((nameField.value !== "") && (passwordField.value === "")) {
        nameField.classList.remove("warning-input");
        usernameEmpty.classList.add("hidden");
        passwordField.classList.add("warning-input");
        passwordEmpty.classList.remove("hidden");
        warningMessage.forEach(element => element.classList.add("hidden"));
    }
    else if ((nameField.value === "") && (passwordField.value === "")) {
        nameField.classList.add("warning-input");
        usernameEmpty.classList.remove("hidden");
        passwordField.classList.add("warning-input");
        passwordEmpty.classList.remove("hidden");
        warningMessage.forEach(element => element.classList.add("hidden"));
    }
    else {
        XHR.open("POST", "/api/logIn");
        XHR.onload = function () {
            if (XHR.status === 200) {
                console.log("Successfully logged in!");
                window.location.href = '/join';
            } else if (XHR.status === 409) {
                warningMessage.forEach(element => element.classList.remove("hidden"));
                nameField.classList.add("warning-input");
                passwordField.classList.add("warning-input");
                usernameEmpty.classList.add("hidden");
                passwordEmpty.classList.add("hidden");
                console.log("Wrong username or password!");
            } else {
                alert("Failed to log in!");
            }
        };
        XHR.send(messageContent);
    }
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

function goToSignUp() {
    window.location.href = '/signUp';
}

passwordField.addEventListener("input", eye);
passwordEye.addEventListener("click", changeEye);
logInBtn.addEventListener("click", logIn);
btnSignUp.addEventListener("click", goToSignUp);