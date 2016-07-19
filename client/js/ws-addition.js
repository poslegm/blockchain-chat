var wsUrl = "ws://" + location.host + "/websocket-addition";
var socket = new WebSocket(wsUrl);

$(document).ready(function() {
    $('#send-contact').click(function () {
        sendMessageContact();
    });

    $('#send-keys').click(function () {
        sendMessageKeys();
    });
});

function sendMessageContact() {
    key = $('#publicKeyContact').val();
    if (key === "") {
        return
    }

    msg = {
        Type: "Contact",
        Key: key
    }
    socket.send(JSON.stringify(msg));
}

function sendMessageKeys() {
    publicKey = $('#publicKey').val();
    privateKey = $('#privateKey').val();
    passphrase = $('#passphrase').val();
    if (publicKey === "" || privateKey === "") {
        return
    }

    msg = {
        Type: "KeyPair",
        PublicKey: publicKey,
        PrivateKey: privateKey,
        Passphrase: passphrase
    }
    socket.send(JSON.stringify(msg));
}

socket.onmessage = function(event) {
    var message = JSON.parse(event.data);

    if (message['Type'] === 'OkContact') {
        alert("Контакт добавлен, хеш: " + message['Key']);
    } else if (message['Type'] === 'BadContact') {
        alert("Не удалось добавить контакт");
    } else if (message['Type'] === 'OkKeyPair') {
        alert("Пара ключей добавлена");
    } else if (message['Type'] === 'BadKeyPair') {
        alert("Не удалось добавить пару ключей");
    }
};