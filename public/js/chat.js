
        /* Initialisation de la page */

function init() {
    const sendBar = document.getElementById("sendBar");
    sendBar.addEventListener("input", onCheckArea);

    onCheckArea();
}

        /* redirection sur la déconnexion */

function redirect_disconnect() {
    document.location.pathname = "/disconnect"
}

        /* Gestion de l'envoi de messages */

function onclickSendMessage() {

    let message = document.getElementById("sendBar").value;

    if (message === "") {
        return;
    }

    let sectionChat = document.getElementById("chat");
    let current_date = getcurrentDate();
    sectionChat.innerHTML += '<div class="userText"><p>USER :</p> <p>' + message + '</p><p class="hour">' + current_date + '</p><button><img src="./public/img/edition.png" alt="image d\'un crayon qui modifie la réponse" width="30" height="30"><p>Modifier le message</p></button></div>';
    document.getElementById("sendBar").value = "";

    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });


    // Délégation d'événements pour gérer les clics sur les boutons
    document.getElementById('chat').addEventListener('click', function (event) {
        if (event.target.tagName === 'BUTTON' || event.target.closest('button')) {
            const button = event.target.tagName === 'BUTTON' ? event.target : event.target.closest('button');
            const userTextDiv = button.closest('.userText');
            const messageContent = userTextDiv.querySelector('p:nth-child(2)').textContent;

            let sendBar = document.getElementById("sendBar");

            sendBar.value = messageContent;
        }
    });

    sendMessageToServer(message, current_date);

    createNewZoneBotMessage();

    disableSendButton();
}

        /* Envoi de messages au serveur */

let botanswer = false;

function sendMessageToServer(message, date) {

    botanswer = true;
    if (socket.readyState === WebSocket.OPEN) {
        const data = JSON.stringify({ message: message, date: date });
        socket.send(data);
        console.log('Message sent to server: ', data);
    } else {
        console.log('WebSocket connection is not open');
    }
}


function createNewZoneBotMessage() {
    let sectionChat = document.getElementById("chat");
    sectionChat.innerHTML += '<div class="botText" id="new"><p>BOT :</p> <p id="answer">' + '</p><p class="hour" id="newhour">' + "" + '</p></div>';

    addWaitingGIF();

    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });
}

function getcurrentDate() {

    const currentDate = new Date();

    const day = String(currentDate.getDate()).padStart(2, '0');
    const month = String(currentDate.getMonth() + 1).padStart(2, '0');
    const year = currentDate.getFullYear();
    const hours = String(currentDate.getHours()).padStart(2, '0');
    const minutes = String(currentDate.getMinutes()).padStart(2, '0');

    const date = `${day}/${month}/${year} ${hours}:${minutes}`;

    return date;
}

        /* Gestion du bouton Envoyer */

function disableSendButton() {
    let sendButton = document.getElementById("sendButton");
    sendButton.disabled = true;
    sendButton.style.cursor = "not-allowed";
}

function enableSendButton() {
    let sendButton = document.getElementById("sendButton");
    sendButton.disabled = false;
    sendButton.style.cursor = "pointer";
}


function onCheckArea(event) {
    const sendBar = document.getElementById("sendBar");
    if (sendBar.value.length > 0 && sendBar.value.trim().length > 0 && botanswer == false){
        enableSendButton();
    } else {
        disableSendButton();
    }
}

        /* Gestion des WebSockets */

const socket = new WebSocket('ws://localhost:3333/ws');
const charFin = "#fin#";

socket.onopen = function (event) {
    console.log('WebSocket connection established');
};


socket.onmessage = function (event) {
    //console.log('Message from server: ', event.data);

    let newDiv = document.getElementById("new");
    let pAnswer = document.getElementById("answer");

    if (isGIF == true) {
        removeWaitingGIF();
    }

    if (document.getElementById("new") != null) {

        //si le message recu n'est pas le char de fin : 
        if (event.data != charFin) {
            if (event.data == "*") {
                pAnswer.innerHTML += " - ";
            } else if (event.data == "\n") {
                pAnswer.innerHTML += "</br>";
            } else {
                pAnswer.innerHTML += event.data;

                if (event.data.length > 0 && event.data[event.data.length - 1] == "\n") {
                    pAnswer.innerHTML += "</br>";
                }
            }
        } else {
            let pHour = document.getElementById("newhour");
            pHour.innerHTML = getcurrentDate();
            pHour.removeAttribute("id");
            newDiv.removeAttribute("id");
            pAnswer.removeAttribute("id");

            botanswer = false;
            onCheckArea();
        }

        window.scrollTo({
            top: document.body.scrollHeight,
            behavior: 'smooth'
        });
    }

};


socket.onerror = function (error) {
    console.log('WebSocket Error: ', error);
};


socket.onclose = function (event) {
    console.log('WebSocket connection closed: ', event);
};



        /* Gestion des GIFs */

let isGIF;

function addWaitingGIF() {
    const gifPath = "./public/img/waiting.gif";

    const imgGIF = document.createElement("img");
    imgGIF.src = gifPath;

    const gifDiv = document.getElementById("answer");
    gifDiv.appendChild(imgGIF);

    isGIF = true;
}

function removeWaitingGIF() {
    const gifDiv = document.getElementById("answer");
    gifDiv.removeChild(gifDiv.firstChild);

    isGIF = false;
}