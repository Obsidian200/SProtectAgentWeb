/**
 * software.js - 软件位管理模块
 */
 
layui.define(function(exports){
  var $ = layui.$;
  
  // 软件位相关功能
  var software = {
    // 获取当前选择的软件位
    getCurrentSoftware: function() {
      var currentSoftware = '';
      try {
        // 尝试从localStorage获取当前软件
        currentSoftware = localStorage.getItem('current_software') || '';
        
        // 如果没有当前软件，则尝试从用户信息中获取第一个软件
        if(!currentSoftware) {
          var userInfo = JSON.parse(localStorage.getItem('user_info') || '{}');
          if(userInfo.software_list && userInfo.software_list.length > 0) {
            currentSoftware = userInfo.software_list[0];
            // 设置为当前软件
            localStorage.setItem('current_software', currentSoftware);
          }
        }
      } catch(e) {
        console.log('获取软件位信息失败:', e);
      }
      return currentSoftware;
    },
    
    // 设置当前软件位
    setCurrentSoftware: function(softwareName) {
      if(softwareName) {
        localStorage.setItem('current_software', softwareName);
        return true;
      }
      return false;
    },
    
    // 获取软件位列表
    getSoftwareList: function() {
      var softwareList = [];
      try {
        var userInfo = JSON.parse(localStorage.getItem('user_info') || '{}');
        if(userInfo.software_list && userInfo.software_list.length > 0) {
          softwareList = userInfo.software_list;
        }
      } catch(e) {
        console.log('获取软件位列表失败:', e);
      }
      return softwareList;
    },
    
    // 刷新当前页面的软件位数据
    refreshWithSoftware: function(softwareName) {
      if(!softwareName) {
        softwareName = this.getCurrentSoftware();
      }
      
      // 设置当前软件位
      this.setCurrentSoftware(softwareName);
      
      // 尝试调用当前iframe的刷新方法
      var currentIframe = $('.layadmin-tabsbody-item.layui-show iframe')[0];
      if(currentIframe) {
        try {
          // 尝试调用iframe中的方法
          if(currentIframe.contentWindow.refreshWithSoftware) {
            currentIframe.contentWindow.refreshWithSoftware(softwareName);
          } else {
            // 如果没有专门的刷新方法，则重新加载iframe
            currentIframe.contentWindow.location.reload();
          }
        } catch(e) {
          console.error('刷新页面失败:', e);
        }
      }
    },
    
    // 初始化软件位选择器
    initSoftwareSelector: function(selectElem, callback) {
      var selector = $(selectElem);
      if(!selector.length) return;
      
      // 清空选择器
      selector.empty();
      
      // 获取软件位列表
      var softwareList = this.getSoftwareList();
      var currentSoftware = this.getCurrentSoftware();
      
      if(softwareList && softwareList.length > 0) {
        // 添加软件位选项
        softwareList.forEach(function(software) {
          var selected = software === currentSoftware ? 'selected' : '';
          selector.append('<option value="' + software + '" ' + selected + '>' + software + '</option>');
        });
        
        // 如果没有选中的软件位且有软件位列表，默认选中第一个
        if(!currentSoftware && softwareList.length > 0) {
          this.setCurrentSoftware(softwareList[0]);
          // 设置第一个选项为选中状态
          selector.val(softwareList[0]);
        }
      }
      
      // 重新渲染表单
      if(layui.form) {
        layui.form.render('select');
      }
      
      // 如果有回调函数，执行回调
      if(typeof callback === 'function') {
        callback(currentSoftware);
      }
    }
  };
  
  //对外暴露的接口
  exports('software', software);
});


