const passwordEye = document.querySelector(".field__eye");
const nameField = document.getElementById("name");
const passwordField = document.getElementById("password");
const signUpBtn = document.getElementById("btn-sign-up");
const warningUsername = document.getElementById("username-taken");
const usernameEmpty = document.getElementById("username-empty");
const passwordEmpty = document.getElementById("password-empty");

nameField.oninput = function() {
    this.value = this.value.substr(0, 10);
}
function signUp() {
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
        warningUsername.classList.add("hidden");
    }
    else if ((nameField.value !== "") && (passwordField.value === "")) {
        nameField.classList.remove("warning-input");
        usernameEmpty.classList.add("hidden");
        passwordField.classList.add("warning-input");
        passwordEmpty.classList.remove("hidden");

    }
    else if ((nameField.value === "") && (passwordField.value === "")) {
        nameField.classList.add("warning-input");
        usernameEmpty.classList.remove("hidden");
        passwordField.classList.add("warning-input");
        passwordEmpty.classList.remove("hidden");
        warningUsername.classList.add("hidden");
    }
    else {

        XHR.open("POST", "/api/signUp");
        XHR.onload = function () {
            if (XHR.status === 200) {
                window.location.href = "/logIn";
                console.log("Successfully registered!");
            } else if (XHR.status === 409) {
                console.clear();
                usernameEmpty.classList.add("hidden")
                nameField.classList.add("warning-input");
                warningUsername.classList.remove("hidden");
                passwordField.classList.remove("warning-input");
                passwordEmpty.classList.add("hidden");
                console.log("Username is taken!");
            } else {
                alert("Failed to register!");
            }
        }
        XHR.send(messageContent);
    }
}

function changeEye(){
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
