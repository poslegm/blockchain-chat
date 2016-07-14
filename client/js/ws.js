var wsUrl = "ws://" + location.host + "/websocket";
var socket = new WebSocket(wsUrl);

var PUBLIC_KEY = "";
var dialogs = {}
var lazyMessages = [];

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
            handleMessages(message['Messages']);
        } else {
            console.log("PUSHED");
            lazyMessages = lazyMessages.concat(message['Messages']);
        }
    } else if (message['Type'] == 'Key') {
        console.log("KEY");
        handlePublicKey(message['Key']);
        handleMessages(lazyMessages);
        lazyMessages = [];
    } else if (message['Type'] == 'NewMessage') {
        console.log("NEW MESSAGE: " + message['']);
        handleMessages(message['Messages']);
        addNewMessagesToViews(message['Messages']);
    }
};

function handleMessages(messages) {
    messages.forEach(function (o) {
        if (o['Sender'] === PUBLIC_KEY || o['Receiver'] === PUBLIC_KEY) {
            if (o['Sender'] === PUBLIC_KEY) {
                dictAppend(dialogs, o['Receiver'], o);
            } else {
                dictAppend(dialogs, o['Sender'], o);
            }
        }
    });

    viewDialogs();
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

function viewDialogs() {
    for (var user in dialogs) {
        console.log(user);
        listElem = $("<li></li>")
            .attr('id', "dialog-" + user)
            .click(function() {
                changeDialog(dialogs, user);
                $("#dialog-" + user).find("div.user").attr("style", "");
            })
            .append($("<div></div>")
                .addClass("info")
                .append($("<div></div>")
                    .addClass("user")
                    .text(user)
                )
            );
        $(".list-friends").append(listElem)
    }
}

function changeDialog(dialogs, userKey) {
    return function () {
        $("#top-name").text(userKey);
        $(".messages").empty();
        for (var i in dialogs[userKey]) {
            msg = dialogs[userKey][i];
            if (msg["Sender"] === PUBLIC_KEY) {
                appendMessage(true, "ME: " + PUBLIC_KEY, msg["Text"]);
            } else {
                appendMessage(false, msg["Sender"], msg["Text"]);
            }
            console.log(dialogs[userKey][i]);
        }
    }
}

function appendMessage(my, sender, text) {
    var messageClass = "i";
    if (my) {
        messageClass = "friend-with-a-SVAGina";
    }

    $(".messages").append($("<li></li>")
        .attr("id", sender)
        .addClass(messageClass)
        .append($("<div></div>")
            .addClass("head")
//            .append($("<span></span>")
//                .addClass("time")
//                .text("10:13, 10.06.2016")
            .append($("<span></span>")
                .addClass("name")
                .text(sender)))
        .append($("<div></div>")
            .addClass("message")
            .text(text)));
}

function addNewMessagesToViews(messages) {
    messages.forEach(function (o) {
        if (o['Sender'] === PUBLIC_KEY || o['Receiver'] === PUBLIC_KEY) {
            if (o['Sender'] === PUBLIC_KEY) {
                addNewMessage(true, o['Receiver'], o);
            } else {
                addNewMessage(false, o['Sender'], o);
            }
        }
    });
}

function addNewMessage(my, user, message) {
    if ($("#top-name").val() === user) {
        appendMessage(my, user, message['Text'])
    } else {
        $("#dialog-" + user).find("div.user").attr("style", "font-weight:bold");
    }
}