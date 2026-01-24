<template>
  <div id="app">
    <!-- Toast Notification -->
    <div v-if="toast.show" :class="['toast', toast.type]">
      {{ toast.message }}
    </div>
    
    <Login v-if="!isLoggedIn && currentView === 'login'" 
           @login-success="handleLoginSuccess" 
           @switch-to-register="currentView = 'register'" />
    
    <Register v-if="!isLoggedIn && currentView === 'register'" 
              @switch-to-login="currentView = 'login'" />
    
    <div v-if="isLoggedIn" class="main-app">
      <header class="app-header">
        <h1>💰 Expense Tracker</h1>
        <div class="user-info">
          <span>Xin chào, {{ username }}!</span>
          <button @click="logout" class="logout-btn">🚪 Đăng xuất</button>
        </div>
      </header>
      
      <main class="app-content">
        <div class="expense-form">
          <h2>📝 Thêm chi phí mới</h2>
          <form @submit.prevent="addExpense">
            <input 
              v-model="newExpense" 
              type="text" 
              placeholder="VD: ăn trưa 50k, mua xăng 200 nghìn"
              required
              :disabled="loading"
            >
            <button type="submit" :disabled="loading || !newExpense.trim()">
              {{ loading ? '⏳ Đang thêm...' : '➕ Thêm' }}
            </button>
          </form>
        </div>
        
        <div class="admin-link">
          <a :href="`${backendUrl}/admin`" target="_blank" class="admin-btn">
            📊 Xem báo cáo Admin
          </a>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import Login from './components/Login.vue'
import Register from './components/Register.vue'

export default {
  name: 'App',
  components: {
    Login,
    Register
  },
  data() {
    return {
      isLoggedIn: false,
      username: '',
      newExpense: '',
      loading: false,
      backendUrl: '',
      currentView: 'login',
      toast: {
        show: false,
        message: '',
        type: 'success'
      }
    }
  },
  async mounted() {
    // Simple fallback - use window.location.hostname
    this.backendUrl = `http://${window.location.hostname}:8081`;
    console.log('Backend URL:', this.backendUrl);
    
    // Check if user is already logged in
    await this.checkLoginStatus();
  },
  methods: {
    handleLoginSuccess(username) {
      this.isLoggedIn = true;
      this.username = username;
      // Save login state
      localStorage.setItem('isLoggedIn', 'true');
      localStorage.setItem('username', username);
    },
    
    async checkLoginStatus() {
      // Check localStorage first
      const savedLogin = localStorage.getItem('isLoggedIn');
      const savedUsername = localStorage.getItem('username');
      
      if (savedLogin === 'true' && savedUsername) {
        // Verify with server
        try {
          const response = await fetch(`${this.backendUrl}/api/health`, {
            credentials: 'include'
          });
          if (response.ok) {
            this.isLoggedIn = true;
            this.username = savedUsername;
            return;
          }
        } catch (error) {
          console.log('Session check failed:', error);
        }
      }
      
      // Clear invalid session
      localStorage.removeItem('isLoggedIn');
      localStorage.removeItem('username');
    },
    
    showToast(message, type = 'success') {
      this.toast.message = message;
      this.toast.type = type;
      this.toast.show = true;
      setTimeout(() => {
        this.toast.show = false;
      }, 3000);
    },
    
    async addExpense() {
      if (!this.newExpense.trim()) return;
      
      this.loading = true;
      try {
        // Ensure backendUrl is set
        if (!this.backendUrl) {
          this.backendUrl = `http://${window.location.hostname}:8081`;
        }
        
        console.log('API URL:', `${this.backendUrl}/api/expense`);
        
        const response = await fetch(`${this.backendUrl}/api/expense`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            message: this.newExpense
          })
        });
        
        if (!response.ok) {
          if (response.status === 401 || response.status === 403) {
            this.showToast('❌ Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.', 'error');
            this.isLoggedIn = false;
            return;
          }
          throw new Error(`HTTP ${response.status}`);
        }
        
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
          throw new Error('Server trả về không phải JSON');
        }
        
        const data = await response.json();
        
        if (data.success) {
          this.newExpense = '';
          
          // Show parsed summary
          if (data.parsed) {
            let summary = `✅ Đã thêm: ${data.parsed.items}`;
            if (data.parsed.quantity && data.parsed.unit) {
              summary += ` (${data.parsed.quantity} ${data.parsed.unit})`;
            } else if (data.parsed.quantity) {
              summary += ` (${data.parsed.quantity})`;
            }
            summary += ` - ${new Intl.NumberFormat('vi-VN').format(data.parsed.amount)} VND`;
            this.showToast(summary);
          } else {
            this.showToast('✅ Đã thêm chi phí thành công!');
          }
        } else {
          this.showToast('❌ Lỗi: ' + (data.error || 'Không thể thêm chi phí'), 'error');
        }
        
      } catch (error) {
        this.showToast('❌ Lỗi: ' + error.message, 'error');
      } finally {
        this.loading = false;
      }
    },
    
    async logout() {
      try {
        if (!this.backendUrl) {
          this.backendUrl = `http://${window.location.hostname}:8081`;
        }
        
        await fetch(`${this.backendUrl}/auth/logout`, {
          credentials: 'include'
        });
      } catch (error) {
        console.error('Logout error:', error);
      }
      
      this.isLoggedIn = false;
      this.username = '';
      // Clear saved login state
      localStorage.removeItem('isLoggedIn');
      localStorage.removeItem('username');
    }
  }
}
</script>

