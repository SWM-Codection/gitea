<script>
import {SvgIcon} from '../svg.js';
import {discussionTreeStore} from '../modules/stores.js';

export default {
  components: {SvgIcon},
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data: () => ({
    store: discussionTreeStore(),
    collapsed: false,
  })
};
</script>
<template>
  <!--title instead of tooltip above as the tooltip needs too much work with the current methods, i.e. not being loaded or staying open for "too long"-->
  <a
    v-if="item.isFile" class="item-file"
    :class="{'selected': store.selectedItem === '#discussion-' + item.file.NameHash}"
    :title="item.name" :href="'#discussion-' + item.file.NameHash"
  >
    <!-- file -->
    <SvgIcon name="octicon-file"/>
    <span class="gt-ellipsis tw-flex-1">{{ item.name }}</span>
    <input type="checkbox" :value="item.file.Name" v-model="store.checkedItems">
  </a>
  <div v-else class="item-directory" :title="item.name" @click.stop="collapsed = !collapsed">
    <!-- directory -->
    <SvgIcon :name="collapsed ? 'octicon-chevron-right' : 'octicon-chevron-down'"/>
    <SvgIcon class="text primary" name="octicon-file-directory-fill"/>
    <span class="gt-ellipsis">{{ item.name }}</span>
  </div>

  <div v-if="item.children?.length" v-show="!collapsed" class="sub-items">
    <DiscussionTreeItem v-for="childItem in item.children" :key="childItem.name" :item="childItem"/>
  </div>
</template>
<style scoped>
a, a:hover {
  text-decoration: none;
  color: var(--color-text);
}

.sub-items {
  display: flex;
  flex-direction: column;
  gap: 1px;
  margin-left: 13px;
  border-left: 1px solid var(--color-secondary);
}

.sub-items .item-file {
  padding-left: 18px;
}

.item-file.selected {
  color: var(--color-text);
  background: var(--color-active);
  border-radius: 4px;
}

.item-file.viewed {
  color: var(--color-text-light-3);
}

.item-file,
.item-directory {
  display: flex;
  align-items: center;
  gap: 0.25em;
  padding: 3px 6px;
}

.item-file:hover,
.item-directory:hover {
  color: var(--color-text);
  background: var(--color-hover);
  border-radius: 4px;
  cursor: pointer;
}
</style>
