<script>
import {discussionDetailStore} from '../modules/stores.js'
export default {
    data: () => ({
        store: discussionDetailStore(), 
    }),
    async mounted() {
        console.log('discussion detail store has been mounted');
    },
    methods: { 

    }
}
</script>

<template>
<div>
    this is discussion detail fuck you guys
</div>
</template>