<style>
/* Toast Notification */
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 15px 20px;
  border-radius: 8px;
  color: white;
  font-weight: 600;
  z-index: 1000;
  animation: slideIn 0.3s ease-out;
  max-width: 300px;
  word-wrap: break-word;
}

.toast.success {
  background: #4CAF50;
}

.toast.error {
  background: #f44336;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  margin: 0;
  padding: 0;
}

.main-app {
  min-height: 100vh;
  background: #f5f5f5;
}

.app-header {
  background: #2196F3;
  color: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.app-header h1 {
  margin: 0;
  font-size: 1.5rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.9rem;
}

.logout-btn {
  background: #f44336;
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 0.8rem;
  white-space: nowrap;
}

.logout-btn:hover {
  background: #d32f2f;
}

.app-content {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.expense-form {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.expense-form h2 {
  margin-bottom: 15px;
  color: #333;
  font-size: 1.3rem;
}

.expense-form form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.expense-form input {
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  -webkit-appearance: none;
}

.expense-form button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 15px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
}

.expense-form button:hover:not(:disabled) {
  background: #45a049;
}

.expense-form button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.admin-link {
  text-align: center;
}

.admin-btn {
  display: inline-block;
  background: #FF9800;
  color: white;
  padding: 15px 25px;
  text-decoration: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: bold;
}

.admin-btn:hover {
  background: #F57C00;
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .toast {
    top: 10px;
    right: 10px;
    left: 10px;
    max-width: none;
    text-align: center;
  }
  
  .app-header {
    padding: 10px 15px;
  }
  
  .app-header h1 {
    font-size: 1.2rem;
  }
  
  .user-info {
    font-size: 0.8rem;
    gap: 8px;
  }
  
  .logout-btn {
    padding: 6px 10px;
    font-size: 0.7rem;
  }
  
  .app-content {
    padding: 15px;
  }
  
  .expense-form {
    padding: 15px;
  }
  
  .expense-form h2 {
    font-size: 1.1rem;
  }
  
  .expense-form form {
    gap: 12px;
  }
  
  .expense-form input,
  .expense-form button {
    padding: 12px;
    font-size: 16px;
  }
  
  .admin-btn {
    padding: 12px 20px;
    font-size: 14px;
  }
}
</style>