export const decodeArticleId = (encodedId) => {
  try {
    const decodedId = parseInt(window.atob(encodedId), 10);
    if (isNaN(decodedId)) throw new Error("非数字ID");
    return decodedId;
  } catch (e) {
    console.error("解码文章ID失败:", e);
    return 95; // 回退默认值
  }
};
