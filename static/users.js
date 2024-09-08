document.addEventListener('DOMContentLoaded', function() {
    fetchJidData();
});

function fetchJidData() {
    fetch('/api/users')
        .then(response => response.json())
        .then(data => {
            populateTable(data);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
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
        jidLink.href = `/all_chats.html?lhid=${encodeURIComponent(item.jid)}&name=${encodeURIComponent(item.name)}`;
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
