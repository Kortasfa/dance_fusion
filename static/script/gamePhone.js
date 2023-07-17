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

                let socket = new WebSocket("wss://" + window.location.hostname + "/ws/joinToRoom/" + userInfo.UserID);

                socket.onopen = function(event) {
                    console.log("WebSocket connection established.");
                };

                socket.onmessage = function(event) {
                    let receivedData = event.data;
                    if (receivedData === 'pause') {
                        console.log('pause');
                    }
                    else if (receivedData === 'resume') {
                        console.log('resume');
                    }
                    else {
                        console.log('Motions:');
                        console.log(JSON.parse(receivedData))//////////////////////////////
                        record()

                    }
                };

                socket.onclose = function(event) {
                    console.log("WebSocket connection closed.");
                };

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







function record() {
    // Check if the device supports the required sensors
    if (window.DeviceMotionEvent && window.DeviceOrientationEvent) {
        const sensorFrequency = 62.5; // Frequency in samples per second
        const interval = 1000 / sensorFrequency; // Interval in milliseconds

        let sensorData = [];

        // Event handler for receiving sensor data
        function handleSensorData(event) {
            const {alpha, beta, gamma} = event.rotationRate; // Gyroscope data
            const {x, y, z} = event.accelerationIncludingGravity; // Accelerometer data

            // Store the sensor data
            sensorData.push({
                alpha,
                beta,
                gamma,
                x,
                y,
                z,
            });
        }

        // Start recording sensor data
        function startRecording() {
            window.addEventListener('devicemotion', handleSensorData, true);
            window.addEventListener('deviceorientation', handleSensorData, true);
            setTimeout(1000);
            setTimeout(stopRecording, 2000);
        }

        function stopRecording(name) {
            window.removeEventListener('devicemotion', handleSensorData, true);
            window.removeEventListener('deviceorientation', handleSensorData, true);

            // Process the recorded data
            let outputString = '0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, ';

            // Loop through the recorded sensor data
            for (let i = 0; i < sensorData.length; i++) {
                const {x, y, z, alpha, beta, gamma} = sensorData[i];

                // Round the values using MyRound10 function
                const roundedX = MyRound10(x);
                const roundedY = MyRound10(y);
                const roundedZ = MyRound10(z);
                const roundedAlpha = MyRound10(alpha);
                const roundedBeta = MyRound10(beta);
                const roundedGamma = MyRound10(gamma);

                // Append the rounded values to the output string
                outputString += `${roundedX}, ${roundedY}, ${roundedZ}, ${roundedAlpha}, ${roundedBeta}, ${roundedGamma}, `;
            }

            // Log the output string
            document.write(outputString);
            let file = {"string": 0, "name": 123};
            let data = JSON.stringify(file);
            sendDataToServer(data);
        }
    } else {
        console.log('The device does not support required sensors.');
    }

    function MyRound10(val) {
        const roundedVal = Math.round(val * 10) / 10;
        const formattedVal = roundedVal.toFixed(4);
        return formattedVal.padStart(6, '0');
    }

    function sendDataToServer(data) {
        // Replace the URL with the appropriate endpoint to handle the data on your server
        let url = '/api/motion';

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: data
        })
            .then(function (response) {
                if (response.ok) {
                    console.log('Data sent successfully.');
                } else {
                    console.log('Error sending data. Status:', response.status);
                }
            })
            .catch(function (error) {
                console.log('Error sending data:', error);
            });
    }

    // Function to handle danceData sequentially
    function handleDanceData(index, danceData) {
        if (index >= danceData.length) {
            // If all danceData objects have been processed, stop recording and exit
            stopRecording();
            return;
        }

        const currentDanceData = danceData[index];
        const {start_time, duration, name} = currentDanceData;

        // Start recording after the start time
        setTimeout(() => {
            startRecording();
            // Stop recording after the duration
            setTimeout(() => {
                // Send data to the server
                stopRecording(name);
                // Move to the next danceData object
                handleDanceData(index + 1, danceData);
            }, duration * 1000);
        }, (start_time - danceData[0].start_time) * 1000);
    }

    fetch('static/motion_list/forgetYou.json')
        .then(response => response.json())
        .then(danceData => {
            // Start handling danceData from the beginning (index 0)
            handleDanceData(0, danceData);
        })
        .catch(error => {
            console.error('Error fetching JSON:', error);
        });
}