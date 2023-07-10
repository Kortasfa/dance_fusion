const btnGo = document.getElementById("btn-join")
const enterInRoom = document.querySelector(".entrance-id-room__field")
const connectionText = document.querySelector(".connection")
const warningID = document.getElementById("id-warning");
const emptyID = document.getElementById("id-empty");
function sendMessage() {
    if (enterInRoom.value === "") {
        warningID.classList.add("hidden");
        emptyID.classList.remove("hidden");
        enterInRoom.classList.add("entrance-id-room__field_warning");
        console.log("1");
    }
    else {
        warningID.classList.add("hidden")
        enterInRoom.classList.remove("entrance-id-room__field_warning")
        let IDField = document.getElementById("id-field");
        let postInfo = {
            "userID": userID,
            "roomID": enterInRoom.value
        }
        let messageContent = JSON.stringify(postInfo);
        let XHR = new XMLHttpRequest();
        XHR.open("POST", "/api/join_to_room");
        XHR.onload = function () {
            if (XHR.status === 200) {
                btnGo.classList.add("hidden");
                enterInRoom.classList.add("hidden");
                connectionText.classList.remove("hidden");
                console.log("Connected to the room!");
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



