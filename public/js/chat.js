
        /* Initialisation de la page */

function init() {
    const sendBar = document.getElementById("sendBar");
    sendBar.addEventListener("input", onCheckArea);

    onCheckArea();

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
}

        /* redirection sur la déconnexion */

function redirect_disconnect() {
    document.location.pathname = "/disconnect"
}

function data_deletion() {
    window.location.href = window.location.origin + "/disconnect?delete=true"
}

        /* Tuto */

const dialogs = [
    ['Pour écrire un message cliquez dans la boite surligné en rouge, puis tapez votre message', document.getElementById("sendBar"), 'Écrire un message'],
    ['Pour envoyer un message à l\'agent cliquez ensuite sur ce bouton après avoir rédigé votre message', document.getElementById("sendButton"), 'Envoyer un message'],
    ['Quand vous avez fini de dialoguer avec l\'agent vous pouvez vous déconnecter en appuyant sur ce bouton', document.getElementById("disconnect"), 'Déconnexion'],
    ['Si vous souhaitez supprimer vos données de l\'application cliquez ici', document.getElementById("data_delete"), 'Suppression des données'],
]
var current_dialog = 0

function start_tutorial () {
    const dialog = document.querySelector("dialog");
    dialog.showModal()
    current_dialog = -1
    next_dialog()
}

function next_dialog () {
    current_dialog++
    if (current_dialog > 0) {
        dialogs[current_dialog - 1][1].style.border = ''
    }
    if (current_dialog >= dialogs.length) {
        document.querySelector("dialog").close()
    }
    document.getElementById("tuto_title").innerHTML = dialogs[current_dialog][2]
    document.getElementById("tuto_text").innerHTML = dialogs[current_dialog][0]
    dialogs[current_dialog][1].style.border = '5px solid red'
}

function close_dialog () {
    if (current_dialog >= 0) {
        dialogs[current_dialog][1].style.border = ''
    }
    document.querySelector('dialog').close()
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
    console.log('Message from server: ', event.data);

    if (historyDisplayed == false /*&& newUser == false*/) {
        history = JSON.parse(event.data);
        displayHistory();
        historyDisplayed = true;
    }

    let newDiv = document.getElementById("new");
    let pAnswer = document.getElementById("answer");

    if (isGIF == true) {
        removeWaitingGIF();
    }

    if (document.getElementById("new") != null) {

        //si le message recu n'est pas le char de fin : 
        if (event.data != charFin) {
            let message = displayMessage(event.data);
            pAnswer.innerHTML += message;

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


        /* Gestion de l'affichage des messages */

function displayMessage(message) {
    if (message == "*") {
        return " - ";
    } else if (message == "\n") {
        return "</br>";
    } else {

        let ret = message;

        if (message.length > 0 && message[message.length - 1] == "\n") {
            ret += "</br>";
        }
        return ret;
    }
}


function processMessage(message) {

    let words = message.split(' ');

    let modifiedWords = words.map(word => displayMessage(word));

    return modifiedWords.join(' ');
}


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


    /* Affichage et gestion historique */

let historyDisplayed = false;
let history = [];

function displayHistory() {

    const chatContainer = document.getElementById('chat');

    //gestion de tous les messages
    history.forEach(message => {

        // Créez une div pour chaque message
        const messageDiv = document.createElement('div');

        myMessage = processMessage(message.content);
         
        let destinataire  = message.destinataire;
         
        if (destinataire === 'assistant') {
            messageDiv.className = 'userText';
            messageDiv.innerHTML = '<p>USER :</p> <p>' + myMessage + '</p><p class="hour">' + message.date + '</p><button><img src="./public/img/edition.png" alt="image d\'un crayon qui modifie la réponse" width="30" height="30"><p>Modifier le message</p></button>';
        } else {
            messageDiv.className = 'botText';
            messageDiv.innerHTML = '<p>BOT :</p> <p>' + myMessage + '</p><p class="hour">' + message.date + '</p>';
        }

        chatContainer.appendChild(messageDiv);
        
    });

    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });

}