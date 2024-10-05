<script>
import { Vue3Marquee } from 'vue3-marquee';
export default {
    components: {Vue3Marquee}
}
</script>

<template>
    <Vue3Marquee :clone="true" :duration="3" :direction="'reverse'">
        <img src="https://sponsors.vuejs.org/images/skweezer_net.avif" height="120">
    </Vue3Marquee>
</template>