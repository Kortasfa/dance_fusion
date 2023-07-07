const btnGo = document.getElementById("btn-join")
const enterInRoom = document.querySelector(".entrance__entrance-id-room")
const connectionText = document.querySelector(".connection")
/*async function startGame(){
    let xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/start");
    xhr.onload = function () {
        if (xhr.status === 200) {
            console.log("Connected to the room!");
        } else {
            console.log("Failed to send room id");
        }
    };
    xhr.send("message=True");

    btnGo.classList.add("hidden");
    enterIdRoom.classList.add("hidden");
    connectionText.classList.remove("hidden");
}*/

function sendMessage() {
    let IDField = document.getElementById("id-field");
    let userID = "{{ .UserID }}";
    let postInfo = {
        "userID": userID,
        "roomID": enterInRoom.value
    }
    let messageContent = JSON.stringify(postInfo);
    let XHR = new XMLHttpRequest();
    XHR.open("POST", "/api/join_to_room");
    XHR.onload = function () {
        if (XHR.status === 200) {
            console.log("Connected to the room!");
        } else if (XHR.status === 404) {
            console.log("Post ID not found!");
        } else if (XHR.status === 409) {
            console.log("The room is full!");
        } else {
            console.log("Failed to send room id");
        }
    };
    XHR.send(messageContent);

    btnGo.classList.add("hidden");
    enterInRoom.classList.add("hidden");
    connectionText.classList.remove("hidden");
}


btnGo.addEventListener("click", sendMessage);




