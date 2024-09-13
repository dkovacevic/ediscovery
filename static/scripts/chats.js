document.addEventListener('DOMContentLoaded', function () {
    const params = new URLSearchParams(window.location.search);
    const lhid = params.get('lhid');
    const name = params.get('name');
    document.getElementById('name').textContent = decodeURIComponent(name);
    fetchChatData(lhid);
});

function fetchChatData(lhid) {
    fetch(`/api/${lhid}/chats`, {
        method: 'GET',
        credentials: 'include'  // Include cookies in the request (important for JWT in cookies)
    })
        .then(response => {
            // If user is not authenticated, redirect to login
            if (response.redirected) {
                window.location.href = response.url;
                return;
            }
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data) {
                populateChatsTable(data, lhid);
            }
        })
}

function populateChatsTable(data, lhid) {
    const tableBody = document.querySelector('#chatsTable tbody');
    const name = document.getElementById('name').textContent;

    if (!tableBody) {
        console.error('Table body not found');
        return;
    }

    data.forEach(item => {
        const row = document.createElement('tr');

        // Create the hyperlink for the entire row
        const rowLink = `messages.html?chatid=${encodeURIComponent(item.chatId)}&lhid=${encodeURIComponent(lhid)}&name=${encodeURIComponent(name)}`;

        // Add Name cell
        const nameCell = document.createElement('td');
        nameCell.textContent = item.name !== "" ? item.name : item.groupName;
        row.appendChild(nameCell);

        // Add PhoneNo cell
        const phoneCell = document.createElement('td');
        phoneCell.textContent = item.phoneNo;
        row.appendChild(phoneCell);

        // Add Chat ID cell
        const chatIdCell = document.createElement('td');
        chatIdCell.textContent = item.chatId;
        row.appendChild(chatIdCell);

        // Make the entire row clickable
        row.addEventListener('click', () => {
            window.location.href = rowLink;
        });

        // Set the cursor to pointer to indicate it's clickable
        row.style.cursor = 'pointer';

        // Append the row to the table body
        tableBody.appendChild(row);
    });
}

