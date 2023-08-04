btnLeave = document.querySelector(".btn-leave-room");

btnLeave.addEventListener("click", function () {
   window.location.href = "join"
});

async function earnPointsForAchievements(achievementID) {
   let response = await fetch("/api/earnPointsForAchievements", {
      method: 'POST',
      headers: {
         'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: `achievement_id=${achievementID}`,
   });
   if (response.ok) {
      console.log('Получил очки за ачивку');
   } else {
      console.log('Не получилось получить баллы за ачивку', response.status);
   }
}