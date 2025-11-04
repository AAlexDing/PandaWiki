export const decodeBase64 = (text: string) => {
  try {
    const buff = Buffer.from(text, 'base64');
    return buff.toString('utf-8');
  } catch (e) {
    // 客户端如果报错，退回到 atob
    if (typeof window !== 'undefined' && window.atob) {
      return window.atob(text);
    }
    // 处理解码失败的情况
    console.error('Base64 decoding failed:', e);
    return '';
  }
};
