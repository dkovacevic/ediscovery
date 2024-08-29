document.addEventListener('DOMContentLoaded', function() {
    fetchJidData();
});

function fetchJidData() {
    fetch('api/devices')
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

    data.forEach(item => {
        const row = document.createElement('tr');

        const jidCell = document.createElement('td');
        jidCell.textContent = item.jid;
        row.appendChild(jidCell);

        const nameCell = document.createElement('td');
        nameCell.textContent = item.name;
        row.appendChild(nameCell);

        const deviceCell = document.createElement('td');
        deviceCell.textContent = item.device;
        row.appendChild(deviceCell);

        tableBody.appendChild(row);
    });
}
