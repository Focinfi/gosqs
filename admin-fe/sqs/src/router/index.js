import Vue from 'vue'
import Router from 'vue-router'
import Try from '@/components/Try'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Try',
      component: Try
    }
  ]
})
