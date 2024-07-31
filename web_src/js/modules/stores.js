import {reactive} from 'vue';

let diffTreeStoreReactive;
export function diffTreeStore() {
  if (!diffTreeStoreReactive) {
    diffTreeStoreReactive = reactive(window.config.pageData.diffFileInfo);
    window.config.pageData.diffFileInfo = diffTreeStoreReactive;
  }
  return diffTreeStoreReactive;
}

let repositoryFileStoreReactive; 
export function repositoryFileStore() {
  if (!repositoryFileStoreReactive) {
    repositoryFileStoreReactive = reactive(window.config.pageData.repositoryFileInfo);
    window.config.pageData.repositoryFileInfo = repositoryFileStoreReactive;
  }
  return repositoryFileStoreReactive; 
}