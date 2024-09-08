
// @ts-check

import { on } from "htmx.org"

export class CodeLinePosition {

    /**
     * 
     * @param {string} lineNumber 
     * @param {string} codeId 
     */
    constructor(lineNumber, codeId) {
        
        this.lineNumber = lineNumber
        this.codeId = codeId
    }


}

export class CodeLineRange {

    /**
     * @param {string} codeId 
     * @param {CodeLinePosition} startPosition 
     * @param {CodeLinePosition} endPosition 
     * 
     */
    constructor(codeId, startPosition, endPosition) {
        this.startPosition = startPosition
        this.endPosition = endPosition
        this.codeId = codeId
        this.elements = new Set()

        const table = document.getElementById(codeId)
        const tableLineElements = table?.children

        if (!tableLineElements) {
            return 
        }

        for (const lineElement of tableLineElements) {
            const lineNumber = lineElement.id.split("-")[1]
            if (this.isLineNumberOutofBound(lineNumber)) {
                return
            }

            this.elements.add(lineElement)
        }
        
    }

    /**
     * 
     * @param {string} lineNumber 
     */
    isLineNumberOutofBound(lineNumber) {
        return parseInt(lineNumber) < parseInt(this.startPosition.lineNumber) || 
        parseInt(lineNumber) > parseInt(this.endPosition.lineNumber)   
    }

    
    validate() {
        
        if (this.startPosition.lineNumber > this.endPosition.lineNumber) {
            return 
        } 
    }




}

let currentDraggedPosition = null 
let isDraggingForComment = false
/** @type {CodeLineRange | null} */
let currentDraggedRange = null
let showMultiLineCommentForm = null
/**
 * 
 * @param {Element} fileElement 
 */
function isFileSelecting(fileElement) {

    return fileElement.closest(".discussion-file-table")?.classList.contains("is-selecting")
}

/**
 * 
 * @param {Element} target 
 * @param {boolean} canExpand 
 */
function setSelection(target, canExpand) {

    const targetLineData = target.id.split("-")
    
    const codeId = targetLineData[0]
    const endLineNumber = targetLineData[1]

    if (currentDraggedRange && currentDraggedRange.codeId == codeId) {

        if (codeId == currentDraggedRange.codeId && endLineNumber > currentDraggedRange.startPosition.lineNumber) {
            return
        }

        const expandedRange = new CodeLineRange(
            codeId,
            currentDraggedRange.startPosition,
            new CodeLinePosition(codeId, endLineNumber)        
        )

        // TODO 

        showMultiLineCommentForm = function() {
            /**@type {HTMLElement | null | undefined}*/
            const button = target.closest("tr")?.querySelector(".add-line-comment")

            if (button && expandedRange) {
                button.click()
            }
        }

    }
    

}

/**
 * 
 * @param {CodeLineRange} expandedRange 
 */
function displayHighlight(expandedRange) {

    if (currentDraggedRange) {
        for (const el of currentDraggedRange.elements) {
            el.classList.remove("selected-line")
        }

        currentDraggedRange = null
    }

    currentDraggedRange = new CodeLineRange(
        expandedRange.codeId, 
        expandedRange.startPosition, 
        expandedRange.endPosition)
    
    
    for (const el of currentDraggedRange.elements) {
        el.classList.add("selected-line")
    }
    
}

/**
 * 
 */
function removeHighlight() {
    const dummyRange = new CodeLineRange("0", 
        new CodeLinePosition("0", "0")
        , new CodeLinePosition("0", "0"))
    
        displayHighlight(dummyRange)
}




/**
 * 
 * @param {Element} codeElement 
 */
function commentDragSelectionIfMouseEnterToCode(codeElement) {
    const target = prevLinkableLine(codeElement)

    if (!target || !isFileSelecting(codeElement)) {
        return 
    }

    setSelection(target, true)
    
}
/**
 * 
 * @param {Element} lineNumberElement 
 */
