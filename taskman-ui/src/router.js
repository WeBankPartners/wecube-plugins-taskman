import Vue from 'vue'
import Router from 'vue-router'
 

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      redirect: '/template-mgmt',
      component: () => import("./index.vue"),
      children: [
        // {
        //   path: '/task-management',
        //   name: 'TaskManagement',
        //   component: () => import("./views/Task-management.vue")
        // },
        // {
        //   path: '/service-catalog',
        //   name: 'ServiceCatalog',
        //   component: () => import("./views/Service-catalog.vue")
        // },
        {
          path: '/template-group',
          name: 'templateGroup',
          component: () => import("./views/template/template-group-mgmt.vue")
        },
        {
          path: '/template-mgmt',
          name: 'templateMgmt',
          component: () => import("./views/template/template-mgmt.vue")
        },
        {
          path: '/request-mgmt',
          name: 'requestMgmt',
          component: () => import("./views/request/request-mgmt.vue")
        },
        {
          path: '/task-mgmt',
          name: 'tasktMgmt',
          component: () => import("./views/task/task-mgmt.vue")
        },
      ]
    }
  ]
})

