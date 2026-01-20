<template>
  <div id="app">
    <Login v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
    <div v-else class="main-app">
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
            >
            <button type="submit" :disabled="loading">➕ Thêm</button>
          </form>
        </div>
        
        <div class="admin-link">
          <a href="http://localhost:8081/admin" target="_blank" class="admin-btn">
            📊 Xem báo cáo Admin
          </a>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import Login from './components/Login.vue'

export default {
  name: 'App',
  components: {
    Login
  },
  data() {
    return {
      isLoggedIn: false,
      username: '',
      newExpense: '',
      loading: false
    }
  },
  methods: {
    handleLoginSuccess(username) {
      this.isLoggedIn = true;
      this.username = username;
    },
    
    async addExpense() {
      if (!this.newExpense.trim()) return;
      
      this.loading = true;
      try {
        const response = await fetch('http://localhost:8081/api/expense', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            message: this.newExpense
          })
        });
        
        const data = await response.json();
        
        if (data.success) {
          this.newExpense = '';
          alert('✅ Đã thêm chi phí thành công!');
        } else {
          alert('❌ Lỗi: ' + (data.error || 'Không thể thêm chi phí'));
        }
      } catch (error) {
        alert('❌ Lỗi kết nối: ' + error.message);
      } finally {
        this.loading = false;
      }
    },
    
    async logout() {
      try {
        await fetch('http://localhost:8081/auth/logout', {
          credentials: 'include'
        });
      } catch (error) {
        console.error('Logout error:', error);
      }
      
      this.isLoggedIn = false;
      this.username = '';
    }
  }
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
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
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-header h1 {
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.logout-btn {
  background: #f44336;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 5px;
  cursor: pointer;
}

.logout-btn:hover {
  background: #d32f2f;
}

.app-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}

.expense-form {
  background: white;
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 30px;
}

.expense-form h2 {
  margin-bottom: 20px;
  color: #333;
}

.expense-form form {
  display: flex;
  gap: 10px;
}

.expense-form input {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 16px;
}

.expense-form button {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
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
  padding: 15px 30px;
  text-decoration: none;
  border-radius: 5px;
  font-size: 16px;
  font-weight: bold;
}

.admin-btn:hover {
  background: #F57C00;
}
</style>