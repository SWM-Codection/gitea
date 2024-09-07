<script>
import { discussionTreeStore } from '../modules/stores.js';
import { GET, POST } from '../modules/fetch.js'; 
import DiscussionTree from './DiscussionTree.vue'; 
import EasyMDE from 'easymde';
import {SvgIcon} from '../svg.js';

export default {
    components: {
        DiscussionTree,
        SvgIcon
    },
    data: () => ({
        store: discussionTreeStore(), 
        branches: [], 
        selectedBranch: null,
        activeMenu: 'file', 
        mde: null, 
        name: '', 
        dragging: false, 
        dragStart: undefined, 
        dragLast: undefined, 
        dragEnd: undefined, 
    }),
    async mounted() {
        await this.fetchBranches();
        this.selectedBranch = this.branches.length > 0 ? this.branches[0] : null;

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
        },
        async fetchFiles() {
            if (!this.selectedBranch) return; 
            const resp = await GET(`${this.store.repoLink}/all-tree-list/branch/${encodeURI(this.selectedBranch)}`);
            const files = await resp.json(); 
            this.store.files = files; 
        },
        async fetchContents() { 
            this.store.contents = [];
            if (!this.selectedBranch) return;
            if (!this.store.selectedItem) return;
            const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
            if (filtered.length <= 0) return;
            const resp = await GET(`${this.store.repoLink}/highlight/branch/${encodeURI(this.selectedBranch)}/${encodeURI(filtered[0].Name)}`);
            const content = await resp.json();
            content.html.forEach((line, idx) => {
                this.store.contents.push({
                    line: idx + 1, 
                    content: line
                });
            });
        },
        async handleSubmit() {
            const name = this.name; 
            const content = this.mde.value(); 
            const branchName = this.selectedBranch;
            const codes = await Promise.all(this.store.checkedItems.map(async (item) => {
                return {
                    filePath: item.file,
                    startLine: item.start, 
                    endLine: item.end, 
                };
            }));
            
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
            const {discussionId} = await resp.json(); 
            window.location.href = `${this.store.repoLink}/discussions/${discussionId}`
        },
        handleLineSelectStart({target}) {
            this.dragging = true; 
                    
            const line = target.getAttribute('data-line-no');
            this.dragStart = parseInt(line); 
            this.dragLast = this.dragStart; 
            this.dragEnd = null; 
        },
        handleLineSelectEnd({target}) { 
            if (!this.dragging) return; 

            const line = target.getAttribute('data-line-no'); 
            this.dragEnd = parseInt(line);

            // invalid drag  
            if (this.dragStart > this.dragEnd) {
                this.dragStart = this.dragEnd = null; 
            }
            // clean up drag behaviour
            // FIXME this dragging logic is not stable
            // it should be handled by pointer events not by mouse events
            this.dragging = false; 
        }, 
        handleLineSelectMove({target}) {
            if (!this.dragging) return ; 
            const line = target.getAttribute('data-line-no'); 
            this.dragLast = parseInt(line); 
        },
        handleAddDiscussionCode() {
            const currentItem = {
                tag: this.store.selectedItem,
                file: this.filename, 
                start: this.dragStart, 
                end: this.dragEnd,
            }

            // XXX idk why, but Array.prototype.includes does not working 
            const hsaItem = this.store.checkedItems.filter(item => (
                item.tag === currentItem.tag && 
                item.file === currentItem.file && 
                item.start === currentItem.start && 
                item.end === currentItem.end 
            )).length > 0;

            if (!hsaItem) {
                this.store.checkedItems = [...this.store.checkedItems, currentItem]; 
            }
            this.dragStart = null; 
            this.dragLast = null; 
            this.dragEnd = null;
        },
        handleGotoCheckedFileRange({target}) {
            console.log('goto file range', target); 
            const checkedItem = target.closest('.checked-item')
            const tag = checkedItem.getAttribute('data-checked-item-tag'); 
            const startNumber = parseInt(checkedItem.getAttribute('data-checked-item-start')); 
            const endNumber = parseInt(checkedItem.getAttribute('data-checked-item-end')); 
            this.store.reset = false; 
            this.store.selectedItem = tag; 
            this.dragStart = startNumber; 
            this.dragLast = endNumber; 
            this.dragEnd = endNumber; 
        },
        handleRemoveCheckedItem({target}) { 
            const checkedItem = target.closest('.checked-item');
            const tag = checkedItem.getAttribute('data-checked-item-tag'); 
            const file = checkedItem.getAttribute('data-checked-item-file'); 
            const startNumber = parseInt(checkedItem.getAttribute('data-checked-item-start')); 
            const endNumber = parseInt(checkedItem.getAttribute('data-checked-item-end')); 

            const filtered = this.store.checkedItems.filter(item => (
                item.tag !== tag || 
                item.file !== file ||
                item.start !== startNumber ||
                item.end !== endNumber 
            ));
            this.store.checkedItems = filtered;
        }
    },
    computed: {
        filename() {
            const filtered = this.store.files.filter(file => '#discussion-' + file.NameHash === this.store.selectedItem)
            if (filtered.length <= 0) return '파일을 선택해주세요'; 
            return filtered[0].Name;
        },
    },
    watch: {
        'store.selectedItem': function() {
            if (this.store.reset) {
                this.dragging = false; 
                this.dragStart = null; 
                this.dragLast = null; 
                this.dragEnd = null; 
            }
            this.fetchContents();
        },
        'selectedBranch': function() {
            this.fetchFiles()
            .then(() => this.store.selectedItem = this.store.files.length > 0 ? '#discussion-' + this.store.files[0].NameHash : null);
            this.store.checkedItems = []; 
        },
    },
}
</script>

