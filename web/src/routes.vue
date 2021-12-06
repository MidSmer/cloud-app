<template>

</template>

<script>
    import * as VueRouter from 'vue-router'
    import Article from "./pages/Article.vue";
    import About from "./pages/About.vue";

    const routes = [
        {path: '/', name: 'Home', component: Article, props: {articleKey: "", suppressEdit: false}},
        {
            path: '/article/:key', name: 'Article', component: Article, props: route => ({
                articleKey: route.params.key,
                suppressEdit: (route.query && route.query.hasOwnProperty('mode')) ? !(String(route.query.mode).toLowerCase() === "edit") : true,
                suppressDataUpdate: route.params.hasOwnProperty('suppressDataUpdate') ? (String(route.params.suppressDataUpdate).toLowerCase() === "true") : false,
            })
        },
        {path: '/about', component: About},
    ]

    const router = VueRouter.createRouter({
        history: VueRouter.createWebHistory(),
        routes,
        scrollBehavior(to, from, savedPosition) {
            if (savedPosition) {
                return savedPosition
            } else {
                if (to.hash) {
                    return {
                        el: to.hash,
                        behavior: 'smooth',
                    }
                }
                return {top: 0}
            }
        },
    })

    export default router
</script>

<style scoped>

</style>