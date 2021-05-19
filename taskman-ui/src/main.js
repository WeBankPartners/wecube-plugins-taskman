import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ViewUI from 'view-design';
import VueI18n from "vue-i18n";
import 'view-design/dist/styles/iview.css';
import PluginSelect from "./components/select.vue";
import locale from "view-design/dist/locale/en-US";
import { quillEditor } from 'vue-quill-editor'
import TaskManQuillEditor from "./components/quillEditor"
import TaskFormItem from "./components/formItem"
import CmdbData from "./components/CmdbData"
import "./i18n"; 
Vue.use(ViewUI,{
  transfer: true,
  size: "default",
  VueI18n,
  locale
});


Vue.config.productionTip = false
Vue.component("PluginSelect", PluginSelect);
Vue.component("quillEditor", quillEditor);
Vue.component("TaskManQuillEditor", TaskManQuillEditor);
Vue.component("TaskFormItem", TaskFormItem);
Vue.component("CmdbData", CmdbData);

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
