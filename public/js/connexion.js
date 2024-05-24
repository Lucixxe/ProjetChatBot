
function changePageToRegister(){
    document.getElementById("title").innerHTML = "Inscription au ChatBot";
    document.getElementById("validate").setAttribute("value", "S'inscrire");
    document.getElementById("register").innerHTML = "Retour à la page de connexion";
    document.title = "ChatBot - Inscription";
    document.getElementById("register").setAttribute("onclick", "changePageToLogin()");
    hideErrorMessage();
}

function changePageToLogin(){
    document.getElementById("title").innerHTML = "Connexion au ChatBot";
    document.getElementById("validate").setAttribute("value", "Se connecter");
    document.getElementById("register").innerHTML = "S'inscrire";
    document.title = "ChatBot - Connexion";
    document.getElementById("register").setAttribute("onclick", "changePageToRegister()");
    hideErrorMessage();
}


function showErrorMessage(){
    document.getElementById("error").removeAttribute("hidden");
}

function hideErrorMessage(){    
    document.getElementById("error").setAttribute("hidden", "true");
}


