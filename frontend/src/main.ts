import { PiniaColada } from "@pinia/colada";
import { createPinia } from "pinia";
import { createApp } from "vue";
import App from "./app/App.vue";
import { router } from "./app/router";
import { configureApiClient } from "./shared/lib/api-client";
import "./app/styles.css";

configureApiClient();

createApp(App).use(createPinia()).use(PiniaColada).use(router).mount("#app");
