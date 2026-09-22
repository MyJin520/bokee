<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { changePassword, getUserInfo, updateUserInfo } from '@/api/user'
import { uploadFiles } from '@/api/file'
import { useUserStore } from '@/stores/user'
import { useUiStore } from '@/stores/ui'
import { formatDate, initialOf } from '@/utils/format'

const router = useRouter()
const userStore = useUserStore()
const ui = useUiStore()

const profileForm = reactive({
  userName: '',
  email: '',
  phone: '',
})
const avatar = ref('')
const savingProfile = ref(false)
const uploadingAvatar = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})
const savingPassword = ref(false)
const passwordErrors = ref<Record<string, string>>({})

const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']

function syncForm() {
  const info = userStore.userInfo
  if (!info) return
  profileForm.userName = info.name
  profileForm.email = info.email
  profileForm.phone = info.phone
  avatar.value = info.avatar
}

async function refreshUserInfo() {
  userStore.userInfo = await getUserInfo()
  syncForm()
}

function pickAvatar() {
  fileInput.value?.click()
}

async function onAvatarChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (!ALLOWED_TYPES.includes(file.type)) {
    ui.toastError('仅支持 JPG / PNG / GIF / WebP 格式')
    return
  }
  if (file.size > 32 * 1024 * 1024) {
    ui.toastError('图片大小不能超过 32MB')
    return
  }
  uploadingAvatar.value = true
  try {
    const files = await uploadFiles([file])
    if (files[0]?.url) {
      await updateUserInfo({ avatar: files[0].url })
      await refreshUserInfo()
      ui.toast('头像已更新')
    }
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '头像上传失败')
  } finally {
    uploadingAvatar.value = false
  }
}

async function saveProfile() {
  if (!profileForm.userName.trim()) {
    ui.toastError('用户名不能为空')
    return
  }
  if (profileForm.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profileForm.email.trim())) {
    ui.toastError('邮箱格式不正确')
    return
  }
  if (profileForm.phone && !/^1\d{10}$/.test(profileForm.phone.trim())) {
    ui.toastError('请输入正确的 11 位手机号')
    return
  }
  savingProfile.value = true
  try {
    await updateUserInfo({
      name: profileForm.userName.trim(),
      email: profileForm.email.trim(),
      phone: profileForm.phone.trim(),
    })
    localStorage.setItem('userName', profileForm.userName.trim())
    await refreshUserInfo()
    ui.toast('资料已保存')
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '保存失败，请稍后重试')
  } finally {
    savingProfile.value = false
  }
}

async function savePassword() {
  passwordErrors.value = {}
  if (!passwordForm.oldPassword) passwordErrors.value.oldPassword = '请输入旧密码'
  if (!passwordForm.newPassword || passwordForm.newPassword.length < 6)
    passwordErrors.value.newPassword = '新密码至少 6 位'
  if (passwordForm.newPassword !== passwordForm.confirmPassword)
    passwordErrors.value.confirmPassword = '两次输入的新密码不一致'
  if (Object.keys(passwordErrors.value).length) return

  savingPassword.value = true
  try {
    await changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    ui.toast('密码修改成功，下次登录请使用新密码')
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error: unknown) {
    ui.toastError(error instanceof Error ? error.message : '密码修改失败')
  } finally {
    savingPassword.value = false
  }
}

async function logout() {
  const ok = await ui.confirm({
    title: '退出登录',
    message: '确定要退出当前账号吗？',
    confirmText: '退出',
    danger: true,
  })
  if (!ok) return
  await userStore.doLogout()
  ui.toast('已退出登录')
  router.push('/')
}

onMounted(async () => {
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  syncForm()
})
</script>

