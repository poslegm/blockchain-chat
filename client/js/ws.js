var wsUrl = "ws://" + location.host + "/websocket";
socket = new WebSocket(wsUrl);

socket.onopen = function() {
    socket.send(JSON.stringify({ Type: "GetMessages" }));
};

function sendMessage() {
    text = $('#texxt').val()
    msg = {
        Type: "SendMessage",
        Messages: [{
            Receiver: "123fdsk124pn12",
            Sender: "123456",
            Text: text
        }]
    }
    socket.send(JSON.stringify(msg))
}

socket.onmessage = function(event) {
    console.log(event.data);
};