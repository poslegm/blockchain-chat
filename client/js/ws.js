var wsUrl = "ws://" + location.host + "/websocket";
socket = new WebSocket(wsUrl);

socket.onopen = function() {
    socket.send(JSON.stringify({Command: "do it", Value: "kek"}));
    socket.send(JSON.stringify({Command: "do it"}));
};

socket.onmessage = function(event) {
    console.log(event.data);
};