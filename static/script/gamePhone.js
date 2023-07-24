const btnGo = document.getElementById("btn-join");
const enterInRoom = document.querySelector(".entrance-id-room__field");
const danceField = document.querySelector(".dance-block");
const entranceField = document.querySelector(".entrance");
const warningID = document.getElementById("id-warning");
const emptyID = document.getElementById("id-empty");
const fullID = document.getElementById("id-full");
const btnLogOut = document.getElementById("logout");
const user = document.querySelector(".users");
const menu = document.querySelector(".menu")

function setJsonCookie(name, value, expirationDays) {
    const jsonValue = JSON.stringify(value);
    const encodedValue = encodeURIComponent(jsonValue);
    document.cookie = `${name}=${encodedValue}; path=/; expires=${getExpirationDate(expirationDays)}`;
}

function getExpirationDate(expirationDays) {
    const date = new Date();
    date.setDate(date.getDate() + expirationDays);
    return date.toUTCString();
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

const userInfo = getJsonCookie("userInfoCookie");
let userID = userInfo.UserID;

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
                danceField.classList.remove("hidden");
                emptyID.classList.add("hidden");
                fullID.classList.add("hidden");
                warningID.classList.add("hidden");
                console.log("Connected to the room!");
                joinRoom(userInfo.UserID)
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

let socket;
function joinRoom(userID) {
    socket = new WebSocket("wss://" + window.location.hostname + "/ws/joinToRoom/" + userID);
    socket.onopen = function(event) {
        console.log("WebSocket connection established.");
    };

    socket.onmessage = function(event) {
        let maxTheory = 5600;
        let maxPractice = maxTheory - maxTheory * 0.2; //4480
        let value = 0
        let receivedData = event.data;
        let receivedJSON = JSON.parse(receivedData);
        if ("point" in receivedJSON) {
            let score = receivedJSON["point"];
            console.log("score: " + score);
            value += score;
            console.log("value: " + value);
            let percentage = (value / maxPractice);
            console.log("percentage: " + percentage);
            let pix = 250 * percentage;
            console.log("pix: " + pix);
            scale.style.height = pix + 'px';
        }
        else {
            handleDanceData(receivedJSON);
            document.querySelector('.dance-block__connection').innerText = 'Dance!';
        }
    };

    socket.onclose = function(event) {
        console.log("WebSocket connection closed.");
    };
}


// window.onbeforeunload = exitFromGame;
async function exitFromGame() {
    const response = await fetch("/api/exitFromGame", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"userID": userID}),
    });
    if (!response.ok) {
        console.log('Не удалось выйти из игры');
    } else {
        if (socket !== undefined) {
            socket.close();
            socket = undefined
        }
        console.log('Вышел из игры');
        entranceField.classList.remove("hidden");
        danceField.classList.add("hidden");
    }
    stop = 1;
}

async function exitFromAccount() {
    const response = await fetch("/api/exitFromAccount", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({"userID": userID}),
    });
    if (!response.ok) {
        if (socket !== undefined) {
            socket.close();
            socket = undefined
        }
        console.log('Не удалось выйти из аккаунта');
    } else {
        console.log('Вышел из аккаунта');
        window.location.href = '/logIn';
    }
}

let isOpen = false;
function userMenu() {
    if (!isOpen) {
        menu.classList.remove("menu-hidden");
        menu.classList.add("menu-open");
        isOpen = true;
    }
    else {
        menu.classList.add("menu-hidden");
        menu.classList.remove("menu-open");
        isOpen = false;
    }
}


btnLogOut.addEventListener("click", exitFromAccount)
btnGo.addEventListener("click", sendMessage);
user.addEventListener("click", userMenu);


if (window.DeviceMotionEvent && window.DeviceOrientationEvent) {
    const sensorFrequency = 62.5;
    const interval = 1000 / sensorFrequency;

    let sensorData = [];

    function handleSensorData(event) {
        try {
            const { alpha, beta, gamma } = event.rotationRate;
            const { x, y, z } = event.accelerationIncludingGravity;
            sensorData.push({
                alpha,
                beta,
                gamma,
                x,
                y,
                z,
            });
        } catch (error) {
            //document.writeln('Ошибка при обработке данных сенсора:', error);
        }
    }

    function startRecording(name, duration) {
        window.addEventListener('devicemotion', handleSensorData, true);
        window.addEventListener('deviceorientation', handleSensorData, true);
        setTimeout(function () {stopRecording(name);}, duration * 1300); // 1000 /// // / / /// / / / / / / / / / / /  / /              /////
    }

    function stopRecording(name) {
        window.removeEventListener('devicemotion', handleSensorData, true);
        window.removeEventListener('deviceorientation', handleSensorData, true);
        let pointCount = 12; // 1146
        let outputString = '0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000, 0.0000';

        for (let i = 0; i < sensorData.length; i++) {
            const {x, y, z, alpha, beta, gamma} = sensorData[i];
            outputString += `, ${MyRound10(x)}, ${MyRound10(y)}, ${MyRound10(z)}, ${MyRound10(alpha)}, ${MyRound10(beta)}, ${MyRound10(gamma)}`;
            pointCount += 6;
        }
        for (pointCount; pointCount < 1146; pointCount++){
            outputString += ', 0.0000';
        }
        sensorData = [];
        //document.write(outputString);

        let data = JSON.stringify({"name": name, "motionString": outputString, "userID": userID});
        //document.writeln(data);
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

let stop = 0;

function sendDataToServer(data) {
    // Replace the URL with the appropriate endpoint to handle the data on your server
    let url = '/api/motion';
    if (stop !== 1) {
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: data
        })
            .then(function (response) {
                if (response.ok) {
                    console.log('Данные успешно отправлены.');
                } else {
                    console.log('Ошибка при отправке данных. Статус:', response.status);
                    if (response.status === 409) {
                        stop = 1;
                        document.querySelector('.dance-block__connection').innerText = 'Комната была закрыта';
                        //window.location.replace("/join")
                        return;
                    }
                }
            })
            .catch(function (error) {
                console.log('Ошибка при отправке данных:', error);
            });
    } else {
        console.log('123');
    }
}

function handleDanceData(danceDataJson) {
    //let oldStartTime = 0;
    for (let danceData of danceDataJson) {
        setTimeout(function () {
                if (stop !== 1) {
                    startRecording(danceData['name'], danceData['duration']);
                }
            },
            (danceData['start_time']) * 1000);
        //oldStartTime += danceData['start_time'];
    }
}

const scale = document.querySelector('.dance-block__ratingScale');

function getRatingScale() {
    let maxTheory = 5600;
    let maxPractice = maxTheory - maxTheory * 0.2; //4480
    let value = 0;
}

function anim(score) {
    console.log("score: " + score);
    if (value > 5600) return
    value += score;
    console.log("value: " + value);
    const percentage = (value / maxPractice);
    console.log("percentage: " + percentage);
    let pix = 250 * percentage;
    console.log("pix: " + pix);
    scale.style.height = pix + 'px';// Рекурсивно вызываем функцию для создания плавной анимации
}
