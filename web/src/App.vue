<template>
    <div>
        <router-view></router-view>
    </div>
</template>

<script>
    import {defineComponent, ref} from 'vue';

    export default {
        name: 'app',
        data() {
            return {
                currentRoute: window.location.pathname,
            }
        },

        created() {
            window.addEventListener('popstate', () => {
                this.currentRoute = window.location.pathname
            })
        }
    }
</script>

<style scoped>

</style>

