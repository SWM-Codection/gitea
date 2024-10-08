import $ from 'jquery';
import {hideElem, showElem, toggleElem} from '../utils/dom.js';
import {PATCH} from '../modules/fetch.js';

async function updateDeadline(deadlineString) {
  hideElem('#deadline-err-invalid-date');
  document.getElementById('deadline-loader')?.classList.add('is-loading');

  let realDeadline = null;
  if (deadlineString !== '') {
    const newDate = Date.parse(deadlineString);

    if (Number.isNaN(newDate)) {
      document.getElementById('deadline-loader')?.classList.remove('is-loading');
      showElem('#deadline-err-invalid-date');
      return false;
    }
    realDeadline = new Date(newDate);
  }

  try {
    const response = await PATCH(document.getElementById('update-discussion-deadline-form').getAttribute('action'), {
      data: {due_date: realDeadline},
      headers: {'Content-Type': 'application/json'},
    });

    if (response.ok) {
      window.location.reload();
    } else {
      throw new Error('Invalid response');
    }
  } catch (err) {
    console.error(err);
    document.getElementById('deadline-loader').classList.remove('is-loading');
    showElem('#deadline-err-invalid-date');
  }
}

// eslint-disable-next-line i/no-unused-modules
export function initDiscussionDue() {
  $(document).on('click', '.discussion-due-edit', () => {
    toggleElem('#deadlineForm');
  });
  $(document).on('click', '.discussion-due-remove', () => {
    updateDeadline('');
  });
  $(document).on('submit', '.discussion-due-form', () => {
    updateDeadline($('#deadlineDate').val());
    return false;
  });
}
