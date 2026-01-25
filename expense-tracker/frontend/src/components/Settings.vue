<template>
  <div class="settings-container">
    <div class="settings-card">
      <h2>⚙️ Cài đặt</h2>
      
      <div class="setting-section">
        <h3>🤖 Gemini API Configuration</h3>
        <p class="description">Cấu hình API key để sử dụng AI parsing cho chi phí</p>
        
        <form @submit.prevent="saveApiKey">
          <div class="form-group">
            <label for="apiKey">Gemini API Key:</label>
            <input 
              id="apiKey"
              v-model="apiKey" 
              type="password"
              placeholder="Nhập Gemini API key của bạn"
              :disabled="loading"
            >
            <small>Lấy API key tại: <a href="https://aistudio.google.com/apikey" target="_blank">Google AI Studio</a></small>
          </div>
          
          <div class="button-group">
            <button type="submit" :disabled="loading || !apiKey.trim()" class="save-btn">
              {{ loading ? '⏳ Đang lưu...' : '💾 Lưu' }}
            </button>
            <button type="button" @click="testApiKey" :disabled="loading || !apiKey.trim()" class="test-btn">
              🧪 Test API
            </button>
            <button type="button" @click="$emit('close')" class="cancel-btn">
              ❌ Đóng
            </button>
          </div>
        </form>
        
        <div v-if="message" :class="['message', messageType]">
          {{ message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Settings',
  props: {
    backendUrl: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      apiKey: '',
      loading: false,
      message: '',
      messageType: 'success'
    }
  },
  async mounted() {
    await this.loadCurrentApiKey();
  },
  methods: {
    async loadCurrentApiKey() {
      try {
        const response = await fetch(`${this.backendUrl}/api/settings/gemini`, {
          credentials: 'include'
        });
        
        if (response.ok) {
          const data = await response.json();
          if (data.hasKey) {
            this.apiKey = '••••••••••••••••';
          }
        }
      } catch (error) {
        console.error('Load API key error:', error);
      }
    },
    
    async saveApiKey() {
      this.loading = true;
      this.message = '';
      
      try {
        const response = await fetch(`${this.backendUrl}/api/settings/gemini`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            apiKey: this.apiKey
          })
        });
        
        const data = await response.json();
        
        if (response.ok && data.success) {
          this.message = '✅ Đã lưu API key thành công!';
          this.messageType = 'success';
          setTimeout(() => {
            this.message = '';
          }, 3000);
        } else {
          throw new Error(data.error || 'Không thể lưu API key');
        }
      } catch (error) {
        this.message = '❌ Lỗi: ' + error.message;
        this.messageType = 'error';
      } finally {
        this.loading = false;
      }
    },
    
    async testApiKey() {
      this.loading = true;
      this.message = '';
      
      try {
        const response = await fetch(`${this.backendUrl}/api/settings/gemini/test`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify({
            apiKey: this.apiKey
          })
        });
        
        const data = await response.json();
        
        if (response.ok && data.success) {
          this.message = '✅ API key hoạt động tốt!';
          this.messageType = 'success';
        } else {
          throw new Error(data.error || 'API key không hợp lệ');
        }
      } catch (error) {
        this.message = '❌ Lỗi: ' + error.message;
        this.messageType = 'error';
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<style scoped>
.settings-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.settings-card {
  background: white;
  border-radius: 10px;
  padding: 30px;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
}

h2 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.5rem;
}

.setting-section {
  margin-bottom: 20px;
}

h3 {
  color: #2196F3;
  margin-bottom: 10px;
  font-size: 1.2rem;
}

.description {
  color: #666;
  margin-bottom: 15px;
  font-size: 0.9rem;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
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

input:focus {
  outline: none;
  border-color: #2196F3;
}

small {
  display: block;
  margin-top: 5px;
  color: #666;
  font-size: 0.85rem;
}

small a {
  color: #2196F3;
  text-decoration: none;
}

small a:hover {
  text-decoration: underline;
}

.button-group {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

button {
  padding: 12px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  flex: 1;
  min-width: 120px;
}

.save-btn {
  background: #4CAF50;
  color: white;
}

.save-btn:hover:not(:disabled) {
  background: #45a049;
}

.test-btn {
  background: #2196F3;
  color: white;
}

.test-btn:hover:not(:disabled) {
  background: #0b7dda;
}

.cancel-btn {
  background: #f44336;
  color: white;
}

.cancel-btn:hover {
  background: #d32f2f;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.message {
  margin-top: 15px;
  padding: 12px;
  border-radius: 5px;
  font-weight: 600;
}

.message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

@media (max-width: 768px) {
  .settings-card {
    padding: 20px;
  }
  
  h2 {
    font-size: 1.3rem;
  }
  
  h3 {
    font-size: 1.1rem;
  }
  
  .button-group {
    flex-direction: column;
  }
  
  button {
    width: 100%;
  }
}
</style>
