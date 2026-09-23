import { marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.min.css'

// marked + highlight.js 全局只配置一次，供文章详情与编辑器预览共用
marked.use(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code: string, lang: string) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      return hljs.highlight(code, { language }).value
    },
  }),
)

marked.setOptions({
  breaks: true,
  gfm: true,
})

/** 将 Markdown 渲染为 HTML（已内置代码高亮） */
export async function renderMarkdown(content: string): Promise<string> {
  if (!content) return ''
  return marked.parse(content) as Promise<string>
}
