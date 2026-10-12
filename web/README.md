# BOKEE Web

BOKEE 个人博客的前端工程，提供首页、文章归档与详情、写作与编辑、收藏关注、登录注册和个人中心等页面。

## 技术栈

- 构建工具：Vite
- 框架：Vue 3（组合式 API + `<script setup>`）
- 语言：TypeScript
- 路由：Vue Router
- 状态管理：Pinia
- 样式：原生 CSS（设计令牌集中在 `src/assets/main.css`）
- 内容渲染：marked + highlight.js

## 目录结构

```
web/
├─ public/              静态资源（站点图标）
├─ src/
│  ├─ api/              后端接口封装
│  ├─ assets/           全局样式与设计令牌
│  ├─ components/       通用组件（页头、提示、确认框、编辑器等）
│  ├─ composables/      组合式函数
│  ├─ router/           路由表
│  ├─ stores/           Pinia 状态
│  ├─ utils/            格式化与 Markdown 渲染
│  └─ views/            页面级组件
├─ index.html           应用入口 HTML
└─ vite.config.ts       构建与本地代理配置
```

## 环境要求

Node.js `^22.18.0 || >=24.12.0`。

## 本地开发

```sh
npm install
npm run dev
```

开发服务器默认运行在 `http://localhost:5173`，并把 `/pub`、`/pri` 前缀的请求代理到后端。

后端地址通过环境变量覆盖，默认 `http://localhost:8888`：

```sh
# .env.local
VITE_API_TARGET=http://localhost:8888
```

## 构建与预览

```sh
npm run build     # 类型检查 + 生产构建
npm run preview   # 本地预览构建产物
```

构建产物输出到 `dist/`，可直接交给静态站点服务托管。

## 代码质量

```sh
npm run lint      # oxlint + eslint
npm run type-check
npm run format    # prettier
```

## 站点图标

`public/favicon.svg`、`public/favicon.ico`、`public/apple-touch-icon.png` 为 BOKEE 品牌图标，
与页头品牌标（`src/components/AppHeader.vue` 中的 `brand-mark`）保持同一配色和几何比例。
