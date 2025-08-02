document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('email-form');
    const emailInput = document.getElementById('email-input');
    const checkButton = document.getElementById('check-button');
    const loading = document.getElementById('loading');
    const result = document.getElementById('result');
    const resultContent = document.getElementById('result-content');
    const checkAnotherButton = document.getElementById('check-another');

    const API_BASE_URL = '/api';

    // Email validation regex
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    // Form submission handler
    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const email = emailInput.value.trim();
        
        // Client-side validation
        if (!email) {
            showError('Please enter an email address');
            return;
        }
        
        if (!validateEmail(email)) {
            showError('Please enter a valid email address');
            return;
        }

        await checkEmail(email);
    });

    // Check another email button
    checkAnotherButton.addEventListener('click', function() {
        resetForm();
    });

    // Real-time email validation
    emailInput.addEventListener('input', function() {
        const email = this.value.trim();
        if (email && !validateEmail(email)) {
            this.setCustomValidity('Please enter a valid email address');
        } else {
            this.setCustomValidity('');
        }
    });

    function validateEmail(email) {
        return email.length >= 3 && 
               email.length <= 254 && 
               emailRegex.test(email);
    }

    async function checkEmail(email) {
        showLoading();

        try {
            const response = await fetch(`${API_BASE_URL}/check-email`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email: email })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'An error occurred while checking the email');
            }

            showResult(data);

        } catch (error) {
            console.error('Error checking email:', error);
            showError(error.message || 'Failed to check email. Please try again.');
        }
    }

    function showLoading() {
        form.style.display = 'none';
        loading.classList.remove('hidden');
        result.classList.add('hidden');
    }

    function showResult(data) {
        loading.classList.add('hidden');
        result.classList.remove('hidden');
        
        // Remove previous result classes
        result.classList.remove('safe', 'compromised', 'error');
        
        if (data.compromised) {
            result.classList.add('compromised');
            resultContent.innerHTML = `
                <h3>⚠️ Email Found in Data Breach</h3>
                <p><strong>${data.email}</strong></p>
                <p>${data.message}</p>
                <div style="margin-top: 1rem; padding: 1rem; background: rgba(0,0,0,0.1); border-radius: 6px;">
                    <strong>What should you do?</strong>
                    <ul style="text-align: left; margin-top: 0.5rem;">
                        <li>Change your password immediately</li>
                        <li>Enable two-factor authentication</li>
                        <li>Check for suspicious account activity</li>
                        <li>Consider using a password manager</li>
                    </ul>
                </div>
            `;
        } else {
            result.classList.add('safe');
            resultContent.innerHTML = `
                <h3>✅ Email Not Found in Known Breaches</h3>
                <p><strong>${data.email}</strong></p>
                <p>${data.message}</p>
                <div style="margin-top: 1rem; padding: 1rem; background: rgba(0,0,0,0.1); border-radius: 6px;">
                    <strong>Stay safe:</strong>
                    <ul style="text-align: left; margin-top: 0.5rem;">
                        <li>Use strong, unique passwords</li>
                        <li>Enable two-factor authentication</li>
                        <li>Keep your software updated</li>
                        <li>Be cautious with suspicious emails</li>
                    </ul>
                </div>
            `;
        }
    }

    function showError(message) {
        loading.classList.add('hidden');
        result.classList.remove('hidden');
        
        // Remove previous result classes and add error class
        result.classList.remove('safe', 'compromised');
        result.classList.add('error');
        
        resultContent.innerHTML = `
            <h3>❌ Error</h3>
            <p>${message}</p>
        `;
    }

    function resetForm() {
        form.style.display = 'block';
        loading.classList.add('hidden');
        result.classList.add('hidden');
        emailInput.value = '';
        emailInput.focus();
    }

    // Focus on email input when page loads
    emailInput.focus();

    // Test API connection on page load
    fetch(`${API_BASE_URL}/health`)
        .then(response => response.json())
        .then(data => {
            console.log('API Health Check:', data);
        })
        .catch(error => {
            console.warn('API connection test failed:', error);
        });
});