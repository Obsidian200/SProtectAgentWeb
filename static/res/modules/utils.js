/**
 * util.js - 自定义工具函数模块
 */
 
layui.define(function(exports){
  // 通用工具函数
  var utils = {
    // 格式化时间戳
    formatTimestamp: function(timestamp) {
      if (!timestamp || timestamp <= 0) {
        return '';
      }
      
      var date = new Date(timestamp * 1000);
      var year = date.getFullYear();
      var month = ('0' + (date.getMonth() + 1)).slice(-2);
      var day = ('0' + date.getDate()).slice(-2);
      var hours = ('0' + date.getHours()).slice(-2);
      var minutes = ('0' + date.getMinutes()).slice(-2);
      var seconds = ('0' + date.getSeconds()).slice(-2);
      
      return year + '-' + month + '-' + day + ' ' + hours + ':' + minutes + ':' + seconds;
    },
    
    // 生成随机字符串
    randomString: function(length) {
      var chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
      var result = '';
      for (var i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
      }
      return result;
    },
    
    // 格式化金额
    formatMoney: function(amount, decimals) {
      decimals = decimals || 2;
      return parseFloat(amount).toFixed(decimals).replace(/\d(?=(\d{3})+\.)/g, '$&,');
    },
    
    // 获取URL参数
    getUrlParam: function(name) {
      var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
      var r = window.location.search.substr(1).match(reg);
      if (r != null) return decodeURIComponent(r[2]); return null;
    },
    
    // 复制文本到剪贴板
    copyToClipboard: function(text, callback) {
      var textarea = document.createElement('textarea');
      textarea.value = text;
      textarea.style.position = 'fixed';
      textarea.style.opacity = 0;
      document.body.appendChild(textarea);
      textarea.select();
      
      try {
        var successful = document.execCommand('copy');
        document.body.removeChild(textarea);
        if(typeof callback === 'function') {
          callback(successful);
        }
        return successful;
      } catch (err) {
        document.body.removeChild(textarea);
        if(typeof callback === 'function') {
          callback(false);
        }
        return false;
      }
    }
  };
  
  //对外暴露的接口
  exports('utils', utils);
});


