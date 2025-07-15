// src/stores/article.js
import { defineStore } from 'pinia';

export const useArticleStore = defineStore('article', {
  state: () => ({
    articleId: null, // 存储当前文章ID
  }),
  actions: {
    // 更新文章ID
    setArticleId(id) {
      this.articleId = id;
      // 可选：同步到localStorage，用于页面刷新后恢复状态
      localStorage.setItem('articleId', id);
    },
  },
});