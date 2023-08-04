const btnLeave = document.querySelector(".btn-leave-room");
const achievements = document.querySelectorAll(".achieve-items__positions");
const inProgressType = document.getElementById("in-progress");
const completedType = document.getElementById("completed");
btnLeave.addEventListener("click", function () {
   window.location.href = "join"
});

inProgressType.addEventListener("click", function (){
   inProgressType.classList.add("type-selected");
   completedType.classList.remove("type-selected");
   checkAchievements();
});

completedType.addEventListener("click", function (){
   completedType.classList.add("type-selected");
   inProgressType.classList.remove("type-selected");
   checkAchievements();
});


function checkAchievements() {
   achievements.forEach(element => {
      if (inProgressType.classList.contains("type-selected")) {
         element.classList.remove("hidden");
         if (element.dataset.complete === "1") {
            element.classList.add("achievements__positions__completed");
         }
         if (element.dataset.collect === "1") {
            element.classList.add("hidden");
         }
      }
      if (completedType.classList.contains("type-selected")) {
         element.classList.remove("hidden")
         if ((element.dataset.complete === "1") || (element.dataset.complete === "0")) {
            element.classList.add("hidden");
         }
         if (element.dataset.collect === "1") {
            element.classList.add("achievements__positions__collected");
         }
      }
   });
}

checkAchievements();

