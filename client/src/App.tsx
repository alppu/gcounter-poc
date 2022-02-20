import { useState } from "react";
import "./App.css";

function App() {
  const [message, setMessage] = useState<string>("");
  const [messages, setMessages] = useState<
    Array<{
      Message: string;
    }>
  >([]);

  return (
    <div className="App">
      <label>Send message:</label>
      <input onChange={(e) => setMessage(e.target.value)}></input>
      <button
        onClick={async () => {
          const response = await fetch("http://localhost:8080/message", {
            method: "POST",
            body: JSON.stringify({
              message,
            }),
          });
        }}
      >
        Send
      </button>

      <button
        onClick={async () => {
          const response = await fetch("http://localhost:8080/messages");
          const messages = await response.json();
          setMessages(messages);
        }}
      >
        Get messages
      </button>

      <ul>
        {messages.map((message) => (
          <li>{message.Message}</li>
        ))}
      </ul>
    </div>
  );
}

export default App;
