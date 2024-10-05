import { createApp } from "vue";
import DiscussionDetail from "../components/DiscussionDetail.vue";

export function initDiscussionForm() {
    const el = document.getElementById('discussion-detail');
    if (!el) return;
    const discussionDetailView = createApp(DiscussionDetail);
    discussionDetailView.mount(el);
}