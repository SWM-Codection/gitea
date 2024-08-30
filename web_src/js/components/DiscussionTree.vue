<script>
import DiscussionTreeItem from './DiscussionTreeItem.vue';
import {toggleElem} from '../utils/dom.js';
import {discussionTreeStore} from '../modules/stores.js';
import {setFileFolding} from '../features/file-fold.js';

export default {
  components: {DiscussionTreeItem},
  data: () => {
    return {store: discussionTreeStore()};
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
      return result;
    },
  },
  mounted() {
    // Default to true if unset
    console.log('discussion tree mounted!')
    this.store.fileTreeIsVisible = true; 

    this.hashChangeListener = () => {
      this.store.selectedItem = window.location.hash;
      this.expandSelectedFile();
    };
    this.hashChangeListener();
    window.addEventListener('hashchange', this.hashChangeListener);
  },
  unmounted() {
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
  },
};
</script>
<template>

<div class="discussion-file-tree-items" >
  <!-- only render the tree if we're visible. in many cases this is something that doesn't change very often -->
  <DiscussionTreeItem v-for="item in fileTree" :key="item.name" :item="item"/>
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
