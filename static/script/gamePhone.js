const btnGo = document.getElementById("btn-join")
const enterInRoom = document.querySelector(".entrance-id-room__field")
const connectionText = document.querySelector(".connection")
const warningID = document.getElementById("id-warning");
const emptyID = document.getElementById("id-empty");

function setJsonCookie(name, value, expirationDays) {
    const jsonValue = JSON.stringify(value);
    const encodedValue = encodeURIComponent(jsonValue);
    document.cookie = `${name}=${encodedValue}; path=/; expires=${getExpirationDate(expirationDays)}`;
}

function getJsonCookie(name) {
    const cookies = document.cookie.split(';');
    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + '=')) {
            const encodedValue = cookie.substring(name.length + 1);
            const decodedValue = decodeURIComponent(encodedValue);
            return JSON.parse(decodedValue);
        }
    }
    return null;
}

function getExpirationDate(expirationDays) {
    const date = new Date();
    date.setDate(date.getDate() + expirationDays);
    return date.toUTCString();
}

const userInfo = getJsonCookie("userInfoCookie");
document.querySelector('.user__name').textContent = userInfo.UserName;
document.querySelector('.user__avatar').src = userInfo.ImgSrc;

function sendMessage() {
    if (enterInRoom.value === "") {
        warningID.classList.add("hidden");
        emptyID.classList.remove("hidden");
        enterInRoom.classList.add("entrance-id-room__field_warning");
    }
    else {
        if (userInfo === null) {
            console.log("Login to your account!");
            return
        }
        if (userInfo.SelectedRoom !== "") {
            if (enterInRoom.value !== userInfo.SelectedRoom) {
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
            "userID": userInfo.UserID,
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
                //setCookie("roomCookieName", enterInRoom.value, 1)
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

btnGo.addEventListener("click", sendMessage);