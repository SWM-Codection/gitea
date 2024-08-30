import {reactive} from 'vue';

let diffTreeStoreReactive;
export function diffTreeStore() {
  if (!diffTreeStoreReactive) {
    diffTreeStoreReactive = reactive(window.config.pageData.diffFileInfo);
    window.config.pageData.diffFileInfo = diffTreeStoreReactive;
  }
  return diffTreeStoreReactive;
}

let discussionTreeStoreReactive;
export function discussionTreeStore() {
  if (!discussionTreeStoreReactive) {
    discussionTreeStoreReactive = reactive({
      repoLink: window.config.pageData.repoLink,
      files: [], 
      selectedItem: null, 
      contents: [], 
      checkedItems: [], 
    });
    window.config.pageData.discussionTreeInfo = discussionTreeStoreReactive;
  }
  return discussionTreeStoreReactive;
}

let discussionDetailStoreReactive; 
export function discussionDetailStore() {
  if (!discussionDetailStoreReactive) {
    discussionDetailStoreReactive = reactive({
      discussion: window.config.pageData.discussion, 
      discussionContent: window.config.pageData.discussionContent, 
    })
  }
  return discussionDetailStoreReactive;
}