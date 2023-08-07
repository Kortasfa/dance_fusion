function scoreGrade(res, score, moveName, songName) {
    switch (songName) {
        case "Forget You":
            return scoreGradeForgetYou(res, score, moveName);
        case "Rasputin":
            return scoreGradeRasputin(res, score, moveName);
        case "There is nothing better":
            return scoreGradeThereIsNothingBetter(res, score, moveName);
        default:
            return scoreGradeForgetYou(res, score, moveName);
    }
}
function scoreGradeForgetYou(res, score, moveName) {
    for (let motionDict of res["results"]) {
            if (motionDict["label"] === moveName) {
                score = motionDict["value"];
                if (moveName === "up-down-click") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    let sameMoveTwo = res["results"].find(item => item.label === "up-down-hands");
                    score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.2;
                    }
                }
                if (moveName === "up-down-hands") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                    if (score > 0.001 && score < 0.1) {
                        score = 0.2;
                    }
                }
                if (moveName === "collapsing-breeding-locks") {
                    if (score > 0.001 && score < 0.005) {
                        score = 0.2
                    }
                    if (score >= 0.005 && score < 0.01) {
                        score = 0.3;
                    }
                    if (score >= 0.01) {
                        score = 0.4;
                    }

                }
                if (moveName === "going-to-the-side-with-clapping") {
                    if (score > 0.001 && score < 0.05) {
                        score = 0.05;
                    }
                }
                if (moveName === "hands-up-and-down") {
                    let sameMoveOne = res["results"].find(item => item.label === "up-down-click");
                    let sameMoveTwo = res["results"].find(item => item.label === "up-down-hands");
                    score = Math.max(motionDict["value"], sameMoveOne.value - 0.1, sameMoveTwo.value - 0.1);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.2;
                    }
                }
                if (moveName === "side-hit") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                    if (score >= 0.05 & score < 0.1) {
                        score = 0.3;
                    }
                    if (score >= 0.1) {
                        score = 0.4;

                    }
                }
                if (moveName === "down-from-the-middle-in-arc") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
                }
                if (moveName === "collapsing-breeding-locks") {
                    let sameMoveOne = res["results"].find(item => item.label === "up-down-hands");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                }
                if (moveName === "draw-circle") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    let sameMoveTwo = res["results"].find(item => item.label === "up-down-hands");
                    score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
                }
                if (moveName === "i-am-strong") {
                    let sameMoveOne = res["results"].find(item => item.label === "gold-turn");
                    let sameMoveTwo = res["results"].find(item => item.label === "gold-circle");
                    let sameMoveThree = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value, sameMoveThree.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
                }
                if (moveName === "broad-back") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                }
                if (moveName === "like-a-chicken") {
                    let sameMoveOne = res["results"].find(item => item.label === "hands-up-and-down");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
                }
                if (moveName === "gold-circle") {
                    let sameMoveOne = res["results"].find(item => item.label === "up-down-click");
                    score = Math.max(motionDict["value"], sameMoveOne.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
                }
                if (moveName === "gold-turn") {
                    let sameMoveOne = res["results"].find(item => item.label === "down-from-the-middle-in-arc");
                    let sameMoveTwo = res["results"].find(item => item.label === "gold-circle");
                    score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value);
                    if (score > 0.01 && score < 0.1) {
                        score = 0.1;
                    }
            }
        }

    }
    return score;
}
function scoreGradeRasputin(res, score, moveName) {
    let sameMoveList = ["i-am-fly", "russian-move", "kalinka-move", "look-at-my-boots", "good-mood"]
    for (let motionDict of res["results"]) {
        if (motionDict["label"] === moveName) {
            score = motionDict["value"];
            if (!sameMoveList.includes(moveName)) {
                if (motionDict["russian-move"] > 0.8) {
                    return 0
                }

                if (motionDict["i-am-fly"] > 0.8) {
                    return 0
                }
            }
            if (moveName === "i-am-fly") {
                let sameMoveOne = res["results"].find(item => item.label === "russian-move");
                let sameMoveTwo = res["results"].find(item => item.label === "clap-clap");
                score = Math.max(motionDict["value"], sameMoveOne.value ,sameMoveTwo.value);
            }
            if (moveName === "russian-move") {
                let sameMoveOne = res["results"].find(item => item.label === "kalinka-move");
                let sameMoveTwo = res["results"].find(item => item.label === "clap-clap");
                score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value);
            }
            if (moveName === "kalinka-move") {
                let sameMoveOne = res["results"].find(item => item.label === "russian-move");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (moveName === "guitar-move") {
                let sameMoveOne = res["results"].find(item => item.label === "kalinka-move");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (moveName === "up-down-move") {
                let sameMoveOne = res["results"].find(item => item.label === "clap-clap");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (moveName === "big-and-strong") {
                score = Math.max(motionDict["value"]);
            }
            if (moveName === "look-at-my-boots") {
                let sameMoveOne = res["results"].find(item => item.label === "russian-move");
                let sameMoveTwo = res["results"].find(item => item.label === "i-am-fly");
                score = Math.max(motionDict["value"], sameMoveOne.value, sameMoveTwo.value);
            }
            if (moveName === "sweeping") {
                let sameMoveOne = res["results"].find(item => item.label === "kalinka-move");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (moveName === "up-down-move") {
                let sameMoveOne = res["results"].find(item => item.label === "kalinka-move");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (moveName === "good-mood") {
                let sameMoveOne = res["results"].find(item => item.label === "russian-move");
                score = Math.max(motionDict["value"], sameMoveOne.value);
            }
            if (score < 0.05)
            {
                score = score * 10;
            }
        }
    }
    return score;
}
function scoreGradeThereIsNothingBetter(res, score, moveName) {
    for (let motionDict of res["results"]) {
        if (motionDict["label"] === moveName) {
            score = motionDict["value"];
        }
    }
    return score;
}
window.scoreGrade = scoreGrade;