const btnGo = document.getElementById("btn-join");
const enterInRoom = document.querySelector(".entrance-id-room__field");
const connectionText = document.querySelector(".connection");
const entranceField = document.querySelector(".entrance");
const warningID = document.getElementById("id-warning");
const emptyID = document.getElementById("id-empty");
const fullID = document.getElementById("id-full");
const btnLogOut = document.querySelector(".btn-log-out");

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
        fullID.classList.add("hidden");
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
                entranceField.classList.add("hidden");
                connectionText.classList.remove("hidden");
                emptyID.classList.add("hidden");
                fullID.classList.add("hidden");
                warningID.classList.add("hidden");
                console.log("Connected to the room!");
            } else if (XHR.status === 404) {
                emptyID.classList.add("hidden");
                warningID.classList.remove("hidden");
                fullID.classList.add("hidden");
                enterInRoom.classList.add("entrance-id-room__field_warning");
                console.log("Room ID not found!");
            } else if (XHR.status === 409) {
                emptyID.classList.add("hidden");
                fullID.classList.remove("hidden");
                warningID.classList.add("hidden");
                enterInRoom.classList.add("entrance-id-room__field_warning");
                console.log("The room is full!");
            } else {
                alert("Failed to send room id");
            }
        };
        XHR.send(messageContent);
    }
}

async function logout() {
    const response = await fetch("/clear");
    if (response.ok) {
        window.location.href = "/logIn";
    }
}

btnLogOut.addEventListener("click", logout)
btnGo.addEventListener("click", sendMessage);