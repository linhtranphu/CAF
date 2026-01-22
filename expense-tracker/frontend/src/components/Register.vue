<template>
  <div class="register-container">
    <div class="register-card">
      <h1>📝 Đăng ký tài khoản</h1>
      
      <form @submit.prevent="register" class="register-form">
        <div class="form-group">
          <label for="username">Tên đăng nhập:</label>
          <input 
            type="text" 
            id="username" 
            v-model="username" 
            required 
            :disabled="loading"
            minlength="3"
            maxlength="20"
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
            minlength="6"
          >
        </div>
        
        <div class="form-group">
          <label for="confirmPassword">Xác nhận mật khẩu:</label>
          <input 
            type="password" 
            id="confirmPassword" 
            v-model="confirmPassword" 
            required 
            :disabled="loading"
          >
        </div>
        
        <button type="submit" class="register-btn" :disabled="loading || !isFormValid">
          {{ loading ? '⏳ Đang đăng ký...' : '📝 Đăng ký' }}
        </button>
        
        <div class="login-link">
          <p>Đã có tài khoản? <a href="#" @click="$emit('switch-to-login')">Đăng nhập</a></p>
        </div>
        
        <div v-if="error" class="error">{{ error }}</div>
        <div v-if="success" class="success">{{ success }}</div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    return {
      username: '',
      password: '',
      confirmPassword: '',
      loading: false,
      error: '',
      success: '',
      backendUrl: ''
    }
  },
  computed: {
    isFormValid() {
      return this.username.length >= 3 && 
             this.password.length >= 6 && 
             this.password === this.confirmPassword;
    }
  },
  async mounted() {
    this.backendUrl = `http://${window.location.hostname}:8081`;
  },
  methods: {
    async register() {
      this.loading = true;
      this.error = '';
      this.success = '';
      
      if (this.password !== this.confirmPassword) {
        this.error = 'Mật khẩu xác nhận không khớp';
        this.loading = false;
        return;
      }
      
      try {
        const response = await fetch(`${this.backendUrl}/auth/register`, {
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
        
        if (response.ok && data.message) {
          this.success = '✅ Đăng ký thành công! Chuyển về trang đăng nhập...';
          setTimeout(() => {
            this.$emit('switch-to-login');
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
.register-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.register-card {
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

input:invalid {
  border-color: #f44336;
}

.register-btn {
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

.register-btn:hover:not(:disabled) {
  background: #357ae8;
}

.register-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.login-link {
  text-align: center;
  margin-top: 20px;
}

.login-link a {
  color: #4285f4;
  text-decoration: none;
  font-weight: 600;
}

.login-link a:hover {
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

/* Mobile optimizations */
@media (max-width: 768px) {
  .register-container {
    padding: 15px;
  }
  
  .register-card {
    padding: 25px 20px;
  }
  
  h1 {
    font-size: 1.5rem;
    margin-bottom: 20px;
  }
  
  input, .register-btn {
    padding: 12px;
    font-size: 16px;
  }
  
  label {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .register-card {
    padding: 20px 15px;
  }
  
  h1 {
    font-size: 1.3rem;
  }
}
</style>