<template>
<div class="discussion-content-left" style="flex-grow: 1; width: calc(100% -316px); border: 1px solid #d0d7de; border-radius: 4px; padding: 16px;">
    <div class="form-content">
        <div class="field">
            <input type="text" placeholder="타이틀" style="border: 1px solid #d0d7de; border-radius: 5px; width: 100%; padding: 10px; margin-bottom: 12px;"  v-model="name">
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
                    <DiscussionTree id="discussion-tree" style="max-width: 220px; width: 220px; max-height: 800px; overflow-y: scroll;"/>
                    <div class="tw-w-full tw-px-1" style="max-height: 800px; overflow: auto;">
                        <div class="file-header ui top attached header tw-items-center tw-justify-between tw-flex-wrap" style="position: sticky; top: 0; z-index: 999;">
                            <div class="file-info tw-font-mono">
                                <div class="file-info-entry">{{ filename }}</div>
                            </div>
                        </div>
                        <div class="ui bottom attached table unstackable segment">
                            <div class="file-view code-view" style="white-space-collapse: preserve;">
                                <table ref="codeTable">
                                    <tbody>
                                        <tr v-for="content in store.contents" 
                                            style="position: relative;" 
                                            :class="[
                                                (content.line === dragEnd) ? 'discussion-line-selected-end': null,
                                                (dragLast && dragStart <= content.line && content.line <= dragLast) ? 'discussion-line-selected' : null,
                                            ]">
                                            <td class="lines-num" style="cursor: ns-resize;" @mousedown="handleLineSelectStart" @mouseup="handleLineSelectEnd" @mousemove="handleLineSelectMove" :data-line-no="content.line">
                                                {{ content.line }}
                                            </td>
                                            <td class="lines-code chroma" v-html="content.content"/>
                                            <button @click.stop.prevent="handleAddDiscussionCode"
                                                    class="discussion-add-button"
                                                    style="
                                                        position: absolute; 
                                                        cursor: pointer;
                                                        color: black; 
                                                        background-color: #f1f3f5;
                                                        top: 16px;
                                                        right: 4px;
                                                        padding: 8px;
                                                        z-index: 1;
                                                        border-radius: 6px;
                                                        border: 1px solid grey;
                                                        justify-content: center;
                                                        font-size: smaller;">
                                                <span>선택 영역 추가하기</span>
                                            </button>
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
                <button class="ui primary button" :class="[store.checkedItems.length > 0 ? null :  'disabled']" @click.prevent="handleSubmit">
                    새로운 디스커션 생성
                </button>
            </div>

        </div>
    </div>
</div>

<div class="discussion-content-right" style="flex-shrink: 0; width: 320px; margin-left: 18px; border: 1px solid #d0d7de; border-radius: 4px; padding: 16px; max-height: 990px; overflow: auto;">
    <span class="text muted flex-text-block" style="margin-bottom: 12px;">
        <strong>브랜치 선택</strong>
    </span>

    <select class="tw-w-full" style="background-color: #f8f9fb; padding: 12px; border-radius: 6px; border: 1px solid #dcdde1;" v-model="selectedBranch" >
        <option disabled value="null">브랜치를 선택해주세요</option>
        <option :value="b" v-for="b in branches">{{b}}</option>
    </select>


    <div class="divider"></div>

    <span class="text muted flex-text-block" style="margin-bottom: 12px;">
        <strong>선택된 파일 목록</strong>
    </span>
    <span style="color: grey;" v-if="store.checkedItems.length === 0">
        선택된 항목이 존재하지 않습니다.
    </span>
    <div v-else>
        <div v-for="item in store.checkedItems" 
            class="checked-item"
            :data-checked-item-tag="item.tag"
            :data-checked-item-file="item.file"
            :data-checked-item-start="item.start"
            :data-checked-item-end="item.end"
            @click="handleGotoCheckedFileRange"
            style="border: 1px solid #dcdde1; padding: 6px; margin: 3px; border-radius: 5px; cursor: pointer;">

            <span style="display: flex; justify-content: space-between;">
                <div>
                    <SvgIcon name="octicon-file"/>{{ item.file }}:{{ item.start }}-{{  item.end }}
                </div>
                <div>
                    <SvgIcon name="octicon-x" @click.prevent.stop="handleRemoveCheckedItem"/>
                </div>  
            </span>
        </div>
    </div>
</div>
</template>

<style>
.discussion-line-selected {
    background-color: rgba(255, 248, 196, 0.7);
}


.discussion-add-button{ 
    visibility: hidden;
}

.discussion-line-selected-end  .discussion-add-button { 
    visibility: visible;
}
</style>