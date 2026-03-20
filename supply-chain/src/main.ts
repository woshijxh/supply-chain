import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import router from './router'
import i18n from './i18n'
import App from './App.vue'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Toast from 'primevue/toast'
import ToastService from 'primevue/toastservice'
import ConfirmDialog from 'primevue/confirmdialog'
import ConfirmationService from 'primevue/confirmationservice'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import DatePicker from 'primevue/datepicker'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import TooltipDirective from 'primevue/tooltip'

import 'primeicons/primeicons.css'
import 'remixicon/fonts/remixicon.css'
import '@fontsource/sora/index.css'
import '@fontsource/noto-sans-sc/index.css'
import './styles/main.scss'

import supplyChainPreset from './styles/primevue-theme'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)
app.use(PrimeVue, {
  ripple: true,
  theme: {
    preset: supplyChainPreset,
    options: {
      darkModeSelector: '.dark-mode'
    }
  }
})
app.use(ToastService)
app.use(ConfirmationService)

app.component('Button', Button)
app.component('InputText', InputText)
app.component('InputNumber', InputNumber)
app.component('Dropdown', Dropdown)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('Dialog', Dialog)
app.component('Toast', Toast)
app.component('ConfirmDialog', ConfirmDialog)
app.component('TabView', TabView)
app.component('TabPanel', TabPanel)
app.component('Card', Card)
app.component('Tag', Tag)
app.component('DatePicker', DatePicker)
app.component('Textarea', Textarea)
app.component('Select', Select)
app.component('Checkbox', Checkbox)
app.directive('tooltip', TooltipDirective)

app.mount('#app')