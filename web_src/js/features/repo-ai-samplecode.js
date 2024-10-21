import {POST} from '../modules/fetch.js';
import $ from 'jquery';

export async function fetchAiSampleCodes(data, aiCodeContainers) {
  try {

    // 로딩 애니메이션 표시
    const loadingImage = document.querySelector('.loading-overlay');
    if (loadingImage) {
      loadingImage.classList.remove('tw-hidden');
    }

    const response = await POST('/ai/samples', {data});

    if (!response.ok) {
      const result = await response.json();
      if (result.message && result.message.includes('already Ai comment')) {
        loadingImage.classList.add('tw-hidden');
        alert('이미 샘플 코멘트가 존재합니다');
        return;
      }
      throw new Error('Failed to fetch AI sample codes');
    }

    const result = await response.json();

    // 응답이 3개의 텍스트를 가진 배열인지 확인
    if (Array.isArray(result) && result.length === 3) {
      for (const [index, sampleCodeObj] of result.entries()) {
        if (aiCodeContainers[index]) {
          // innerHTML을 사용하여 마크다운이 적용된 HTML을 표시하고, 원본 마크다운도 저장
          aiCodeContainers[index].innerHTML = sampleCodeObj.sample_code;
          aiCodeContainers[index].setAttribute('data-original-markdown', sampleCodeObj.original_markdown);
          aiCodeContainers[index].setAttribute('data-discussionId', data.discussionId)
          aiCodeContainers[index].setAttribute('data-codeId', data.codeId)
          aiCodeContainers[index].setAttribute('data-endLine', data.endLine)
          aiCodeContainers[index].setAttribute('data-startLine', data.startLine)
        }
      }
    } else {
      throw new Error('Unexpected response format');
    }

    // 응답이 성공적으로 왔을 때 모달을 표시
    const aiCodeModal = document.querySelector('.ai-code-modal');
    if (aiCodeModal.classList.contains('tw-hidden')) {
      aiCodeModal.classList.remove('tw-hidden');
    }
  } catch (error) {
    console.error('Error fetching AI sample codes:', error);
    for (const container of aiCodeContainers) {
      container.innerHTML = 'Failed to load AI code samples.';
    }
  } finally {
    // 로딩 애니메이션을 다시 숨김
    const loadingImage = document.querySelector('.loading-overlay');
    if (loadingImage) {
      loadingImage.classList.add('tw-hidden');
    }
  }
}

async function saveAiDiscussionSampleCode(data, aiCodeModal) {
  const response = await POST('/ai/discussion/sample', {data});

  try {
    if (!response.ok) {
      throw new Error('Failed to save AI discussion sample code');
    }

    const result = await response.text();

    const $newCommentHolder = $(result);

    // Replace the current comment with the new one
    const selector = `#line-${data.codeId}-${data.endLine}`;
    const $targetLine = $(selector);

    if ($targetLine.length) {
      const tr = document.createElement("tr")
      const td = document.createElement("td")
      td.setAttribute("colspan", "3")
      td.append($newCommentHolder.get(0))
      tr.append(td);
      const lineElement = $targetLine.get(0)
      lineElement.insertAdjacentElement("afterend", tr);
    } else {
      console.warn('Could not find the discussion comment holder with the given selector.');
    }

    // Activate dropdown functionality for the new comment
    $newCommentHolder.find('.dropdown').dropdown();

    if (!aiCodeModal.classList.contains('tw-hidden')) {
      aiCodeModal.classList.add('tw-hidden');
    }
  } catch (error) {
    console.error('Error saving AI discussion sample code:', error);
    alert(`Error saving AI discussion sample code: ${error.message}`);
  }
}


async function saveAiPullSampleCode(data, aiCodeModal) {
  try {
    const response = await POST('/ai/pull/sample', {data});

    if (!response.ok) {

      throw new Error('Failed to save AI sample codes');
    }

    const $newConversationHolder = $(await response.text());
    const {path, side, idx} = $newConversationHolder.data();

    // 현재 코멘트 위치를 새로운 코멘트로 교체
    const selector = `.conversation-holder[data-path="${path}"][data-side="${side}"][data-idx="${idx}"]`;
    const $currentCommentHolder = $(selector);

    if ($currentCommentHolder.length) {
      $currentCommentHolder.replaceWith($newConversationHolder);
    } else {
      console.warn('Could not find the comment holder with the given selector.');
    }

    // 새로 추가된 코멘트에 대한 드롭다운 기능 활성화
    $newConversationHolder.find('.dropdown').dropdown();

    if (!aiCodeModal.classList.contains('tw-hidden')) {
      aiCodeModal.classList.add('tw-hidden');
    }
  } catch (error) {
    console.error('Error saving AI sample codes:', error);
    alert(`Error saving AI sample codes: ${error.message}`);
  }
}

export function initAiSampleCodeModal() {
  const aiCodeModal = document.querySelector('.ai-code-modal');
  const aiCodeModalClose = document.querySelector('.ai-code-modal-close');
  const aiCodeModalInsert = document.querySelector('.ai-code-modal-insert');
  const aiCodeContainers = document.querySelectorAll('.ai-code-area');
  let selectedCodeContainer = null;

  if (!aiCodeModal) return;
  if (!aiCodeModalClose) return;
  if (!aiCodeContainers.length) return;


  aiCodeModalClose.addEventListener('click', () => {
    const isHidden = aiCodeModal.classList.contains('tw-hidden');
    if (!isHidden) aiCodeModal.classList.add('tw-hidden');
  });

  for (const container of aiCodeContainers) {
    container.addEventListener('click', () => {
      for (const c of aiCodeContainers) {
        c.classList.remove('tw-border-green-500', 'tw-border-4');
        c.style.borderColor = '';
      }

      container.classList.add('tw-border-green-500', 'tw-border-4');
      container.style.borderColor = '#22c55e';

      selectedCodeContainer = container;
    });
  }

  aiCodeModalInsert.addEventListener('click', async () => {
    if (!selectedCodeContainer) {
      alert('코드 영역을 선택하세요.');
      return;
    }

    const originalMarkdown = selectedCodeContainer.getAttribute('data-original-markdown');
    const startLine = selectedCodeContainer.getAttribute('data-startLine');
    const endLine = selectedCodeContainer.getAttribute('data-endLine');
    const discussionId = selectedCodeContainer.getAttribute(`data-discussionId`);
    const codeId = selectedCodeContainer.getAttribute('data-codeId');

    const data = {
      origin_data: "",
      codeId: Number(codeId),
      discussionId: Number(discussionId),
      startLine: Number(startLine),
      endLine: Number(endLine),
      sample_code_content: originalMarkdown,
      type: "discussion"
    };

    await saveAiDiscussionSampleCode(data, aiCodeModal);

  });
}
