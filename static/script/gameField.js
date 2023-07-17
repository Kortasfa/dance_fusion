function getUsersByCookie() {
    let numberOfUser = 0;
    let allCookies = document.cookie;
    let cookiesArray = allCookies.split(';');
    for (let i = 0; i < cookiesArray.length; i++) {
        let cookie = cookiesArray[i].split('=')
        let name = cookie[0];
        let findUser = name.indexOf('User');
        if (findUser === 1) {
            let parts = cookie[1].split(',');
            let userName = parts[1];
            let imgSrc = parts[2];
            numberOfUser = numberOfUser + 1;
            let indexUser = document.getElementById('hero' + numberOfUser);
            let indexUserName = document.getElementById('heroName' + numberOfUser);
            let indexUserImg = document.getElementById('heroImg' + numberOfUser);
            indexUser.classList.remove('hidden');
            indexUserImg.src =  '../' + imgSrc;
            indexUserName.innerText = userName;
            redyPlayer = true;
            ChengeBtn();
        }
    }
}
getUsersByCookie();