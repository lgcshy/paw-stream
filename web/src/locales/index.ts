import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import en from './en'

const LOCALE_KEY = 'pawstream_locale'

function getDefaultLocale(): string {
  const saved = localStorage.getItem(LOCALE_KEY)
  if (saved) return saved

  const browserLang = navigator.language
  if (browserLang.startsWith('zh')) return 'zh-CN'
  return 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    en: en,
  },
})

export function setLocale(locale: string) {
  ;(i18n.global.locale as any).value = locale
  localStorage.setItem(LOCALE_KEY, locale)
  document.documentElement.setAttribute('lang', locale === 'zh-CN' ? 'zh' : 'en')
}

export function getLocale(): string {
  return (i18n.global.locale as any).value
}

export default i18n
