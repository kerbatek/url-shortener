const form = document.getElementById('shorten-form');
const result = document.getElementById('result');

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const url = document.getElementById('url-input').value;

    try {
        const res = await fetch('/shorten', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url }),
        });
        const data = await res.json();

        if (!res.ok) {
            result.className = 'error';
            result.textContent = data.error;
            result.style.display = 'block';
            return;
        }

        const shortURL = `${location.origin}/${data.code}`;
        result.className = '';
        result.innerHTML = `
            Short URL: <a href="${shortURL}" target="_blank">${shortURL}</a>
            <br><button class="delete-btn" onclick="deleteURL(${data.id})">Delete</button>
        `;
        result.style.display = 'block';
    } catch (err) {
        result.className = 'error';
        result.textContent = 'Something went wrong';
        result.style.display = 'block';
    }
});

async function deleteURL(id) {
    const res = await fetch(`/url/${id}`, { method: 'DELETE' });
    if (res.ok) {
        result.className = '';
        result.textContent = 'Deleted successfully';
    } else {
        result.className = 'error';
        result.textContent = 'Failed to delete';
    }
}
