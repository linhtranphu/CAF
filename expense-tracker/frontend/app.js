const API_BASE = 'http://localhost:8081/api';

class ExpenseTracker {
    constructor() {
        this.offlineQueue = JSON.parse(localStorage.getItem('offlineQueue') || '[]');
        this.init();
    }

    init() {
        this.bindEvents();
        this.checkOnlineStatus();
        this.loadUserData();
        this.processOfflineQueue();
    }

    bindEvents() {
        document.getElementById('expenseForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        window.addEventListener('online', () => {
            this.updateOfflineIndicator(false);
            this.processOfflineQueue();
        });

        window.addEventListener('offline', () => {
            this.updateOfflineIndicator(true);
        });
    }

    loadUserData() {
        const savedUserId = localStorage.getItem('userId');
        if (savedUserId) {
            document.getElementById('userId').value = savedUserId;
        }
    }

    checkOnlineStatus() {
        this.updateOfflineIndicator(!navigator.onLine);
    }

    updateOfflineIndicator(isOffline) {
        const indicator = document.getElementById('offlineIndicator');
        if (isOffline) {
            indicator.classList.add('show');
        } else {
            indicator.classList.remove('show');
        }
    }

    async handleSubmit() {
        const userId = document.getElementById('userId').value.trim();
        const message = document.getElementById('message').value.trim();

        if (!userId || !message) {
            this.showStatus('Vui lòng điền đầy đủ thông tin', 'error');
            return;
        }

        // Save user ID
        localStorage.setItem('userId', userId);

        const expenseData = {
            message,
            userId
        };

        this.setLoading(true);

        if (navigator.onLine) {
            try {
                await this.sendToServer(expenseData);
                this.showStatus('✅ Đã ghi nhận chi phí thành công!', 'success');
                this.clearForm();
            } catch (error) {
                console.error('Error:', error);
                this.saveToOfflineQueue(expenseData);
                this.showStatus('📱 Đã lưu offline, sẽ đồng bộ khi có mạng', 'success');
                this.clearForm();
            }
        } else {
            this.saveToOfflineQueue(expenseData);
            this.showStatus('📱 Đã lưu offline, sẽ đồng bộ khi có mạng', 'success');
            this.clearForm();
        }

        this.setLoading(false);
    }

    async sendToServer(data) {
        const response = await fetch(`${API_BASE}/expense`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return response.json();
    }

    saveToOfflineQueue(data) {
        data.timestamp = Date.now();
        this.offlineQueue.push(data);
        localStorage.setItem('offlineQueue', JSON.stringify(this.offlineQueue));
    }

    async processOfflineQueue() {
        if (!navigator.onLine || this.offlineQueue.length === 0) {
            return;
        }

        const queue = [...this.offlineQueue];
        this.offlineQueue = [];
        localStorage.setItem('offlineQueue', JSON.stringify(this.offlineQueue));

        for (const item of queue) {
            try {
                await this.sendToServer(item);
                console.log('Synced offline item:', item);
            } catch (error) {
                console.error('Failed to sync offline item:', error);
                this.offlineQueue.push(item);
            }
        }

        if (this.offlineQueue.length > 0) {
            localStorage.setItem('offlineQueue', JSON.stringify(this.offlineQueue));
        }

        if (queue.length > 0) {
            this.showStatus(`🔄 Đã đồng bộ ${queue.length - this.offlineQueue.length} mục từ offline`, 'success');
        }
    }

    showStatus(message, type) {
        const statusDiv = document.getElementById('status');
        statusDiv.innerHTML = `<div class="status ${type}">${message}</div>`;
        
        setTimeout(() => {
            statusDiv.innerHTML = '';
        }, 5000);
    }

    setLoading(loading) {
        const submitBtn = document.getElementById('submitBtn');
        if (loading) {
            submitBtn.disabled = true;
            submitBtn.textContent = '⏳ Đang xử lý...';
        } else {
            submitBtn.disabled = false;
            submitBtn.textContent = '📝 Ghi nhận chi phí';
        }
    }

    clearForm() {
        document.getElementById('message').value = '';
    }
}

// Initialize app
new ExpenseTracker();