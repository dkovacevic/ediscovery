document.addEventListener('DOMContentLoaded', function () {
    fetchJidData();
});

function fetchJidData() {
    fetch('/api/users', {
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
            populateTable(data);
        })
}

function populateTable(data) {
    const tableBody = document.querySelector('#jidTable tbody');

    if (!tableBody) {
        console.error('Table body not found');
        return;
    }

    data.forEach(item => {
        const row = document.createElement('tr');

        // Create a table cell for JID and make it a hyperlink
        const jidCell = document.createElement('td');
        const jidLink = document.createElement('a');
        jidLink.href = `/chats.html?lhid=${encodeURIComponent(item.jid)}&name=${encodeURIComponent(item.name)}`;
        jidLink.textContent = item.jid;
        jidCell.appendChild(jidLink);  // Append the link to the cell
        row.appendChild(jidCell);

        // Create a table cell for Name
        const nameCell = document.createElement('td');
        nameCell.textContent = item.name;
        row.appendChild(nameCell);

        // Create a table cell for Device
        const deviceCell = document.createElement('td');
        deviceCell.textContent = item.device;
        row.appendChild(deviceCell);

        // Append the row to the table body
        tableBody.appendChild(row);
    });
}
