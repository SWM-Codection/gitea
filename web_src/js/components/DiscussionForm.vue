<script>
export default {
    data: () => {
        store: {}
    },
    mounted() {
        console.log('vue component has been mounted')
    }
}
</script>

<template>
<form method="post">
    <div class="form-content tw-flex tw-p-4">
        <input type="text" placeholder="Title Here">
    </div>
</form>
</template>