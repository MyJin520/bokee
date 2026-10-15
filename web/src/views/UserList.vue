<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getUserList } from '@/api/user'
import { ApiError } from '@/api/request'
import { useUserStore } from '@/stores/user'
import type { UserInfo } from '@/types'
import { formatDate, formatNumber, initialOf } from '@/utils/format'

const userStore = useUserStore()

/** 状态筛选项：后端 status 为精确匹配，空串表示不筛选 */
const STATUS_OPTIONS = [
  { label: '全部状态', value: '' },
  { label: '正常', value: 'normal' },
  { label: '已禁用', value: 'disabled' },
]

const STATUS_LABELS: Record<string, string> = {
  normal: '正常',
  disabled: '已禁用',
}

const pageSize = 10

const users = ref<UserInfo[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const errorMessage = ref('')
const forbidden = ref(false)
const filters = ref({ name: '', phone: '', email: '', status: '' })

const pages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

function statusLabel(status: string) {
  return STATUS_LABELS[status] || status || '未知'
}

async function load() {
  if (!userStore.isAdmin) {
    forbidden.value = true
    return
  }
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await getUserList({
      name: filters.value.name.trim(),
      phone: filters.value.phone.trim(),
      email: filters.value.email.trim(),
      status: filters.value.status,
      page: page.value,
      pageSize,
    })
    users.value = res.list
    total.value = res.total
  } catch (error: unknown) {
    users.value = []
    total.value = 0
    // 该接口仅对拥有权限的角色开放，无权限时引导提示而非报错
    if (error instanceof ApiError && error.code === 403) {
      forbidden.value = true
    } else {
      errorMessage.value = error instanceof Error ? error.message : '用户加载失败，请稍后重试。'
    }
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

function resetFilters() {
  filters.value = { name: '', phone: '', email: '', status: '' }
  search()
}

function changePage(next: number) {
  if (next < 1 || next > pages.value || next === page.value) return
  page.value = next
  load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  // 页面直接刷新进入时，等待根组件恢复登录态后再判断权限
  if (!userStore.userInfo && userStore.isLoggedIn) await userStore.restoreSession()
  if (!userStore.isAdmin) {
    forbidden.value = true
    return
  }
  load()
})
</script>

<template>
  <div class="users-page">
    <header class="page-header">
      <span class="eyebrow">Admin · users</span>
      <h1>用户管理</h1>
      <p>共 {{ formatNumber(total) }} 位用户，支持按用户名、手机号、邮箱模糊检索。</p>
    </header>

    <div v-if="forbidden" class="state-block">
      <p class="state-title">没有访问权限</p>
      <p>用户列表接口仅对拥有该权限的角色开放，当前账号无法查看。</p>
      <router-link class="btn btn-coral" to="/" style="margin-top: 18px">返回首页 →</router-link>
    </div>

    <template v-else>
      <form class="filter-bar" @submit.prevent="search">
        <label class="filter-field">
          <span class="field-label">用户名</span>
          <input v-model="filters.name" class="field-input" type="text" placeholder="模糊搜索" />
        </label>
        <label class="filter-field">
          <span class="field-label">手机号</span>
          <input v-model="filters.phone" class="field-input" type="text" placeholder="模糊搜索" />
        </label>
        <label class="filter-field">
          <span class="field-label">邮箱</span>
          <input v-model="filters.email" class="field-input" type="text" placeholder="模糊搜索" />
        </label>
        <label class="filter-field">
          <span class="field-label">状态</span>
          <select v-model="filters.status" class="field-input">
            <option v-for="option in STATUS_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </label>
        <div class="filter-actions">
          <button class="btn btn-coral btn-sm" type="submit" :disabled="loading">搜索</button>
          <button class="btn btn-outline btn-sm" type="button" :disabled="loading" @click="resetFilters">
            重置
          </button>
        </div>
      </form>

      <div v-if="loading" class="state-block">
        <div class="spinner"></div>
        <p>正在读取用户…</p>
      </div>
      <div v-else-if="errorMessage" class="state-block">
        <p class="state-title">加载遇到问题</p>
        <p>{{ errorMessage }}</p>
        <button class="text-button" type="button" style="margin-top: 12px" @click="load">重新加载 →</button>
      </div>
      <template v-else-if="users.length">
        <div class="table-wrap">
          <table class="user-table">
            <thead>
              <tr>
                <th>用户</th>
                <th>角色</th>
                <th>状态</th>
                <th>手机号</th>
                <th>邮箱</th>
                <th>粉丝</th>
                <th>注册时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>
                  <div class="user-cell">
                    <span class="user-avatar">
                      <img v-if="user.avatar" :src="user.avatar" :alt="user.name" />
                      <template v-else>{{ initialOf(user.name) }}</template>
                    </span>
                    <span>
                      <span class="user-name">{{ user.name }}</span>
                      <span class="user-id mono">ID {{ user.id }}</span>
                    </span>
                  </div>
                </td>
                <td>
                  <span v-if="user.roles.length" class="role-cell">
                    <span v-for="role in user.roles" :key="role.id" class="tag tag-coral">{{ role.name }}</span>
                  </span>
                  <span v-else class="cell-muted">未分配</span>
                </td>
                <td>
                  <span :class="['status-tag', user.status === 'normal' ? 'ok' : 'off']">
                    {{ statusLabel(user.status) }}
                  </span>
                </td>
                <td>{{ user.phone || '—' }}</td>
                <td>{{ user.email || '—' }}</td>
                <td class="mono">{{ formatNumber(user.followersCount) }}</td>
                <td class="cell-muted mono">{{ formatDate(user.createdAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="pagination">
          <button class="page-btn" type="button" :disabled="page <= 1" @click="changePage(page - 1)">← 上一页</button>
          <span class="page-info">{{ page }} / {{ pages }}</span>
          <button class="page-btn" type="button" :disabled="page >= pages" @click="changePage(page + 1)">
            下一页 →
          </button>
        </div>
      </template>
      <div v-else class="state-block empty-state">
        <p class="state-title">没有找到用户</p>
        <p>换一个关键词，或清除筛选条件试试。</p>
        <button class="text-button" type="button" style="margin-top: 12px" @click="resetFilters">清除筛选 →</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.users-page {
  padding: 48px 0 80px;
}

.page-header {
  margin-bottom: 30px;
}

.page-header h1 {
  margin: 14px 0 8px;
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(34px, 4.6vw, 54px);
  font-weight: 500;
  letter-spacing: -0.02em;
}

.page-header p {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.filter-bar {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr)) auto;
  align-items: end;
  gap: 18px;
  padding: 22px 24px;
  border: 1px solid var(--line);
  background: var(--surface);
}

.filter-field {
  display: block;
}

.filter-actions {
  display: flex;
  gap: 8px;
}

.table-wrap {
  margin-top: 22px;
  overflow-x: auto;
  border: 1px solid var(--line);
  background: var(--surface);
}

.user-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
  font-size: 12px;
}

.user-table th {
  padding: 12px 16px;
  border-bottom: 1px solid var(--line);
  color: var(--muted);
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-align: left;
  text-transform: uppercase;
  white-space: nowrap;
}

.user-table td {
  padding: 14px 16px;
  border-bottom: 1px solid var(--line);
  vertical-align: middle;
}

.user-table tbody tr:last-child td {
  border-bottom: 0;
}

.user-table tbody tr:hover {
  background: var(--paper);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 11px;
}

.user-avatar {
  width: 34px;
  height: 34px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 50%;
  overflow: hidden;
  background: var(--ink);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  display: block;
  font-weight: 700;
}

.user-id {
  display: block;
  margin-top: 3px;
  color: var(--muted);
  font-size: 10px;
}

.role-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.status-tag {
  padding: 3px 8px;
  font-size: 10px;
  font-weight: 800;
}

.status-tag.ok {
  background: var(--lime);
  color: var(--lime-dark);
}

.status-tag.off {
  background: var(--coral-tint);
  color: var(--coral-dark);
}

.cell-muted {
  color: var(--muted);
}

.empty-state {
  border: 1px dashed var(--line);
  background: var(--surface);
  margin-top: 22px;
}

@media (max-width: 1024px) {
  .filter-bar {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .filter-bar {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
