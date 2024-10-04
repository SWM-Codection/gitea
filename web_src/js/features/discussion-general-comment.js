import {DELETE, POST, PUT} from "../modules/fetch.js";
import {hideElem, showElem} from '../utils/dom.js';
import {getComboMarkdownEditor, initComboMarkdownEditor} from './comp/ComboMarkdownEditor.js';

export function initDiscussionCommentDelete() {
  const deleteButtons = document.querySelectorAll('.discussion-delete-comment');
  console.log('searching delete comments...')
  if (!deleteButtons) return;
  console.log('find delete buttons done!');
  deleteButtons.forEach((deleteButton) => {
    deleteButton.addEventListener('click', async (event) => {
      event.preventDefault();
  
      const confirmationMessage = deleteButton.getAttribute('data-locale');
      const deleteUrl = deleteButton.getAttribute('data-url');
      const commentId = deleteButton.getAttribute('data-comment-id');
  
      if (window.confirm(confirmationMessage)) {
        try {
          const response = await DELETE(deleteUrl);
          if (!response.ok) throw new Error('Failed to delete comment');
  
          const commentElement = document.querySelector(`#${commentId}`);
          commentElement?.parentElement.remove();
        } catch (error) {
          console.error(error);
        }
      }
    });
  })
}


async function onEditContent(event) {
  event.preventDefault();

  const segment = this.closest('.header').nextElementSibling;
  const editContentZone = segment.querySelector('.edit-content-zone');
  const renderContent = segment.querySelector('.render-content');
  const rawContent = segment.querySelector('.raw-content');

  let comboMarkdownEditor;

  /**
   * @param {HTMLElement} dropzone
   */
  const setupDropzone = async (dropzone) => {
    if (!dropzone) return null;

    let disableRemovedfileEvent = false; // when resetting the dropzone (removeAllFiles), disable the "removedfile" event
    let fileUuidDict = {}; // to record: if a comment has been saved, then the uploaded files won't be deleted from server when clicking the Remove in the dropzone
    const dz = await createDropzone(dropzone, {
      url: dropzone.getAttribute('data-upload-url'),
      headers: {'X-Csrf-Token': csrfToken},
      maxFiles: dropzone.getAttribute('data-max-file'),
      maxFilesize: dropzone.getAttribute('data-max-size'),
      acceptedFiles: ['*/*', ''].includes(dropzone.getAttribute('data-accepts')) ? null : dropzone.getAttribute('data-accepts'),
      addRemoveLinks: true,
      dictDefaultMessage: dropzone.getAttribute('data-default-message'),
      dictInvalidFileType: dropzone.getAttribute('data-invalid-input-type'),
      dictFileTooBig: dropzone.getAttribute('data-file-too-big'),
      dictRemoveFile: dropzone.getAttribute('data-remove-file'),
      timeout: 0,
      thumbnailMethod: 'contain',
      thumbnailWidth: 480,
      thumbnailHeight: 480,
      init() {
        this.on('success', (file, data) => {
          file.uuid = data.uuid;
          fileUuidDict[file.uuid] = {submitted: false};
          const input = document.createElement('input');
          input.id = data.uuid;
          input.name = 'files';
          input.type = 'hidden';
          input.value = data.uuid;
          dropzone.querySelector('.files').append(input);
        });
        this.on('removedfile', async (file) => {
          document.getElementById(file.uuid)?.remove();
          if (disableRemovedfileEvent) return;
          if (dropzone.getAttribute('data-remove-url') && !fileUuidDict[file.uuid].submitted) {
            try {
              await POST(dropzone.getAttribute('data-remove-url'), {data: new URLSearchParams({file: file.uuid})});
            } catch (error) {
              console.error(error);
            }
          }
        });
        this.on('submit', () => {
          for (const fileUuid of Object.keys(fileUuidDict)) {
            fileUuidDict[fileUuid].submitted = true;
          }
        });
        this.on('reload', async () => {
          try {
            const response = await GET(editContentZone.getAttribute('data-attachment-url'));
            const data = await response.json();
            // do not trigger the "removedfile" event, otherwise the attachments would be deleted from server
            disableRemovedfileEvent = true;
            dz.removeAllFiles(true);
            dropzone.querySelector('.files').innerHTML = '';
            for (const el of dropzone.querySelectorAll('.dz-preview')) el.remove();
            fileUuidDict = {};
            disableRemovedfileEvent = false;

            for (const attachment of data) {
              const imgSrc = `${dropzone.getAttribute('data-link-url')}/${attachment.uuid}`;
              dz.emit('addedfile', attachment);
              dz.emit('thumbnail', attachment, imgSrc);
              dz.emit('complete', attachment);
              fileUuidDict[attachment.uuid] = {submitted: true};
              dropzone.querySelector(`img[src='${imgSrc}']`).style.maxWidth = '100%';
              const input = document.createElement('input');
              input.id = attachment.uuid;
              input.name = 'files';
              input.type = 'hidden';
              input.value = attachment.uuid;
              dropzone.querySelector('.files').append(input);
            }
            if (!dropzone.querySelector('.dz-preview')) {
              dropzone.classList.remove('dz-started');
            }
          } catch (error) {
            console.error(error);
          }
        });
      },
    });
    dz.emit('reload');
    return dz;
  };

  const cancelAndReset = (e) => {
    e.preventDefault();
    showElem(renderContent);
    hideElem(editContentZone);
    comboMarkdownEditor.attachedDropzoneInst?.emit('reload');
  };

  const saveAndRefresh = async (e) => {
    e.preventDefault();
    showElem(renderContent);
    hideElem(editContentZone);
    const dropzoneInst = comboMarkdownEditor.attachedDropzoneInst;
    try {
      const discussionId = parseInt(editContentZone.getAttribute('data-discussion-id')); 
      const discussionCommentId = parseInt(editContentZone.getAttribute('data-discussion-comment-id'));
      const response = await PUT(editContentZone.getAttribute('data-update-url'), {data: {
        discussionId, discussionCommentId, content: comboMarkdownEditor.value(),
      }});

      const data = await response.json();
      if (!data.content) {
        renderContent.innerHTML = document.getElementById('no-content').innerHTML;
        rawContent.textContent = '';
      } else {
        renderContent.innerHTML = data.content;
        rawContent.textContent = comboMarkdownEditor.value();
      }
      dropzoneInst?.emit('submit');
      dropzoneInst?.emit('reload');
    } catch (error) {
      console.error(error);
    }
  };

  comboMarkdownEditor = getComboMarkdownEditor(editContentZone.querySelector('.combo-markdown-editor'));
  if (!comboMarkdownEditor) {
    editContentZone.innerHTML = document.getElementById('discussion-comment-editor-template').innerHTML;
    comboMarkdownEditor = await initComboMarkdownEditor(editContentZone.querySelector('.combo-markdown-editor'));
    comboMarkdownEditor.attachedDropzoneInst = await setupDropzone(editContentZone.querySelector('.dropzone'));
    editContentZone.querySelector('.ui.cancel.button').addEventListener('click', cancelAndReset);
    editContentZone.querySelector('.ui.primary.button').addEventListener('click', saveAndRefresh);
  }

  // Show write/preview tab and copy raw content as needed
  showElem(editContentZone);
  hideElem(renderContent);
  if (!comboMarkdownEditor.value()) {
    comboMarkdownEditor.value(rawContent.textContent);
  }
  comboMarkdownEditor.focus();
}

export function initDiscussionGeneralEditContent() {
  console.log('init discussion general edit content')
  const editContents = document.querySelectorAll('.discussion-general-edit-content');
  editContents.forEach(($el) => {
    $el.addEventListener('click', onEditContent);
  })
}

export function initDiscussionCommentReaction() { 
  const $reactionButtons = document.querySelectorAll('.discussion-comment-reaction-button');
  $reactionButtons.forEach(($el) => {
    $el.onclick = async (e) => {
      const reactionType = $el.getAttribute('data-tooltip-content');
      const dataUrl = $el.getAttribute('data-url');
      const reacted = $el.getAttribute('data-has-reacted');
      const reactUrl = `${dataUrl}/${reacted ? 'unreact' : 'react'}`;

      const resp = await POST(reactUrl, {data: {'Content': reactionType}}); 

      const data = await resp.json();
      const html = data.html; 
      if (!html) return; 

      // handle reaction add 
      const $commentContainer = $el.closest('.comment-container');
      const $bottomReactions = $commentContainer.querySelector('.bottom-reactions');
      $bottomReactions?.remove(); 
      $commentContainer.insertAdjacentHTML('beforeend', html);
      initDiscussionCommentReaction();
    };
  });
}
