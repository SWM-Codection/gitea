import { createApp } from "vue";
import DiscussionFileTree from '../components/DiscussionFileTree.vue'
import DiscussionFileList from '../components/DiscussionFileList.vue'

export function initDiscussionFileTree() {
    const el = document.getElementById('discussion-file-tree');
    if (!el) return;
    console.log('creating discussion file tree')
    const fileTreeView = createApp(DiscussionFileTree);
    fileTreeView.mount(el);

    const fileListElement = document.getElementById('discussion-file-list'); 
    if (!fileListElement) return; 

    const fileListView = createApp(DiscussionFileList);
    fileListView.mount(fileListElement);
}