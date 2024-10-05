import { createApp } from "vue";
import DiscussionFileDetail from "../components/DiscussionFileDetail.vue";

export function initDiscussionFileView() {
    const el = document.getElementById('discussion-file-view');
    if (!el) return;
    const discussionFormView = createApp(DiscussionFileDetail);
    discussionFormView.mount(el);
}