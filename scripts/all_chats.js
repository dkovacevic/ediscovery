document.addEventListener("DOMContentLoaded", function() {
    // Add event listener to the form
    document.getElementById("chatsForm").addEventListener("submit", function(e) {
        e.preventDefault(); // Prevent default form submission behavior

        const lhid = document.getElementById("lhid").value;

        // Fetch all chats using the provided lhid
        fetchChats(lhid);
    });
});

// Function to fetch all chats for a given lhid from the REST API
function fetchChats(lhid) {
    fetch(`/api/all-chats?lhid=${lhid}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // Expecting a JSON response
        })
        .then(data => {
            renderChats(data);
        })
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
}

// Function to render chats into HTML
function renderChats(chats) {
    const chatsDiv = document.getElementById("chats");
    chatsDiv.innerHTML = ""; // Clear previous content

    chats.forEach(chat => {
        const chatDiv = document.createElement("div");
        chatDiv.classList.add("chat");

        chatDiv.innerHTML = `
            <div class="chat-id">Chat ID: ${chat.ChatID}</div>
            <div class="group-name">Group Name: ${chat.GroupName}</div>
            <div class="participants">Participants: ${chat.Participants}</div>
        `;

        chatsDiv.appendChild(chatDiv);
    });
}
