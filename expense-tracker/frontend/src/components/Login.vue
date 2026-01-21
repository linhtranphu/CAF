<template>
  <div class="login-container">
    <div class="login-card">
      <h1>💰 Expense Tracker</h1>
      
      <div class="users-info">
        <h4>👥 Demo Users:</h4>
        <p><strong>admin</strong> / admin123 (Admin)</p>
        <p><strong>linh</strong> / linh123 (User)</p>
        <p><strong>toan</strong> / toan123 (User)</p>
      </div>
      
      <form @submit.prevent="login" class="login-form">
        <div class="form-group">
          <label for="username">Tên đăng nhập:</label>
          <input 
            type="text" 
            id="username" 
            v-model="username" 
            required 
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="password">Mật khẩu:</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            required 
            :disabled="loading"
          >
        </div>
        
        <button type="submit" class="login-btn" :disabled="loading">
          {{ loading ? '⏳ Đang đăng nhập...' : '🔐 Đăng nhập' }}
        </button>
        
        <div v-if="error" class="error">{{ error }}</div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      username: '',
      password: '',
      loading: false,
      error: ''
    }
  },
  methods: {
    async login() {
      this.loading = true;
      this.error = '';
      
      try {
        const response = await fetch(`http://${window.location.hostname}:8081/auth/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            username: this.username,
            password: this.password
          })
        });
        
        const data = await response.json();
        
        if (data.message) {
          this.$emit('login-success', this.username);
        } else {
          this.error = data.error || 'Đăng nhập thất bại';
        }
      } catch (error) {
        this.error = 'Lỗi kết nối: ' + error.message;
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-card {
  background: white;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  max-width: 400px;
  width: 100%;
}

h1 {
  color: #333;
  margin-bottom: 30px;
  text-align: center;
}

.users-info {
  background: #f5f5f5;
  padding: 15px;
  border-radius: 5px;
  margin-bottom: 20px;
  font-size: 14px;
}

.users-info h4 {
  margin-bottom: 10px;
  color: #333;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 16px;
  box-sizing: border-box;
}

.login-btn {
  background: #4285f4;
  color: white;
  padding: 12px 24px;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
  width: 100%;
  margin-top: 10px;
}

.login-btn:hover:not(:disabled) {
  background: #357ae8;
}

.login-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.error {
  color: #f44336;
  margin-top: 10px;
  text-align: center;
  padding: 10px;
  background: #ffebee;
  border-radius: 5px;
}
</style>