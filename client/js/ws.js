var wsUrl = "ws://" + location.host + "/websocket";
var socket = new WebSocket(wsUrl);

var PUBLIC_KEY = ""
var lazyMessages = []

socket.onopen = function() {
    if (PUBLIC_KEY == "") {
        socket.send(JSON.stringify({ Type: "GetMyKey" }));
    }
    socket.send(JSON.stringify({ Type: "GetMessages" }));
};

socket.onclose = function() {
    console.log("WS closed");
    socket = new WebSocket(wsUrl);
}

function sendMessage() {
    text = $('#texxt').val()
    msg = {
        Type: "SendMessage",
        Messages: [{
            Receiver: "123fdsk124pn12",
            Sender: PUBLIC_KEY,
            Text: text
        }]
    }
    socket.send(JSON.stringify(msg))
}

socket.onmessage = function(event) {
    var message = JSON.parse(event.data);

    if (message['Type'] == 'AllMessages') {
        console.log("ALL MESSAGES");
        if (PUBLIC_KEY != "") {
            handleAllMessages(message['Messages']);
        } else {
            console.log("PUSHED");
            lazyMessages = lazyMessages.concat(message['Messages']);
        }
    } else if (message['Type'] == 'Key') {
        console.log("KEY");
        handlePublicKey(message['Key']);
        handleAllMessages(lazyMessages);
        lazyMessages = [];
    }
};

function handleAllMessages(messages) {
    var dialogs = {};
    messages.forEach(function (o) {
        console.log(o['Sender'] + " / " + PUBLIC_KEY + " / " + (o['Sender'] === PUBLIC_KEY))
        if (o['Sender'] === PUBLIC_KEY || o['Receiver'] === PUBLIC_KEY) {
            if (o['Sender'] === PUBLIC_KEY) {
                dictAppend(dialogs, o['Receiver'], o);
            } else {
                dictAppend(dialogs, o['Sender'], o);
            }
        }
    });

    console.log(dialogs);
}

function handlePublicKey(key) {
    PUBLIC_KEY = key;
    console.log(PUBLIC_KEY);
}

function dictAppend(dict, key, value) {
    if (dict[key] !== undefined) {
        dict[key].push(value)
    } else {
        dict[key] = [value]
    }
}