document.addEventListener('DOMContentLoaded', function() {
    const params = new URLSearchParams(window.location.search);
    const chatId = params.get('chatid');
    const lhid = params.get('lhid');

    if (chatId && lhid) {
        fetchChatMessages(chatId, lhid);
    } else {
        console.error('ChatID or LHID not provided');
    }
});

function fetchChatMessages(chatId, lhid) {
    fetch(`/api/${lhid}/chats/${encodeURIComponent(chatId)}/messages`)
        .then(response => response.json())
        .then(data => {
            displayChatMessages(data);
        })
        .catch(error => {
            console.error('Error fetching chat messages:', error);
        });
}

function displayChatMessages(messages) {
    const chatMessagesDiv = document.getElementById('chat-messages');

    if (!chatMessagesDiv) {
        console.error('Chat messages container not found');
        return;
    }

    messages.forEach(message => {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'chat-message';

        const senderDiv = document.createElement('div');
        senderDiv.className = 'sender';
        senderDiv.textContent = message.sender;

        const textDiv = document.createElement('div');
        textDiv.className = 'text';
        textDiv.textContent = message.text;

        const sentDateDiv = document.createElement('div');
        sentDateDiv.className = 'date';
        sentDateDiv.textContent = message.date;

        messageDiv.appendChild(senderDiv);
        messageDiv.appendChild(textDiv);
        messageDiv.appendChild(sentDateDiv);

        chatMessagesDiv.appendChild(messageDiv);
    });
}
