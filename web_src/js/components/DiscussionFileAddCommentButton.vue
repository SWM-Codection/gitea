<script>

import {SvgIcon} from '../svg.js';

export default {
  

    components : {SvgIcon},

    data : () => ({
        isHovered : false
    }),


    methods : {
        
        hovered() {

            this.isHovered = true
        },
        notHovered() {
            this.isHovered = false
        }
    },
    mounted() {

    }

}
</script>

<template>

        <button class="ui primary button add-code-comment add-code-comment-right" > 
            <SvgIcon name="octicon-plus"/>
        </button>

</template>

<style>


</style>