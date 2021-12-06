import {createApp} from 'vue';
import App from './App.vue';
import routes from "./routes.vue";

createApp(App).use(routes).mount('#app')
