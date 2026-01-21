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
      error: '',
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
      
      try {
        // Ensure backendUrl is set
        if (!this.backendUrl) {
          this.backendUrl = `http://${window.location.hostname}:8081`;
        }
        
        console.log('Login URL:', `${this.backendUrl}/auth/login`);
        
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
        
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
          throw new Error('Server trả về không phải JSON');
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

.users-info {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 13px;
  border-left: 4px solid #4285f4;
}

.users-info h4 {
  margin-bottom: 10px;
  color: #333;
  font-size: 14px;
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

.login-btn {
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

.login-btn:hover:not(:disabled) {
  background: #357ae8;
}

.login-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
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
  
  .users-info {
    padding: 12px;
    font-size: 12px;
  }
  
  .users-info h4 {
    font-size: 13px;
  }
  
  input, .login-btn {
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
  
  .users-info {
    padding: 10px;
    font-size: 11px;
  }
}
</style>