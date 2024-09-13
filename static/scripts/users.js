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

        // Add a click event to the entire row
        row.addEventListener('click', function() {
            window.location.href = `/chats.html?lhid=${encodeURIComponent(item.jid)}`;
        });

        // Ensure that the cursor indicates the row is clickable
        row.style.cursor = 'pointer';

        // Create table cells
        const nameCell = document.createElement('td');
        nameCell.textContent = item.name !== "" ? item.name : item.groupName;
        row.appendChild(nameCell);

        const phoneCell = document.createElement('td');
        phoneCell.textContent = "+" + item.user;
        row.appendChild(phoneCell);

        const deviceCell = document.createElement('td');
        deviceCell.textContent = item.device;
        row.appendChild(deviceCell);

        const jidCell = document.createElement('td');
        jidCell.textContent = item.jid;
        row.appendChild(jidCell);

        // Append the row to the table body
        tableBody.appendChild(row);
    });
}
