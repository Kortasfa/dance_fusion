const passwordEye = document.querySelector(".field__eye");
const nameField = document.getElementById("name");
const passwordField = document.getElementById("password");
const logInBtn = document.getElementById("btn-sign-up");
const warningMessage = document.querySelectorAll(".field__warning")

function logIn() {
    let userInfo = {
        "userName": nameField.value,
        "password": passwordField.value
    }
    let messageContent = JSON.stringify(userInfo);
    let XHR = new XMLHttpRequest();
    XHR.open("POST", "/api/logIn");
    XHR.onload = function () {
        if (XHR.status === 200) {
            alert("Successfully logged in!");
        } else if (XHR.status === 409) {
            warningMessage.forEach(element => element.classList.remove("hidden"));
            nameField.classList.add("warning-input");
            passwordField.classList.add("warning-input");
            console.log("Wrong username or password!");
        } else {
            console.log("Failed to log in!");
        }
    };
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
logInBtn.addEventListener("click", logIn)
