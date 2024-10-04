

<script>
import { discussionFileTreeStore, discussionResponseDummy } from '../modules/stores.js';
import { GET, POST } from '../modules/fetch.js'; 
import DiscussionFileDetailTreeView from './DiscussionFileDetailTreeView.vue';
import DiscussionCodeLineSelector from './DiscussionCodeLineSelector.vue';

const {pageData} = window.config

export default {
    

    components: {
        DiscussionFileDetailTreeView,
        DiscussionCodeLineSelector
    },
    data: () => ({
        store: discussionFileTreeStore(), 
        dummy :  discussionResponseDummy(),
        name: '', 
        
    }),
    async mounted() {
        console.log('vue component has been mounted'); 
        await this.fetchDiscussion();
    },

    methods: {

        // TODO globalComments, globalReaction 처리 추가
        async fetchDiscussion() {
            const resp =  await GET(`${pageData.RepoLink}/discussions/${pageData.DiscussionId}/contents`)
            // const result = resp.json()
            const result = await resp.json()
            let contents = []
            let files = []
            result["contents"].forEach((content) => {
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


        <div style="margin-top: 1rem;">
            <div id="file-tab-content" style="display:flex; width: 1500px">
                <DiscussionFileDetailTreeView id="discussion-tree" style="width: 200px;"/>
                <div style="display: flex; flex-wrap: wrap;">
                    <div v-for ="content in store.contents" class="tw-w-full tw-px-1" style="margin-bottom: 2rem;">
                        <DiscussionCodeLineSelector :content = "content" />

                    </div>
                </div>
            </div>
        </div>
    </div>
</template>