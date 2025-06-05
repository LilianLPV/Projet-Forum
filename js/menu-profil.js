const menuProfil = document.getElementById("photo-de-profil");
const menu = document.getElementById("menu-profil");

menuProfil.addEventListener("click", function() {
    console.log("Menu Profil clicked");
    if (menu.style.display === "block") {
        menu.style.display = "none";
    } else {
        menu.style.display = "block";
    }
    console.log("ça marche");
});