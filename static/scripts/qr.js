document.addEventListener('DOMContentLoaded', function () {
    // Fetch QR code from the backend
    fetchQRCode();
});

function fetchQRCode() {
    fetch('/api/code', {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => {
            // If user is not authenticated, redirect to login
            if (response.redirected) {
                const currentPage = window.location.href;
                window.location.href = `/login.html?redirect=${encodeURIComponent(currentPage)}`;
                return;
            }
            if (!response.ok) {
                throw new Error('Failed to fetch QR code');
            }
            return response.json();
        })
        .then(data => {
            if (data.qr_code) {
                document.getElementById('qrCode').textContent = data.qr_code;
            } else {
                document.getElementById('qrCode').textContent = 'QR code could not be loaded.';
            }
        })
        .catch(error => {
            console.error('Error fetching QR code:', error);
            document.getElementById('qrCode').textContent = 'Error loading QR code.';
        });
}
