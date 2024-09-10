document.addEventListener('DOMContentLoaded', function () {
    const params = new URLSearchParams(window.location.search);
    const lhid = params.get('lhid');
    const name = params.get('name');
    document.getElementById('name').textContent = decodeURIComponent(name);
    fetchChatData(lhid);
});

function fetchChatData(lhid) {
    fetch(`/api/${lhid}/chats`)
        .then(response => response.json())
        .then(data => {
            populateChatsTable(data, lhid);
        })
        .catch(error => {
            console.error('Error fetching chats:', error);
        });
}

function populateChatsTable(data, lhid) {
    const tableBody = document.querySelector('#chatsTable tbody');
    const name = document.getElementById('name').textContent

    if (!tableBody) {
        console.error('Table body not found');
        return;
    }

    data.forEach(item => {
        const row = document.createElement('tr');

        // Create a hyperlink for chatId
        const chatIdCell = document.createElement('td');
        const chatIdLink = document.createElement('a');
        chatIdLink.href = `messages.html?chatid=${encodeURIComponent(item.chatId)}&lhid=${encodeURIComponent(lhid)}&name=${encodeURIComponent(name)}`;
        chatIdLink.textContent = item.chatId;
        chatIdCell.appendChild(chatIdLink);
        row.appendChild(chatIdCell);

        // Add group name cell
        const groupNameCell = document.createElement('td');
        groupNameCell.textContent = item.groupName;
        row.appendChild(groupNameCell);

        // Add participants cell
        const participantsCell = document.createElement('td');
        // Check if participants is an array
        if (Array.isArray(item.participants)) {
            participantsCell.textContent = item.participants.join(', ');
        } else {
            participantsCell.textContent = 'N/A'; // or handle the case appropriately
        }
        row.appendChild(participantsCell);

        // Append the row to the table body
        tableBody.appendChild(row);
    });
}
