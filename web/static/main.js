// Vote functionality
document.addEventListener('DOMContentLoaded', function() {
    const voteBtn = document.getElementById('voteBtn');
    const voteCountEl = document.getElementById('voteCount');
    const thankYouModal = document.getElementById('thankYouModal');
    const alreadyVotedModal = document.getElementById('alreadyVotedModal');
    const errorModal = document.getElementById('errorModal');
    const errorMessage = document.getElementById('errorMessage');

    if (!voteBtn) return;

    // Disable button if already voted (persistent across page refreshes)
    if (typeof hasVoted !== 'undefined' && hasVoted) {
        voteBtn.disabled = true;
        voteBtn.textContent = 'Sudah Vote';
        voteBtn.classList.add('disabled');
    }

    // Function to add vote notice info
    function addVoteNotice() {
        const voteSection = voteBtn.closest('.vote-section');
        if (voteSection && !voteSection.querySelector('.vote-notice')) {
            const notice = document.createElement('p');
            notice.className = 'vote-notice';
            notice.textContent = '✓ Anda sudah memberikan vote untuk 1 karya inovasi';
            notice.style.cssText = 'color: #10b981; font-size: 0.875rem; margin-top: 0.75rem; text-align: center; font-weight: 500;';
            
            const subNotice = document.createElement('p');
            subNotice.className = 'vote-notice-sub';
            subNotice.textContent = 'Sistem ini hanya mengizinkan 1 vote per pengguna';
            subNotice.style.cssText = 'color: #6b7280; font-size: 0.75rem; margin-top: 0.25rem; text-align: center;';
            
            voteSection.appendChild(notice);
            voteSection.appendChild(subNotice);
        }
    }

    // Add notice if already voted on page load
    if (typeof hasVoted !== 'undefined' && hasVoted) {
        addVoteNotice();
    }

    voteBtn.addEventListener('click', async function() {
        // Confirm vote with copywriting
        if (!confirm('Yakin ingin vote untuk inovasi ini?\n\n⚠️ PERHATIAN: Anda hanya bisa vote untuk 1 karya saja. Pastikan pilihan Anda sudah tepat!')) {
            return;
        }

        // Disable button and show loading state
        voteBtn.disabled = true;
        voteBtn.classList.add('loading');
        const originalText = voteBtn.innerHTML;
        voteBtn.textContent = 'Memproses...';

        try {
            const response = await fetch(`/api/vote/${groupSlug}/${slug}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': csrfToken
                },
                credentials: 'same-origin'
            });

            const data = await response.json();

            // Reset button state
            voteBtn.classList.remove('loading');

            if (response.ok && data.success) {
                // Update vote count
                if (data.vote_count !== undefined) {
                    voteCountEl.textContent = data.vote_count;
                }

                // Show thank you modal
                thankYouModal.showModal();

                // Keep button disabled permanently
                voteBtn.disabled = true;
                voteBtn.textContent = '✓ Sudah Vote';
                voteBtn.classList.add('disabled');
                
                // Add vote notice
                addVoteNotice();
            } else if (response.status === 409 || data.error === 'already_voted') {
                // Already voted
                if (data.vote_count !== undefined) {
                    voteCountEl.textContent = data.vote_count;
                }

                alreadyVotedModal.showModal();
                voteBtn.disabled = true;
                voteBtn.textContent = '✓ Sudah Vote';
                voteBtn.classList.add('disabled');
                
                // Add vote notice
                addVoteNotice();
            } else {
                // Other error
                throw new Error(data.message || data.error || 'Terjadi kesalahan');
            }
        } catch (error) {
            console.error('Vote error:', error);
            
            // Reset button
            voteBtn.disabled = false;
            voteBtn.innerHTML = originalText;
            voteBtn.classList.remove('loading');

            // Show error modal
            errorMessage.textContent = error.message || 'Terjadi kesalahan. Silakan coba lagi.';
            errorModal.showModal();
        }
    });
});

// Close modal function
function closeModal() {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        if (modal.open) {
            modal.close();
        }
    });
}

// Close modal on backdrop click
document.addEventListener('click', function(event) {
    if (event.target.tagName === 'DIALOG') {
        const rect = event.target.getBoundingClientRect();
        const isInDialog = (rect.top <= event.clientY && event.clientY <= rect.top + rect.height &&
            rect.left <= event.clientX && event.clientX <= rect.left + rect.width);
        
        if (!isInDialog) {
            event.target.close();
        }
    }
});


