/**
 * auth.js - 认证管理模块
 */
 
layui.define(function(exports){
  var $ = layui.$;
  var layer = layui.layer;
  
  // 认证相关功能
  var auth = {
    // 登录方法
    login: function(data){
      data = data || {};
  
      // 获取表单数据
      var username = data.username || $('#LAY-user-login-username').val();
      var password = data.password || $('#LAY-user-login-password').val();
  
      // 简单验证
      if(!username) {
        layer.msg('请输入用户名', {icon: 2, anim: 6});
        return;
      }
      if(!password) {
        layer.msg('请输入密码', {icon: 2, anim: 6});
        return;
      }
  
      // 显示loading
      var loadIndex = layer.load(2);
  
      // 发送登录请求
      $.ajax({
        url: '/api/auth/login',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          username: username,
          password: password
        }),
        success: function(res){
          layer.close(loadIndex);
          if(res.code === 200){
            // 登录成功，直接跳转到主页
            layer.msg('登录成功', {
              offset: '15px',
              icon: 1,
              time: 1000
            }, function(){
              location.href = './index.html'; // 后台主页
            });
          } else {
            layer.msg(res.message || '登录失败', {
              offset: '15px',
              icon: 2,
              time: 2000
            });
          }
        },
        error: function(xhr){
          layer.close(loadIndex);
          layer.msg('登录请求失败', {
            offset: '15px',
            icon: 2,
            time: 2000
          });
        }
      });
    },
    
    // 退出登录
    logout: function(){
      layer.confirm('确定要退出吗？', {
        icon: 3,
        title: '提示'
      }, function(index){
        // 发送退出请求
        $.ajax({
          url: '/api/auth/logout',
          type: 'POST',
          contentType: 'application/json',
          xhrFields: {
            withCredentials: true  // 确保跨域请求时携带cookie
          },
          success: function(){
            // 清除本地存储的用户信息
            localStorage.removeItem('user_info');
            localStorage.removeItem('current_software');
            // 跳转到登录页
            location.href = './login.html';
          },
          error: function(){
            // 即使退出接口失败，也清除本地存储并跳转
            localStorage.removeItem('user_info');
            localStorage.removeItem('current_software');
            location.href = './login.html';
          }
        });
        layer.close(index);
      });
    },
    
    // 获取用户信息
    getUserInfo: function(callback) {
      $.ajax({
        url: '/api/auth/getUserInfo',
        type: 'POST',
        contentType: 'application/json',
        success: function(res){
          if(res.code === 200){
            // 将用户信息存储到localStorage
            localStorage.setItem('user_info', JSON.stringify({
              username: res.data.username,
              software_list: res.data.software_list,
              software_agent_info: res.data.software_agent_info
            }));
            
            if(typeof callback === 'function') {
              callback(res.data);
            }
          } else {
            // Session无效
            if(typeof callback === 'function') {
              callback(null, res.message || '获取用户信息失败');
            }
          }
        },
        error: function(xhr){
          if(typeof callback === 'function') {
            callback(null, '请求失败');
          }
        }
      });
    },
    
    // 检查登录状态
    checkLoginStatus: function(callback) {
      this.getUserInfo(function(userData, error) {
        if(userData) {
          if(typeof callback === 'function') {
            callback(true, userData);
          }
        } else {
          if(typeof callback === 'function') {
            callback(false, error);
          }
        }
      });
    },
    
    // 初始化AJAX拦截器
    initAjaxInterceptor: function() {
      $.ajaxSetup({
        complete: function(xhr, status) {
          // 检查是否是未授权响应
          if (xhr.status === 401) {
            // 清除本地存储
            localStorage.removeItem('user_info');
            localStorage.removeItem('current_software');
            
            // 显示提示
            layer.msg('登录已过期，请重新登录', {
              offset: '15px',
              icon: 2,
              time: 2000
            }, function() {
              // 跳转到登录页
              location.href = './login.html';
            });
          }
        }
      });
    }
  };
  
  // 初始化AJAX拦截器
  auth.initAjaxInterceptor();
  
  //对外暴露的接口
  exports('auth', auth);
});

