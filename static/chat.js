document.addEventListener("DOMContentLoaded", function() {
    // Add event listener to the form
    document.getElementById("chatForm").addEventListener("submit", function(e) {
        e.preventDefault(); // Prevent default form submission behavior

        const lhid = document.getElementById("lhid").value;
        const chatid = document.getElementById("chatid").value;

        // Fetch chat messages using the provided lhid and chatid
        fetchChatMessages(lhid, chatid);
    });
});

// Function to fetch chat messages from the REST API
function fetchChatMessages(lhid, chatid) {
    fetch(`/api/chat?lhid=${lhid}&chatid=${chatid}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // Expecting a JSON response
        })
        .then(data => {
            renderChatMessages(data);
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
}

// Function to render chat messages into HTML
function renderChatMessages(messages) {
    const chatMessagesDiv = document.getElementById("chat-messages");
    chatMessagesDiv.innerHTML = ""; // Clear previous content

    messages.forEach(message => {
        const messageDiv = document.createElement("div");
        messageDiv.classList.add("chat-message");

        messageDiv.innerHTML = `
            <div class="sender">${message.SenderName}</div>
            <div class="text">${message.Text}</div>
            <div class="sent-date">${message.SentDate}</div>
        `;

        chatMessagesDiv.appendChild(messageDiv);
    });
}
