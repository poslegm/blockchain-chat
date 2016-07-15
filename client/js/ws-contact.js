var wsUrl = "ws://" + location.host + "/websocket-addition";
var socket = new WebSocket(wsUrl);

$(document).ready(function() {
    $('#send-contact').click(function () {
        sendMessage();
    });
});

function sendMessage() {
    key = $('#publicKey').val();
    if (key === "") {
        return
    }

    msg = {
        Type: "Key",
        Key: key
    }
    socket.send(JSON.stringify(msg));
}

socket.onmessage = function(event) {
    var message = JSON.parse(event.data);

    if (message['Type'] === 'Ok') {
        alert("Контакт добавлен, хеш: " + message['Key']);
    } else {
        alert("Не удалось добавить контакт");
    }
};