function commentDragSelectionIfMouseEnterToLineNumber(lineNumberElement) {
    setSelection(lineNumberElement, true)
}

/**
 * 
 * @param {MouseEvent} event 
 */
function handleDragMouseEvent(event) {
    const target = event.target
    if (!(target instanceof Element)) {
        return
    }

    if (currentDraggedPosition) {
        beginDrag()
    }

    const cell = target.closest(".lines-num, .lines-code")

    if (!cell) {
        return 
    }

    if (cell.classList.contains(".lines-num")) {
        commentDragSelectionIfMouseEnterToLineNumber(cell)
        return 
    }
    if (cell.classList.contains(".lines-code")) {
        commentDragSelectionIfMouseEnterToCode
    }


}

// 드래그 시작할 때 이벤트
function beginDrag() {

    if (!currentDraggedPosition) {
        return
    }

    setSelection(currentDraggedPosition, false)
    // TODO 현재 선택된 테이블 가져오기
    /** @type {HTMLElement} */
    const table = currentDraggedPosition.closest(".discussion-file-table")

    table.classList.add("is-selecting")
    currentDraggedPosition = null
    

    document.addEventListener("mouseup", (event) => {
        table.classList.remove("is-selecting")
        table.classList.remove("is-selecting","is-commenting")
        showMultiLineCommentForm && showMultiLineCommentForm()
        showMultiLineCommentForm = null
        
        event.preventDefault()
        
    })


    // TODO 해당 드래그에서 시작점 지정하기

}

/**
 * @param {Element} element
 */
function prevLinkableLine(element) {
    
    if (!(element instanceof HTMLElement)) {
        return
    }

    if (element.classList.contains("lines-num")) {
        return element
    }

    const previousElementSibling = element.previousElementSibling


    if (previousElementSibling) {
        return prevLinkableLine(previousElementSibling)
    }

    return null
    

}

/**
 * 
 * @param {MouseEvent} e 
 */
function handleClick(e) {

    if (!currentDraggedRange) {
        return
    }
    
    
    const target = /**@type {Element} */ (e.target)

    if (target?.closest(".discussion-file-table")) {
        
        return
    }

}

/**
 * 
 * @param {HTMLElement} table 
 */
function addCommentDragSelectionEvent(table) {

    table.addEventListener("mouseenter", handleDragMouseEvent, {capture : true })
} 

/**
 * 
 * @param {HTMLElement} table 
 */
function removeCommentDragSelectionEvent(table) {

    isDraggingForComment = false
    table.removeEventListener("mouseenter", handleDragMouseEvent, {capture : true }) 
    setTimeout(() => {
        document.addEventListener("click", handleClick, {once: true})
    }, 0)

    removeHighlight()
}


/**
 *  
 * @param {MouseEvent} event
 * 

*/
// 버튼을 누를 때 이벤트

on("add-code-line-comment", "mousedown", function(event) {
    if (!(event instanceof MouseEvent)) {
        return
    }

    if (event.button !== 0) {
        return
    }

    const targetElement = /** @type {Element} */ (event.target);
    const parent = targetElement.parentElement;

    if (!parent) {
        return
    }

    const lineNumber =  prevLinkableLine(parent)

    if (!lineNumber) {
        return
    }
    

    /**@type {HTMLElement | null} */
    const table = /** @type {HTMLElement} */ (event.target).closest(".discussion-file-table")

    if (!table) {
        return
    }

    addCommentDragSelectionEvent(table)
    currentDraggedPosition = lineNumber
    isDraggingForComment = true

    event.target?.addEventListener(
        "mouseup",
        function() {
            removeCommentDragSelectionEvent(table)
            currentDraggedPosition = null
            isDraggingForComment = false
        },
        {once: true},
    )

    if (currentDraggedRange && currentDraggedRange.elements.size > 1) {
        event.preventDefault()
    }
    

})