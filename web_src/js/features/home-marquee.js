import { createApp } from 'vue';
import HomeMarquee from '../components/HomeMarquee.vue'; 

export function initHomeMarquee() {
    const el = document.getElementById('home-marquee'); 
    if (!el) return; 
    const homeMarqueView = createApp(HomeMarquee); 
    homeMarqueView.mount(el); 
}