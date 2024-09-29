import {POST} from '../modules/fetch.js';

async function fetchAiSampleCodes(data, aiCodeContainers) {
  try {
    // 로딩 애니메이션 표시
    const loadingImage = document.querySelector('.loading-overlay');
    if (loadingImage) {
      loadingImage.classList.remove('tw-hidden');
    }

    const response = await POST('/ai/discussion/samples', {data});

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

async function saveAiSampleCode(data, aiCodeModal) {
  try {
    const response = await POST('/ai/discussion/sample', {data});

    if (!response.ok) {
      throw new Error('Failed to save AI sample codes');
    }

    // 데이터 저장 성공 후 모달 닫기 (선택 사항)
    const isHidden = aiCodeModal.classList.contains('tw-hidden');
    if (!isHidden) aiCodeModal.classList.add('tw-hidden');
  } catch (error) {
    console.error('Error saving AI sample codes:', error);
    alert(`Error saving AI sample codes: ${error.message}`);
  }
}

export function initAiSampleCodeModal() {
  const modalShowBtns = document.querySelectorAll('.show-ai-code-modal');
  const aiCodeModal = document.querySelector('.ai-code-modal');
  const aiCodeModalClose = document.querySelector('.ai-code-modal-close');
  const aiCodeModalInsert = document.querySelector('.ai-code-modal-insert');
  const aiCodeContainers = document.querySelectorAll('.ai-code-area');
  let selectedCodeContainer = null;
  let commentId = null;

  if (!modalShowBtns.length) return;
  if (!aiCodeModal) return;
  if (!aiCodeModalClose) return;
  if (!aiCodeContainers.length) return;

  for (const modalShowBtn of modalShowBtns) {
    modalShowBtn.addEventListener('click', async ({target}) => {
      const tag = target.getAttribute('data-comment-id');
      commentId = parseInt(tag.split('-')[1]);

      const data = {
        target_comment_id: commentId.toString(),
        type: 'pull',
      };

      await fetchAiSampleCodes(data, aiCodeContainers);
    });
  }

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

    const data = {
      target_comment_id: commentId.toString(),
      sample_code_content: originalMarkdown, // original_markdown 값을 sample_code_content로 전달
      type: 'pull',
    };
    await saveAiSampleCode(data, aiCodeModal);
  });
}
