const btnGo = document.getElementById("btn-go")

async function startGame(){
    const start = {
       button: true,
    }
    let startJson = JSON.stringify(start);
    console.log(startJson);

    const response = await fetch("/api/start", {
       method: "POST",
       headers: {
           'Content-Type': 'application/json;charset=utf-8'
       },
       body: startJson,
    });

}


btnGo.addEventListener("click", startGame);




