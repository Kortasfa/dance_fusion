const btnGo = document.getElementById("btn-join")
const enterInRoom = document.querySelector(".entrance-id-room__field")
const connectionText = document.querySelector(".connection")
const warningID = document.getElementById("id-warning");
const emptyID = document.getElementById("id-empty");
function getCookieValue(cookieName) {
    let allCookies = document.cookie;
    let cookiesArray = allCookies.split(';');
    let name = cookieName + "=";
    for (let i = 0; i < cookiesArray.length; i++) {
        let cookie = cookiesArray[i].trim();
        if (cookie.indexOf(name) === 0) {
            return cookie.substring(name.length, cookie.length);
        }
    }
    return "";
}
function setCookie(cookieName, value, days) {
    let d = new Date();
    d.setTime(d.getTime() + (days*24*60*60*1000));
    let expires = "expires="+ d.toUTCString();
    document.cookie = cookieName + "=" + value + ";" + expires + ";path=/";
}
function sendMessage() {
    if (enterInRoom.value === "") {
        warningID.classList.add("hidden");
        emptyID.classList.remove("hidden");
        enterInRoom.classList.add("entrance-id-room__field_warning");
    }
    else {
        let authCookieValue = getCookieValue("authCookieName");
        if (authCookieValue === "") {
            console.log("Login to your account!");
            return
        }
        let roomCookieValue = getCookieValue("roomCookieName");
        if (roomCookieValue !== "") {
            if (enterInRoom.value !== roomCookieValue) {
                console.log("You are already connected to another room!");
            }
            else {
                console.log("You are already connected to this room!");
            }
            return
        }
        warningID.classList.add("hidden")
        enterInRoom.classList.remove("entrance-id-room__field_warning")
        let IDField = document.getElementById("id-field");
        let postInfo = {
            "userID": authCookieValue,
            "roomID": enterInRoom.value
        }
        let messageContent = JSON.stringify(postInfo);
        let XHR = new XMLHttpRequest();
        XHR.open("POST", "/api/joinToRoom");
        XHR.onload = function () {
            if (XHR.status === 200) {
                btnGo.classList.add("hidden");
                enterInRoom.classList.add("hidden");
                connectionText.classList.remove("hidden");
                emptyID.classList.add("hidden");
                console.log("Connected to the room!");
                setCookie("roomCookieName", enterInRoom.value, 1)
            } else if (XHR.status === 404) {
                emptyID.classList.add("hidden");
                warningID.classList.remove("hidden");
                enterInRoom.classList.add("entrance-id-room__field_warning");
                console.log("Room ID not found!");
            } else if (XHR.status === 409) {
                console.log("The room is full!");
            } else {
                console.log("Failed to send room id");
            }
        };
        XHR.send(messageContent);
    }
}
