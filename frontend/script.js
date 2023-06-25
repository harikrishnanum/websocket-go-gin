var input = document.getElementById("input");
var output = document.getElementById("output");
var socket = new WebSocket("ws://localhost:8448/todo");

socket.onopen = function () {
    output.innerHTML += "Status: Connected\n";
};

socket.onmessage = function (e) {
    output.innerHTML = "\nServer: " + e.data + "\n";
};

socket.onclose = function (event) {
    if (event.wasClean) {
        output.innerHTML += "\nStatus: Connection closed cleanly\n";
    } else {
        output.innerHTML += "\nStatus: Connection died\n";
    }
    output.innerHTML += "\nStatus Code: " + event.code + ", Reason: " + event.reason + "\n";
};
function send() {
    socket.send(input.value);
    input.value = "";
}
function handleKeyDown(event) {
    if (event.keyCode === 13) { // Enter key code is 13
        event.preventDefault(); // Prevent form submission or new line insertion
        send();
    }
}