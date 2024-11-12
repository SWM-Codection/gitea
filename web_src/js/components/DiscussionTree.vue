<script>
import DiscussionTreeItem from './DiscussionTreeItem.vue';
import {discussionTreeStore} from '../modules/stores.js';
import {setFileFolding} from '../features/file-fold.js';
import {debounce} from 'lodash'; 

export default {
  components: {DiscussionTreeItem},
  data: () => {
    return {
      localSearchInput: '', 
      store: discussionTreeStore(), 
    };
  },
  computed: {
    fileTree() {
      const result = [];
      for (const file of this.store.files) {
        // Split file into directories
        const filterEnabled = this.store.searchInput.length !== 0; 
        
        if (filterEnabled && !file.Name.includes(this.store.searchInput)) continue;
        const names = file.Name.split('/');
        let index = 0;
        let parent = null;
        let isFile = false;
        for (const name of names) {
          index += 1;
          // reached the end
          if (index === names.length) {
            isFile = true;
          }
          let newParent = {
            name: name,
            children: [],
            isFile,
          };

          if (isFile === true) {
            newParent.file = file;
          }

          if (parent) {
            // check if the folder already exists
            const existingFolder = parent.children.find(
              (x) => x.name === name,
            );
            if (existingFolder) {
              newParent = existingFolder;
            } else {
              parent.children.push(newParent);
            }
          } else {
            const existingFolder = result.find((x) => x.name === name);
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
    this.store.fileTreeIsVisible = true; 

    this.hashChangeListener = () => {
      this.store.reset = true; 
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
  watch: {
    localSearchInput: debounce(function (value) {
      this.store.searchInput = value; 
    }, 33),
  }
};
</script>
<template>

<div class="discussion-file-tree-items" >
  <!-- only render the tree if we're visible. in many cases this is something that doesn't change very often -->
  <input type="text" v-model="localSearchInput" placeholder="원하는 경로를 검색해보세요"  class="search-input">
  <DiscussionTreeItem v-for="item in fileTree" :key="item.name" :item="item" />
</div>

</template>
<style scoped>
.discussion-file-tree-items {
  display: flex;
  flex-direction: column;
  gap: 1px;
  margin-right: .5rem;
}

.search-input {
  border: 1px solid #d0d7de; 
  border-radius: 5px; 
  width: 100%; 
  padding: 5px; 
  margin-bottom: 6px; 

  position: sticky;
  top: 0;
}

.search-input::placeholder {
  font-size: 12px;
}
</style>