<template>
  <div class="profile-page">
    <header class="profile-header">
      <span class="eyebrow">My account</span>
      <h1>个人中心</h1>
      <p v-if="userStore.userInfo">
        注册于 {{ formatDate(userStore.userInfo.createdAt) }} ·
        角色：{{ userStore.userInfo.roles.map((r) => r.name).join('、') || '普通用户' }}
      </p>
    </header>

    <div class="profile-grid">
      <section class="surface-card profile-card-block">
        <p class="side-title">Profile · 个人资料</p>

        <div class="avatar-row">
          <span class="big-avatar">
            <img v-if="avatar" :src="avatar" alt="头像" />
            <template v-else>{{ initialOf(profileForm.userName || '我') }}</template>
          </span>
          <div>
            <button class="btn btn-outline btn-sm" type="button" :disabled="uploadingAvatar" @click="pickAvatar">
              {{ uploadingAvatar ? '上传中…' : '更换头像' }}
            </button>
            <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/gif,image/webp" hidden @change="onAvatarChange" />
            <p class="form-hint">JPG / PNG / GIF / WebP，≤32MB</p>
          </div>
        </div>

        <label class="field-label">用户名</label>
        <input v-model="profileForm.userName" class="field-input" type="text" maxlength="10" />

        <label class="field-label">邮箱</label>
        <input v-model="profileForm.email" class="field-input" type="email" placeholder="选填" />

        <label class="field-label">手机号</label>
        <input v-model="profileForm.phone" class="field-input" type="tel" maxlength="11" placeholder="选填" />

        <div style="margin-top: 22px">
          <button class="btn" type="button" :disabled="savingProfile" @click="saveProfile">
            {{ savingProfile ? '保存中…' : '保存资料' }}
          </button>
        </div>
      </section>

      <div class="profile-side">
        <section class="surface-card profile-card-block">
          <p class="side-title">Password · 修改密码</p>
          <label class="field-label">旧密码</label>
          <input v-model="passwordForm.oldPassword" class="field-input" type="password" autocomplete="current-password" />
          <p v-if="passwordErrors.oldPassword" class="form-error">{{ passwordErrors.oldPassword }}</p>

          <label class="field-label">新密码</label>
          <input v-model="passwordForm.newPassword" class="field-input" type="password" autocomplete="new-password" />
          <p v-if="passwordErrors.newPassword" class="form-error">{{ passwordErrors.newPassword }}</p>

          <label class="field-label">确认新密码</label>
          <input v-model="passwordForm.confirmPassword" class="field-input" type="password" autocomplete="new-password" />
          <p v-if="passwordErrors.confirmPassword" class="form-error">{{ passwordErrors.confirmPassword }}</p>

          <div style="margin-top: 20px">
            <button class="btn btn-coral" type="button" :disabled="savingPassword" @click="savePassword">
              {{ savingPassword ? '提交中…' : '修改密码' }}
            </button>
          </div>
        </section>

        <section class="surface-card profile-card-block quick-card">
          <p class="side-title">Quick links · 快捷入口</p>
          <button class="quick-row" type="button" @click="router.push('/articles/create')">＋ 写新文章</button>
          <button class="quick-row" type="button" @click="router.push('/my-articles')">📄 我的文章</button>
          <button class="quick-row danger" type="button" @click="logout">退出登录</button>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 48px 0 80px;
}

.profile-header {
  margin-bottom: 32px;
}

.profile-header h1 {
  margin: 14px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(34px, 4.6vw, 52px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.profile-header p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
}

.profile-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(280px, 0.9fr);
  gap: 26px;
  align-items: start;
}

.profile-side {
  display: flex;
  flex-direction: column;
  gap: 26px;
}

.profile-card-block {
  padding: 26px;
}

.avatar-row {
  display: flex;
  align-items: center;
  gap: 18px;
  margin: 14px 0 6px;
}

.big-avatar {
  width: 72px;
  height: 72px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--ink);
  color: #fff;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 30px;
}

.big-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.quick-card {
  display: flex;
  flex-direction: column;
}

.quick-row {
  padding: 13px 4px;
  text-align: left;
  background: transparent;
  border-bottom: 1px solid var(--line);
  color: var(--ink);
  font-size: 13px;
}

.quick-row:last-child {
  border-bottom: 0;
}

.quick-row:hover {
  color: var(--coral-dark);
}

.quick-row.danger {
  color: var(--color-danger);
}

@media (max-width: 820px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }
}
</style>
