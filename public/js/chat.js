
function onclickSendMessage() {
    
    let message = document.getElementById("sendBar").value;
    
    if (message === "") {
        return;
    }

    let sectionChat = document.getElementById("chat");
    sectionChat.innerHTML += '<div class="userText"><p>USER :</p> <p>' + message + '</p><p class="hour">'+ getcurrentDate() +'</p><button><img src="./public/img/edition.png" alt="image d\'un crayon qui modifie la réponse" width="30" height="30"><p>Modifier le message</p></button></div>';
    document.getElementById("sendBar").value = "";

    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });


    // Délégation d'événements pour gérer les clics sur les boutons
    document.getElementById('chat').addEventListener('click', function(event) {
        if (event.target.tagName === 'BUTTON' || event.target.closest('button')) {
            const button = event.target.tagName === 'BUTTON' ? event.target : event.target.closest('button');
            const userTextDiv = button.closest('.userText');
            const messageContent = userTextDiv.querySelector('p:nth-child(2)').textContent;

            let sendBar = document.getElementById("sendBar");

            sendBar.value = messageContent;
        } 
    });


    sendMessageToServer(message);

    createNewZoneBotMessage();

}

function getcurrentDate(){

    const currentDate = new Date();

    const day = String(currentDate.getDate()).padStart(2, '0'); 
    const month = String(currentDate.getMonth() + 1).padStart(2, '0'); 
    const year = currentDate.getFullYear(); 
    const hours = String(currentDate.getHours()).padStart(2, '0');
    const minutes = String(currentDate.getMinutes()).padStart(2, '0'); 

    const date = `${day}/${month}/${year} ${hours}:${minutes}`;

    return date;
}


const socket = new WebSocket('ws://localhost:3333/ws');

socket.onopen = function(event) {
    console.log('WebSocket connection established');
};


const charFin = "#fin#";

socket.onmessage = function(event) {
    console.log('Message from server: ', event.data);

    let sectionChat = document.getElementById("chat");
    let newDiv = document.getElementById("new");
    let pAnswer = document.getElementById("answer");

    if(isGIF == true){
        removeWaitingGIF();
    }

    if(document.getElementById("new") != null){

        //si le message recu n'est pas le char de fin : 
        if (event.data != charFin) {
            pAnswer.innerHTML += event.data;
        } else {
            let pHour = document.getElementById("newhour");
            pHour.innerHTML = getcurrentDate();
            pHour.removeAttribute("id");
            newDiv.removeAttribute("id");
            pAnswer.removeAttribute("id");
            window.scrollTo({
                top: document.body.scrollHeight,
                behavior: 'smooth'
            });
        }
    }
};


socket.onerror = function(error) {
    console.log('WebSocket Error: ', error);
};


socket.onclose = function(event) {
    console.log('WebSocket connection closed: ', event);
};

function sendMessageToServer(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(message);
        console.log('Message sent to server: ', message);
    } else {
        console.log('WebSocket connection is not open');
    }
}


function createNewZoneBotMessage(){
    let sectionChat = document.getElementById("chat");
    sectionChat.innerHTML += '<div class="botText" id="new"><p>BOT :</p> <p id="answer">' + '</p><p class="hour" id="newhour">'+ "" +'</p></div>';
    addWaitingGIF();
}


let isGIF;

function addWaitingGIF(){
    const gifPath = "./public/img/waiting.gif";

    const imgGIF = document.createElement("img");
    imgGIF.src = gifPath;

    const gifDiv = document.getElementById("answer");
    gifDiv.appendChild(imgGIF);

    isGIF = true;
}

function removeWaitingGIF(){
    const gifDiv = document.getElementById("answer");
    gifDiv.removeChild(gifDiv.firstChild);

    isGIF = false;
}