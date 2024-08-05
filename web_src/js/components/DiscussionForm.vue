<script>
import { discussionTreeStore } from '../modules/stores.js';
import { GET, POST } from '../modules/fetch.js'; 
import DiscussionTree from './DiscussionTree.vue'; 
import EasyMDE from 'easymde';
export default {
    components: {
        DiscussionTree,
    },
    data: () => ({
        store: discussionTreeStore(), 
        branches: [], 
        selectedBranch: null,
        activeMenu: 'file', 
        mde: null, 
        name: '', 
    }),
    async mounted() {
        console.log('vue component has been mounted'); 
        await this.fetchBranches();
        await this.fetchFiles();
        this.mde = new EasyMDE({
            element: this.$refs.mde,
            forceSync: true, 
            renderingConfig: {singleLineBreaks: false}, 
            indentWithTabs: false, 
            tabSize: 4, 
            spellChecker: false, 
            inputStyle: 'contenteditable',
            nativeSpellCheck: true, 
        });
    },
    methods: {
        async fetchBranches() {
            const resp = await GET(`${this.store.repoLink}/branches/list`);
            const { results } = await resp.json(); 
            this.branches = results; 
            this.selectedBranch = results.length > 0 ? results[0] : null;
        },
        async fetchFiles() {
            if (!this.selectedBranch) return; 
            const resp = await GET(`${this.store.repoLink}/all-tree-list/branch/${this.selectedBranch}`);
            const files = await resp.json(); 
            this.store.files = files; 
            this.store.selectedItem = this.store.files.length > 0 ? '#discussion-' + this.store.files[0].NameHash : null; 
        },
        async fetchContents() { 
            this.store.contents = [];
            if (!this.selectedBranch) return;
            if (!this.store.selectedItem) return;
            const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
            if (filtered.length <= 0) return;
            const resp = await GET(`${this.store.repoLink}/raw/branch/${this.selectedBranch}/${filtered[0].Name}`);
            const content = await resp.text();
            content.split('\n').forEach((line, idx) => {
                this.store.contents.push({
                    line: idx + 1, 
                    content: line
                });
            });
        },
        async handleSumbit() {
            console.log('handle submit has been called~')
            const name = this.name; 
            const content = this.mde.value(); 
            const branchName = this.selectedBranch;
            const codes = await Promise.all(this.store.checkedItems.map(async (filePath) => {
                const filtered = this.store.files.filter(file => file.Name === filePath)
                if (filtered.length <= 0) return; 
                const resp = await GET(`${this.store.repoLink}/raw/branch/${this.selectedBranch}/${filePath}`);
                const content = await resp.text();
                const lines = content.split('\n').length + 1;  
                return {
                    filePath,
                    startLine: 1, 
                    endLine: lines + 1, 
                };
            }));


            // const formData = new FormData();
            // formData.append('name', name);
            // formData.append('content', content);
            // formData.append('branchName', branchName);
            // formData.append('codes', JSON.stringify(codes));            

            // i'm working on here! 
            const resp = await POST(`${this.store.repoLink}/discussions/new`, {
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name, 
                    content, 
                    branchName, 
                    codes
                }),
            });
            const result = await resp.text(); 
            console.log(result);
        }
    },
    computed: {
        filename() {
            const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
            if (filtered.length <= 0) return '파일을 선택해주세요'; 
            return filtered[0].Name;
        },
        contents() {
            return [
                {
                    line: 1, 
                    content: 'hello'
                },
                {
                    line: 2, 
                    content: 'world'
                },
            ]
        },
    },
    watch: {
        'store.selectedItem': function() {
            this.fetchContents();
        },
    }
}
</script>

<template>
<div class="form-content">
    <div class="field">
        <input type="text" placeholder="타이틀" v-model="name">
    </div>
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
        <div id="file-tab-content" v-show="activeMenu === 'file'">
            <div style="display: flex;">
                <DiscussionTree id="discussion-tree" style="max-height: 600px; width: 200px; overflow: auto;"/>
                <div class="tw-w-full tw-px-1" style="max-height: 600px; overflow: auto;">
                    <div class="file-header ui top attached header tw-items-center tw-justify-between tw-flex-wrap" style="position: sticky; top: 0; z-index: 999;">
                        <div class="file-info tw-font-mono">
                            <div class="file-info-entry">{{ filename }}</div>
                        </div>
                    </div>
                    <div class="ui bottom attached table unstackable segment">
                        <div class="file-view code-view">
                            <table>
                                <tbody>
                                    <tr v-for="content in store.contents">
                                        <td class="lines-num">{{ content.line }}</td>
                                        <td class="lines-code chroma">{{ content.content }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div id="write-tab-content" v-show="activeMenu === 'write'">
            <textarea class="EasyMDEContainer" ref="mde"></textarea>
        </div>

        <div class="text right tw-px-1 tw-mt-4">
            <button class="ui primary button" @click.prevent="handleSumbit">
                repo.discussion.new
            </button>
        </div>

    </div>
</div>
</template>