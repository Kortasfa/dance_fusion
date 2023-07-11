function logIn() {
    let nameField = document.querySelector(".name-field");
    let passwordField = document.querySelector(".password-field");
    let userInfo = {
        "userName": nameField.value,
        "password": passwordField.value
    }
    let messageContent = JSON.stringify(userInfo);
    let XHR = new XMLHttpRequest();
    XHR.open("POST", "/api/logIn");
    XHR.onload = function () {
        if (XHR.status === 200) {
            console.log("Successfully logged in!");
        } else if (XHR.status === 409) {
            console.log("Wrong username or password!");
        } else {
            console.log("Failed to log in!");
        }
    };
    XHR.send(messageContent);
}