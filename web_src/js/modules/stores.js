import {reactive} from 'vue';

let diffTreeStoreReactive;
export function diffTreeStore() {
  if (!diffTreeStoreReactive) {
    diffTreeStoreReactive = reactive(window.config.pageData.diffFileInfo);
    window.config.pageData.diffFileInfo = diffTreeStoreReactive;
  }
  return diffTreeStoreReactive;
}

let discussionTreeStoreReactive;
let discussionTreeFileStoreReactive

export function discussionResponseDummy() {

  return {
    "discussionId": 9999,
    "contents": [
      {
        "filePath": "src/main/java/Main.java",
        "codeBlocks": [
          {
            "codeId": 1,
            "lines": [
              { "lineNumber": 1, "content": "public class Main {" },
              { "lineNumber": 2, "content": "    public static void main(String[] args) {" },
              { "lineNumber": 3, "content": "        System.out.println(\"Hello, World!\");" },
              { "lineNumber": 4, "content": "    }" },
              { "lineNumber": 5, "content": "}" }
            ],
            "comments": [
              {
                "id": 1001,
                "scope": "GLOBAL",
                "startLine": 2,
                "endLine": 3,
                "content": "Consider adding more logging for debugging.",
                "reactions": [
                  { "id": 2001, "type": "+1", "discussionId": 9999, "commentId": 1001, "userId": 3001 },
                  { "id": 2002, "type": "hooray", "discussionId": 9999, "commentId": 1001, "userId": 3002 }
                ]
              },
              {
                "id": 1002,
                "scope": "GLOBAL",
                "startLine": 3,
                "endLine": 3,
                "content": "Nice work on this output!",
                "reactions": [
                  { "id": 2003, "type": "heart", "discussionId": 9999, "commentId": 1002, "userId": 3003 }
                ]
              }
            ]
          }
        ]
      },
      {
        "filePath": "csdw/eeewq/test/Helper.java",
        "codeBlocks": [
          {
            "codeId": 2,
            "lines": [
              { "lineNumber": 1, "content": "public class Helper {" },
              { "lineNumber": 2, "content": "    public static String help() {" },
              { "lineNumber": 3, "content": "        return \"Helping...\";" },
              { "lineNumber": 4, "content": "    }" },
              { "lineNumber": 5, "content": "}" }
            ],
            "comments": [
              {
                "id": 1003,
                "scope": "line",
                "startLine": 2,
                "endLine": 2,
                "content": "Consider renaming the method to be more descriptive.",
                "reactions": [
                  { "id": 2004, "type": "eyes", "discussionId": 9999, "commentId": 1003, "userId": 3004 }
                ]
              }
            ]
          }
        ]
      }
    ],
    "globalComments" : {
      
    },
    "globalReactions" : {

    }
  }
}



export function discussionFileTreeStore() {
  if (!discussionTreeFileStoreReactive) {
    discussionTreeFileStoreReactive = reactive({
      repoLink: window.config.pageData.repoLink,
      files: [], 
      selectedItem: null, 
      contents: [], 
      checkedItems: [], 
    });
    window.config.pageData.discussionTreeInfo = discussionTreeFileStoreReactive;
  }
  return discussionTreeFileStoreReactive;
}


export function discussionTreeStore() {
  if (!discussionTreeStoreReactive) {
    discussionTreeStoreReactive = reactive({
      repoLink: window.config.pageData.repoLink,
      files: [], 
      selectedItem: null, 
      contents: [], 
      checkedItems: [], 
    });
    window.config.pageData.discussionTreeInfo = discussionTreeStoreReactive;
  }
  return discussionTreeStoreReactive;
}

let discussionDetailStoreReactive; 
export function discussionDetailStore() {
  if (!discussionDetailStoreReactive) {
    discussionDetailStoreReactive = reactive({
      discussion: window.config.pageData.discussion, 
      discussionContent: window.config.pageData.discussionContent, 
    })
  }
  return discussionDetailStoreReactive;
}