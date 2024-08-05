import { createApp } from "vue";
import DiscussionForm from "../components/DiscussionForm.vue";

export function initDiscussionForm() {
    const el = document.getElementById('discussion-form');
    if (!el) return;
    const discussionFormView = createApp(DiscussionForm);
    discussionFormView.mount(el);
}