// Function to get URL query parameters
function getQueryParams() {
    const params = new URLSearchParams(window.location.search);
    return {
        lhid: params.get('lhid'),
        chatID: params.get('chatid'),
    };
}

// Extract `lhid` and `chatID` from URL
const {lhid, chatID} = getQueryParams();
let currentPage = 1;  // Start with page 1
const limit = 15;  // Number of messages per page

// Fetch initial messages on page load
fetchMessages(currentPage);

// Function to fetch chat messages
async function fetchMessages(page = 1) {
    try {
        const response = await fetch(`/api/${encodeURIComponent(lhid)}/chats/${encodeURIComponent(chatID)}/messages?page=${page}&limit=${limit}`, {
            method: 'GET',
            credentials: 'include'  // Include cookies in the request (important for JWT in cookies)
        });

        // If user is not authenticated, redirect to login
        if (response.redirected) {
            window.location.href = response.url;
            return;
        }

        const data = await response.json();

        if (data.messages && data.messages.length > 0) {
            renderMessages(data.messages);
            updatePagination(data.total_pages, page);
        } else {
            document.getElementById('chat-messages').innerHTML = '<p>No messages available</p>';
        }
    } catch (error) {
        console.error('Error fetching messages:', error);
    }
}

// Function to render messages
function renderMessages(messages) {
    const chatMessagesDiv = document.getElementById('chat-messages');
    chatMessagesDiv.innerHTML = '';  // Clear previous messages
    document.getElementById('back-to-chats').href = `chats.html?lhid=${lhid}`;

    messages.forEach(message => {
        const messageDiv = document.createElement('div');
        messageDiv.classList.add('chat-message');
        messageDiv.innerHTML = `
            <div class="sender">${message.sender}</div>
            <div class="text">${message.text}</div>
            <div class="date">${message.date}</div>
        `;
        chatMessagesDiv.appendChild(messageDiv);
    });
}

// Function to update pagination controls
function updatePagination(totalPages, page) {
    const prevButton = document.getElementById('prev-page');
    const nextButton = document.getElementById('next-page');
    const pageInfo = document.getElementById('page-info');

    // Update page info
    pageInfo.textContent = `Page ${page} of ${totalPages}`;

    // Enable/disable buttons based on the current page
    prevButton.disabled = page === 1;
    nextButton.disabled = page === totalPages;
}

// Event listeners for pagination buttons
document.getElementById('prev-page').addEventListener('click', () => {
    if (currentPage > 1) {
        currentPage--;
        fetchMessages(currentPage);
    }
});

document.getElementById('next-page').addEventListener('click', () => {
    currentPage++;
    fetchMessages(currentPage);
});


