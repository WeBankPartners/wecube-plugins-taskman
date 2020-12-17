import TaskManagement from "./views/Task-management.vue"
import ServiceCatalog from "./views/Service-catalog.vue"
import TemplateGroup from "./views/template/template-group-mgmt.vue"
import TemplateMgmt from "./views/template/template-mgmt.vue"
import RequestMgmt from "./views/request/request-mgmt.vue"

export default [
    {
      path: '/template-group',
      name: 'templateGroup',
      component: TemplateGroup
    },
    {
      path: '/template-mgmt',
      name: 'templateMgmt',
      component: TemplateMgmt
    },
    {
      path: '/request-mgmt',
      name: 'requestMgmt',
      component: RequestMgmt
    },
  ]