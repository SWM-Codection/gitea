<script>
import DiscussionFileTreeItem from './DiscussionFileTreeItem.vue';
import {loadMoreFiles} from '../features/repo-diff.js';
import {toggleElem} from '../utils/dom.js';
import {setFileFolding} from '../features/file-fold.js';
import {repositoryFileStore} from '../modules/stores.js';
const LOCAL_STORAGE_KEY = 'discussion_file_tree_visible';

export default {
  components: {DiscussionFileTreeItem},
  data: () => {
    return { store: repositoryFileStore() };
  },
  computed: {
    fileTree() {
      const result = [];
      for (const file of this.store.files) {
        // Split file into directories
        const splits = file.Name.split('/');
        let index = 0;
        let parent = null;
        let isFile = false;
        for (const split of splits) {
          index += 1;
          // reached the end
          if (index === splits.length) {
            isFile = true;
          }
          let newParent = {
            name: split,
            children: [],
            isFile,
          };

          if (isFile === true) {
            newParent.file = file;
          }

          if (parent) {
            // check if the folder already exists
            const existingFolder = parent.children.find(
              (x) => x.name === split,
            );
            if (existingFolder) {
              newParent = existingFolder;
            } else {
              parent.children.push(newParent);
            }
          } else {
            const existingFolder = result.find((x) => x.name === split);
            if (existingFolder) {
              newParent = existingFolder;
            } else {
              result.push(newParent);
            }
          }
          parent = newParent;
        }
      }
      const mergeChildIfOnlyOneDir = (entries) => {
        for (const entry of entries) {
          if (entry.children) {
            mergeChildIfOnlyOneDir(entry.children);
          }
          if (entry.children.length === 1 && entry.children[0].isFile === false) {
            // Merge it to the parent
            entry.name = `${entry.name}/${entry.children[0].name}`;
            entry.children = entry.children[0].children;
          }
        }
      };

      // Merge folders with just a folder as children in order to
      // reduce the depth of our tree.
      mergeChildIfOnlyOneDir(result);
      console.log('fileTree result is ', result);
      return result;
    },
  },
  mounted() {
    // document.querySelector('.discussion-toggle-file-tree-button').addEventListener('click', this.toggleVisibility);

    this.hashChangeListener = () => {
      this.store.selectedItem = window.location.hash;
      this.expandSelectedFile();
    };
    this.hashChangeListener();
    window.addEventListener('hashchange', this.hashChangeListener);
  },
  unmounted() {
    // document.querySelector('.diff-toggle-file-tree-button').removeEventListener('click', this.toggleVisibility);
    window.removeEventListener('hashchange', this.hashChangeListener);
  },
  methods: {
    expandSelectedFile() {
      // expand file if the selected file is folded
      if (this.store.selectedItem) {
        const box = document.querySelector(this.store.selectedItem);
        const folded = box?.getAttribute('data-folded') === 'true';
        if (folded) setFileFolding(box, box.querySelector('.fold-file'), false);
      }
    },
    updateState(visible) {
      // const btn = document.querySelector('.diff-toggle-file-tree-button');
      const [toShow, toHide] = btn.querySelectorAll('.icon');
      const tree = document.getElementById('discussion-file-tree');
      const newTooltip = btn.getAttribute(visible ? 'data-hide-text' : 'data-show-text');
      btn.setAttribute('data-tooltip-content', newTooltip);
      toggleElem(tree, visible);
      toggleElem(toShow, !visible);
      toggleElem(toHide, visible);
    },
    loadMoreData() {
      loadMoreFiles(this.store.linkLoadMore);
    },
  },
};
</script>
<template>
  <div class="discussion-file-tree-items">
    <!-- only render the tree if we're visible. in many cases this is something that doesn't change very often -->
    <DiscussionFileTreeItem v-for="item in fileTree" :key="item.name" :item="item"/>
    <div v-if="store.isIncomplete" class="tw-pt-1">
      <a :class="['ui', 'basic', 'tiny', 'button', store.isLoadingNewData ? 'disabled' : '']" @click.stop="loadMoreData">{{ store.showMoreMessage }}</a>
    </div>
  </div>
</template>
<style scoped>
.discussion-file-tree-items {
  display: flex;
  flex-direction: column;
  gap: 1px;
  margin-right: .5rem;
}
</style>
