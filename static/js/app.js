// Handle closing modal on HTMX trigger
document.addEventListener('DOMContentLoaded', function() {
    // Listen for HTMX events to close modal
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.xhr.getResponseHeader('HX-Trigger') === 'closeModal') {
            const modal = document.getElementById('modal');
            if (modal) {
                modal.classList.add('hidden');
            }
        }
    });
});