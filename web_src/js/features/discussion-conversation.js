import {PATCH} from '../modules/fetch.js';

export function initDiscussionStatusButton() {
  const discussionStatusButton = document.querySelector('#discussion-status-button');

  if (!discussionStatusButton) return;
  const mutedLink = document.querySelector('a[class="muted"]').getAttribute('href');

  discussionStatusButton.addEventListener('click', async (e) => {
    e.preventDefault();

    const target = e.target;
    const discussionId = target.getAttribute('data-discussion-id');
    const isClosed = target.getAttribute('data-is-closed') === 'true';

    // 서버에 상태 변경 요청을 보냄
    const params = new URLSearchParams({
      isClosed: (!isClosed).toString(),
    });

    try {
      const response = await PATCH(`${mutedLink}/discussions/state/${discussionId}?${params.toString()}`);

      if (!response.ok) {
        throw new Error('Failed to Change Status');
      }

      window.location.reload();
    } catch (error) {
      console.error('Error:', error);
    }
  });
}


