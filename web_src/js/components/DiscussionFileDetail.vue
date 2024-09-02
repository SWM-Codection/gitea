

<script>
import { discussionFileTreeStore, discussionResponseDummy } from '../modules/stores.js';
import { GET, POST } from '../modules/fetch.js'; 
import DiscussionFileDetailTreeView from './DiscussionFileDetailTreeView.vue';
export default {
    components: {
        DiscussionFileDetailTreeView,
    },
    data: () => ({
        store: discussionFileTreeStore(), 
        dummy :  discussionResponseDummy(),
        activeMenu: 'file', 
        name: '', 
    }),
    async mounted() {
        console.log('vue component has been mounted'); 
        await this.fetchDiscussion();
    },

    methods: {

        // TODO globalComments, globalReaction 처리 추가
        async fetchDiscussion() {
            const resp =  this.dummy
            // const result = resp.json()
            const result = resp
            let contents = []
            let files = []
            result["contents"].forEach((content, idx) => {
                const nameHash = crypto.randomUUID()
                contents.push({
                    Name : content["filePath"],
                    codeBlocks : content["codeBlocks"],
                    NameHash: nameHash
                })
                files.push({
                    Name : content["filePath"],
                    NameHash : nameHash
                })
            })
            

            this.store.contents = contents
            this.store.files = files
            this.store.selectedItem = this.store.files.length > 0 ? '#discussion-' + this.store.files[0].NameHash : null; 
        },

        // async fetchBranches() {
        //     const resp = await GET(`${this.store.repoLink}/branches/list`);
        //     const { results } = await resp.json(); 
        //     this.branches = results; 
        //     this.selectedBranch = results.length > 0 ? results[0] : null;
        // },

        // async fetchFiles() {
        //     if (!this.selectedBranch) return; 
        //     const resp = await 
        //     const files = await resp.json(); 
        //     this.store.files = files; 
        //     this.store.selectedItem = this.store.files.length > 0 ? '#discussion-' + this.store.files[0].Name : null; 
        // },


        // async fetchContents() { 

        //     if (!this.selectedBranch) return;
        //     if (!this.store.selectedItem) return;
        //     const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
        //     if (filtered.length <= 0) return;
        //     const resp = await GET(`${this.store.repoLink}/raw/branch/${this.selectedBranch}/${filtered[0].NameHash}`);
        //     const content = await resp.text();
        //     const tmp = [] 
        //     content.split('\n').forEach((lineContent, idx) => {
        //         tmp.push({
        //             line: idx + 1, 
        //             content: lineContent
        //         });
        //     });
        //     this.store.contents = tmp; 
        // },
    },

    computed: {
        filename() {
            const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
            if (filtered.length <= 0) return '파일을 선택해주세요'; 
            return filtered[0].NameHash;
        },
    },
}

</script>



<template>
    <div class="form-content">

        <div class="field">
            <div class="ui top tabular menu">
                <a class="item" :class="{active: activeMenu === 'file'}" @click="activeMenu = 'file'">
                    <span class="resize-for-semibold">파일 선택</span>
                </a>
                <a class="item" :class="{active: activeMenu === 'write'}" @click="activeMenu = 'write'">
                    <span class="resize-for-semibold">쓰기</span>
                </a>
            </div>
        </div>

        <div style="margin-top: 1rem;">
            <div id="file-tab-content" style="display:flex;" v-show="activeMenu === 'file'">
                <DiscussionFileDetailTreeView id="discussion-tree" style="max-height: 600px; width: 200px; overflow: auto;"/>
                <div style="display: flex; flex-wrap: wrap;">
                    <div v-for ="content in store.contents" class="tw-w-full tw-px-1" style="max-height: 600px; margin-bottom: 2rem; overflow: auto;">
                        <div class="file-header ui top attached header tw-items-center tw-justify-between tw-flex-wrap" style="position: sticky; top: 0; z-index: 999;">
                            <div class="file-info tw-font-mono">
                                <div :href="`#discussion-${content.NameHash}`" class="file-info-entry">{{ content.Name }}</div>
                            </div>
                        </div>
                        <div class="ui bottom attached table unstackable segment">
                            <div class="file-view code-view" style="display: flex;">
                                <table>
                                    <tbody v-for = "codeBlock in content.codeBlocks" >
                                        <input type="hidden" :value = "`${codeBlock.codeId}`">
                                        <tr v-for="line in codeBlock.lines" :key ="`${codeBlock.codeId}-${codeBlock.lineNumber}`">
                                            <td class="lines-num">{{ line.lineNumber }}</td>
                                            <td class="lines-code chroma">{{ line.content }}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>