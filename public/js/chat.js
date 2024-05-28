
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


socket.onmessage = function(event) {
    console.log('Message from server: ', event.data);
    let sectionChat = document.getElementById("chat");
    sectionChat.innerHTML += '<div class="botText"><p>BOT :</p> <p>' + event.data + '</p><p class="hour">'+ getcurrentDate() +'</p></div>';
    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });
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
