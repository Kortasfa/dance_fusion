function signUp() {
    let nameField = document.querySelector(".name-field");
    let passwordField = document.querySelector(".password-field");
    let userInfo = {
        "userName": nameField.value,
        "password": passwordField.value
    }
    let messageContent = JSON.stringify(userInfo);
    let XHR = new XMLHttpRequest();
    XHR.open("POST", "/api/signUp");
    XHR.onload = function () {
        if (XHR.status === 200) {
            console.log("Successfully registered!");
        } else if (XHR.status === 409) {
            console.log("Username is taken!");
        } else {
            console.log("Failed to register!");
        }
    };
    XHR.send(messageContent);
}