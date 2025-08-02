document.addEventListener('DOMContentLoaded', function() {
    const TARGET_DATE = '2025-08-02 23:59:59'; // To update
    
    // Elements
    const daysElement = document.getElementById('days');
    const hoursElement = document.getElementById('hours');
    const minutesElement = document.getElementById('minutes');
    const secondsElement = document.getElementById('seconds');
    const countdownMessage = document.getElementById('countdown-message');
    const targetDateElement = document.getElementById('target-date');
    const accessButton = document.getElementById('access-button');
    
    // Convert target date string to Date object
    const targetDate = new Date(TARGET_DATE);
    
    // Display the target date in a readable format
    const options = { 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        timeZoneName: 'short'
    };
    targetDateElement.textContent = targetDate.toLocaleDateString('en-US', options);
    
    // Countdown timer function
    function updateCountdown() {
        const now = new Date().getTime();
        const timeRemaining = targetDate.getTime() - now;
        
        if (timeRemaining <= 0) {
            // Countdown has expired
            handleExpiry();
            return;
        }
        
        // Calculate time units
        const days = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));
        const hours = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((timeRemaining % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((timeRemaining % (1000 * 60)) / 1000);
        
        // Update display with leading zeros
        daysElement.textContent = String(days).padStart(2, '0');
        hoursElement.textContent = String(hours).padStart(2, '0');
        minutesElement.textContent = String(minutes).padStart(2, '0');
        secondsElement.textContent = String(seconds).padStart(2, '0');
        
        // Update access button state based on time remaining
        updateAccessButton(timeRemaining);
    }
    
    function updateAccessButton(timeRemaining) {
        const hoursRemaining = timeRemaining / (1000 * 60 * 60);
        
        if (hoursRemaining < 0) {
            accessButton.classList.add('disabled');
            accessButton.textContent = 'Demo Expired';
        } else {
            accessButton.classList.remove('disabled');
            accessButton.textContent = 'Launch Email Checker';
        }
    }
    
    function handleExpiry() {
        // Clear the countdown display
        daysElement.textContent = '00';
        hoursElement.textContent = '00';
        minutesElement.textContent = '00';
        secondsElement.textContent = '00';
        
        // Update message
        countdownMessage.innerHTML = '<strong>Demo has expired</strong>';
        countdownMessage.classList.add('expired');
        
        // Disable access button
        accessButton.classList.add('disabled');
        accessButton.textContent = 'Demo Expired';
        
        // Optionally, show expired message
        showExpiredMessage();
        
        // Stop the timer
        clearInterval(countdownInterval);
    }
    
    function showExpiredMessage() {
        const accessSection = document.querySelector('.access-section');
        const expiredDiv = document.createElement('div');
        expiredDiv.className = 'expired-message';
        expiredDiv.innerHTML = `
            <h3>Demo Period Ended</h3>
            <p>The demonstration period for this Razer Assignment has concluded. 
            The email checker tool is no longer accessible through the tunnel.\n
            Please contact me if this needs to be extended.</p>
            <p>Thank you for reviewing this technical demonstration.</p>
        `;
        
        // Replace access section with expired message
        accessSection.innerHTML = '';
        accessSection.appendChild(expiredDiv);
    }
    
    // Handle access button clicks
    accessButton.addEventListener('click', function(e) {
        if (this.classList.contains('disabled')) {
            e.preventDefault();
            return false;
        }
        
        console.log('User accessed email checker tool');
        // Note: Link now goes to razerassignmentapp.yoongjiahui.com
    });
    
    // Initialize countdown
    updateCountdown();
    
    // Start the countdown timer (update every second)
    const countdownInterval = setInterval(updateCountdown, 1000);
    
    // Utility function to check if date is valid
    function isValidDate(date) {
        return date instanceof Date && !isNaN(date);
    }
    
    // Validate target date
    if (!isValidDate(targetDate)) {
        console.error('Invalid target date format. Please use YYYY-MM-DD HH:MM:SS format.');
        countdownMessage.innerHTML = '<strong>Configuration Error: Invalid target date</strong>';
        countdownMessage.classList.add('expired');
    }
     
    // Handle page visibility changes (pause/resume timer when tab not active)
    let isPageVisible = true;
    
    document.addEventListener('visibilitychange', function() {
        if (document.hidden) {
            isPageVisible = false;
        } else {
            isPageVisible = true;
            // Immediately update when page becomes visible again
            updateCountdown();
        }
    });
});