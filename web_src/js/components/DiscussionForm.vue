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
        collapseMenu: false, 
        selectedPreview: undefined, 
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
        window.addEventListener('resize', this.handleResize); 
    },
    beforeDestroy() {
        window.removeEventListener('resize', this.handleResize);
    }, 
    methods: {
        handleResize() {
            if (window.innerWidth < 1060) this.collapseMenu = true; 
        }, 
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
            this.store.isBin = content.isBin; 


            if (!this.store.isBin) {
                content.html.forEach((line, idx) => {
                    this.store.contents.push({
                        line: idx + 1, 
                        content: line
                    });
                });
            }
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
                codes: Array.from({length: this.dragEnd - this.dragStart + 1}, (_, i) => 1 + i)
                            .map(i => this.$refs.codeTable.querySelector(`td[data-line-no="${i + this.dragStart - 1}"] ~ td.lines-code`).outerHTML)
            }

            // XXX idk why, but Array.prototype.includes does not working 
            const hasItem = this.store.checkedItems.filter(item => (
                item.tag === currentItem.tag && 
                item.file === currentItem.file && 
                item.start === currentItem.start && 
                item.end === currentItem.end 
            )).length > 0;

            if (!hasItem) {
                this.store.checkedItems = [...this.store.checkedItems, currentItem]; 
            }
            this.dragStart = null; 
            this.dragLast = null; 
            this.dragEnd = null;
        },
        handleGotoCheckedFileRange({target}) {
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
                            <div class="file-view code-view" style="white-space-collapse: preserve;" v-if="!store.isBin">
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
                                                    class="discussion-add-button ui primary button tw-absolute tw-right-0 tw-bottom-0 tw-text-xs tw-opacity-50 hover:tw-opacity-100 tw-p-2">
                                                <span>선택 영역 추가하기</span>
                                            </button>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                            <div class="file-view code-view tw-px-4 tw-py-2" style="color: grey; " v-else>
                                binary file has been selected 
                            </div>

                        </div>
                    </div>
                </div>
            </div>

            <div id="write-tab-content" v-show="activeMenu === 'write'">
                <div class="code-preview tw-flex tw-max-h-[250px] tw-w-full tw-my-3" >
                    <div class="tw-flex tw-flex-col tw-overflow-auto" >
                        <span class="tw-flex tw-align-center tw-p-1" v-for="item in store.checkedItems" @click="selectedPreview = item"
                            :title="`${item.file}:${item.start}-${item.end}`"
                            :style="[selectedPreview == item && 'background-color: #eff0f2;' ]">
                            <SvgIcon name="octicon-file"/>
                            <div>{{ item.file }}:{{ item.start }}-{{  item.end }}</div>
                        </span>
                    </div> 

                    <div class="code-segment chroma tw-overflow-auto tw-grow tw-shrink-0 tw-bg-[#fbfdff] tw-px-2">
                        <tr v-for="(td, i) in selectedPreview?.codes">
                            <td class="tw-px-2" style="color: grey;">{{ selectedPreview?.start + i }}</td>
                            <td class="tw-px-2 tw-whitespace-pre" v-html="td"></td>
                        </tr>
                    </div>
                </div>
                <textarea class="EasyMDEContainer" ref="mde"></textarea>
            </div>

            <div class="text right tw-px-1 tw-mt-4">
                <button class="ui primary button" :class="[store.checkedItems.length > 0 && name.length > 0 ? null :  'disabled']" @click.prevent="handleSubmit">
                    새로운 디스커션 생성
                </button>
            </div>

        </div>
    </div>
</div>

<div class="discussion-content-right" 
    style="overflow: auto; flex-shrink: 0; max-width: 320px;"
    :style="[
        collapseMenu 
            ? 'position: absolute; right: 0; height: 100px; '
            : 'padding: 16px; max-height: 24px; margin-left: 18px; max-height: 990px; border: 1px solid #d0d7de; border-radius: 4px; '
    ]">
    <div aria-label="collapsable-menu" class="tw-flex tw-align-center tw-justify-end">
        <button @click.prevent="collapseMenu=!collapseMenu" class="tw-p-2" :style="[
            collapseMenu
                ? 'background: #f1f3f5; border-radius: 12px 0 0 12px; height: 100px; width: 24px; border: 1px solid #d0d7de;'
                : 'background: transparent;'
        ]">
            <SvgIcon :name="collapseMenu ? 'octicon-chevron-left' : 'octicon-x'"></SvgIcon>
        </button>
    </div>


    <div v-show="!collapseMenu">
        <span class="text muted flex-text-block" style="margin-bottom: 12px;">
            <SvgIcon name="octicon-git-branch"></SvgIcon>
            <strong>선택된 브랜치</strong>
        </span>

        <select class="tw-w-full" style="background-color: #f8f9fb; padding: 8px 12px;  border-radius: 6px; border: 1px solid #dcdde1;" v-model="selectedBranch" >
            <option disabled value="null">브랜치를 선택해주세요</option>
            <option :value="b" v-for="b in branches">{{b}}</option>
        </select>
        <div class="divider"></div>

        <span class="text muted flex-text-block" style="margin-bottom: 12px;">
            <SvgIcon name="octicon-file"></SvgIcon>
            <strong>선택된 파일 목록</strong>
        </span>

        <span style="color: grey;" v-if="store.checkedItems.length === 0">
            선택된 항목이 존재하지 않습니다.
        </span>
        <div v-else>
            <div v-for="item in store.checkedItems" 
                class="checked-item tw-relative"
                :title="`${item.file}:${item.start}-${item.end}`"
                :data-checked-item-tag="item.tag"
                :data-checked-item-file="item.file"
                :data-checked-item-start="item.start"
                :data-checked-item-end="item.end"
                @click="handleGotoCheckedFileRange"
                style="border: 1px solid #dcdde1; padding: 6px; margin: 3px; border-radius: 5px; cursor: pointer;">
                <SvgIcon name="octicon-file" class="tw-absolute tw-top-[8px] tw-left-[2px]"/>
                <span class="tw-inline-block tw-overflow-hidden tw-whitespace-nowrap tw-text-ellipsis tw-shrink tw-px-[16px] tw-max-w-[220px]">
                    {{ item.file }}:{{ item.start }}-{{  item.end }}
                </span>
                <SvgIcon name="octicon-x" class="tw-absolute tw-top-[8px] tw-right-[2px]" @click.prevent.stop="handleRemoveCheckedItem"/>
            </div>
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