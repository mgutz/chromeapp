/* globals io */

"use strict";

function main() {
    var socket = io();

    socket.on('connect', function() {

        // for some reason emit will send to the server correctly but if the server replies
        // the message is dropped unless we handle it outside of the connect event
        setTimeout(function() {
            socket.emit("echo", "hello world!");
        }, 0);
    });

    socket.on("echo", function(msg) {
        console.log("received message", msg);
    });
}

main();
