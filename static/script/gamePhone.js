const btnGo = document.getElementById("btn-join")

function sendMessage() {
    let IDField = document.querySelector('.entrance__entrance-id-room');
    let postInfo = {
        "userID": userID,
        "roomID": IDField.value
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
}


btnGo.addEventListener("click", sendMessage);




