import consumer from "./consumer"


window.onload = function() {
  var chatForm = document.getElementById("chat-form")

  if (!chatForm) {
    return
  }

  var messages = document.getElementById("messages")

  var roomId = chatForm.dataset.roomId

  var wsConn = consumer.subscriptions.create({ channel: "RoomChannel", id: roomId }, {
    received(data) {
      // Called when there's incoming data on the websocket for this channel
      if (data["message"]) {
        messages.insertAdjacentHTML("beforeend", this.createLine(data, ""))
      }

      if (data["user_join"]) {
        data["send_by"] = "SYSTEM"
        data["message"] = data["user_join"]
        messages.insertAdjacentHTML("beforeend", this.createLine(data, "is-success"))
      }

      if (data["user_leave"]) {
        data["send_by"] = "SYSTEM"
        data["message"] = data["user_leave"]
        messages.insertAdjacentHTML("beforeend", this.createLine(data, "is-warning"))
      }
    },

    createLine(data, messageType) {
      return '<article class="message ' + messageType + '"><div class="message-header">' +
          data["send_by"] + ' Says' +
        '</div><div class="message-body">' +
          data["message"] + "</div></article>"
    },

    connected() {
      var data = {
        "send_by": "SYSTEM",
        "message": "RoomChannel is connected."
      }
      if (messages) {
        messages.insertAdjacentHTML("beforeend", this.createLine(data, "is-info"))
      }
    },

    disconnected() {
      var data = {
        "send_by": "SYSTEM",
        "message": "Websocket connection is closed by the server."
      }
      if (messages) {
        messages.insertAdjacentHTML("beforeend", this.createLine(data, "is-danger"))
      }
    },

    rejected() {
      var data = {
        "send_by": "SYSTEM",
        "message": "The subscription is rejected by the server."
      }
      if (messages) {
        messages.insertAdjacentHTML("beforeend", this.createLine(data, "is-warning"))
      }

      if (message) {
        message.setAttribute("disabled", true)
        message.setAttribute("placeholder", "The subscription is rejected by the server.")
      }

      document.querySelectorAll(".button.is-dark").forEach((btn) => btn.setAttribute("disabled", true))
    },

    stopStream() {
      this.perform("stop_stream")
    }
  })

  if (roomId === "fake") {
    return
  }

  var message = chatForm.querySelector("textarea")

  chatForm.onsubmit = function(e) {
    e.preventDefault()

    wsConn.perform("send_message", {message: message.value})

    message.value = ""
  }

  message.addEventListener("keypress", function(e) {
    if(e.code === "Enter" && !e.shiftKey) {
      e.preventDefault();
      message.closest("form").dispatchEvent(new SubmitEvent("submit"));
    }
  })

  var kickForm = document.getElementById("kick-form")
  var name = kickForm.querySelector("input")
  kickForm.onsubmit = function(e) {
    e.preventDefault()
    wsConn.perform("kick", {name: name.value})

    name.value = ""
  }

  var unsubscribeBtn = document.getElementById("unsubscribe-btn")

  if (unsubscribeBtn) {
    unsubscribeBtn.onclick = function() {
      wsConn.unsubscribe()
    }
  }

  var stopStreamBtn = document.getElementById("stop-stream-btn")

  if (stopStreamBtn) {
    stopStreamBtn.onclick = function() {
      wsConn.stopStream()
    }
  }
}

