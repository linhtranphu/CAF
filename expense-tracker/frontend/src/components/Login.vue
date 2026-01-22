<template>
  <div class="login-container">
    <div class="login-card">
      <h1>💰 Expense Tracker</h1>
      
      <div class="auth-tabs">
        <button 
          @click="activeTab = 'login'" 
          :class="{ active: activeTab === 'login' }"
          class="tab-btn"
        >
          🔐 Đăng nhập
        </button>
        <button 
          @click="activeTab = 'register'" 
          :class="{ active: activeTab === 'register' }"
          class="tab-btn"
        >
          📝 Đăng ký
        </button>
      </div>
      
      <!-- Login Form -->
      <form v-if="activeTab === 'login'" @submit.prevent="login" class="auth-form">
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
        
        <button type="submit" class="auth-btn" :disabled="loading">
          {{ loading ? '⏳ Đang đăng nhập...' : '🔐 Đăng nhập' }}
        </button>
        
        <div class="register-link">
          <p>Chưa có tài khoản? <a href="#" @click="$emit('switch-to-register')">Đăng ký ngay</a></p>
        </div>
      </form>
      
      <!-- Register Form -->
      <form v-if="activeTab === 'register'" @submit.prevent="register" class="auth-form">
        <div class="form-group">
          <label for="reg-username">Tên đăng nhập:</label>
          <input 
            type="text" 
            id="reg-username" 
            v-model="regUsername" 
            required 
            :disabled="loading"
          >
        </div>
        
        <div class="form-group">
          <label for="reg-password">Mật khẩu:</label>
          <input 
            type="password" 
            id="reg-password" 
            v-model="regPassword" 
            required 
            :disabled="loading"
          >
        </div>
        
        <button type="submit" class="auth-btn" :disabled="loading">
          {{ loading ? '⏳ Đang đăng ký...' : '📝 Đăng ký' }}
        </button>
      </form>
        
      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="success" class="success">{{ success }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      activeTab: 'login',
      username: '',
      password: '',
      regUsername: '',
      regPassword: '',
      loading: false,
      error: '',
      success: '',
      backendUrl: ''
    }
  },
  async mounted() {
    // Simple fallback - use window.location.hostname
    this.backendUrl = `http://${window.location.hostname}:8081`;
    console.log('Backend URL:', this.backendUrl);
  },
  methods: {
    async login() {
      this.loading = true;
      this.error = '';
      this.success = '';
      
      try {
        if (!this.backendUrl) {
          this.backendUrl = `http://${window.location.hostname}:8081`;
        }
        
        const response = await fetch(`${this.backendUrl}/auth/login`, {
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
        
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`);
        }
        
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
    },
    
    async register() {
      this.loading = true;
      this.error = '';
      this.success = '';
      
      try {
        if (!this.backendUrl) {
          this.backendUrl = `http://${window.location.hostname}:8081`;
        }
        
        const response = await fetch(`${this.backendUrl}/auth/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            username: this.regUsername,
            password: this.regPassword
          })
        });
        
        const data = await response.json();
        
        if (response.ok && data.message) {
          this.success = 'Đăng ký thành công! Bạn có thể đăng nhập ngay.';
          this.regUsername = '';
          this.regPassword = '';
          setTimeout(() => {
            this.activeTab = 'login';
            this.success = '';
          }, 2000);
        } else {
          this.error = data.error || 'Đăng ký thất bại';
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
  padding: 30px;
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  max-width: 400px;
  width: 100%;
}

h1 {
  color: #333;
  margin-bottom: 25px;
  text-align: center;
  font-size: 1.8rem;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
  font-size: 14px;
}

input {
  width: 100%;
  padding: 15px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
  transition: border-color 0.3s;
  -webkit-appearance: none;
}

input:focus {
  outline: none;
  border-color: #4285f4;
}

.auth-tabs {
  display: flex;
  margin-bottom: 20px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid #e1e5e9;
}

.tab-btn {
  flex: 1;
  padding: 12px;
  border: none;
  background: #f8f9fa;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.3s;
}

.tab-btn.active {
  background: #4285f4;
  color: white;
}

.tab-btn:hover:not(.active) {
  background: #e9ecef;
}

.auth-form {
  margin-top: 20px;
}

.success {
  color: #28a745;
  margin-top: 15px;
  text-align: center;
  padding: 12px;
  background: #d4edda;
  border: 1px solid #c3e6cb;
  border-radius: 8px;
  font-size: 14px;
}

.auth-btn {
  background: #4285f4;
  color: white;
  padding: 15px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  width: 100%;
  margin-top: 10px;
  transition: background-color 0.3s;
}

.auth-btn:hover:not(:disabled) {
  background: #357ae8;
}

.auth-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 20px;
}

.register-link a {
  color: #4285f4;
  text-decoration: none;
  font-weight: 600;
}

.register-link a:hover {
  text-decoration: underline;
}

.error {
  color: #dc3545;
  margin-top: 15px;
  text-align: center;
  padding: 12px;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 8px;
  font-size: 14px;
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .login-container {
    padding: 15px;
  }
  
  .login-card {
    padding: 25px 20px;
  }
  
  h1 {
    font-size: 1.5rem;
    margin-bottom: 20px;
  }
  
  input, .auth-btn {
    padding: 12px;
    font-size: 16px;
  }
  
  label {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .login-card {
    padding: 20px 15px;
  }
  
  h1 {
    font-size: 1.3rem;
  }
}
</style>