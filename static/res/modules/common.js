/**
 * common demo
 */
 
layui.define(function(exports){
  var $ = layui.$
  ,layer = layui.layer
  ,laytpl = layui.laytpl
  ,setter = layui.setter
  ,view = layui.view
  ,admin = layui.admin
  
  //公共业务的逻辑处理可以写在此处，切换任何页面都会执行
  //……
  
  // 登录提交
  admin.events.login = function(data){
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
  };

  //退出
  admin.events.logout = function(){
    layui.layer.confirm('确定要退出吗？', {
      icon: 3,
      title: '提示'
    }, function(index){
      // 发送退出请求
      $.ajax({
        url: '/api/auth/logout',
        type: 'POST',
        contentType: 'application/json',
        xhrFields: {
          withCredentials: true  // 添加这一行确保跨域请求时携带cookie
        },
        success: function(){
          // 清除本地存储的用户信息
          localStorage.removeItem('user_info');
          // 跳转到登录页
          location.href = './login.html';
        },
        error: function(){
          // 即使退出接口失败，也清除本地存储并跳转
          localStorage.removeItem('user_info');
          location.href = './login.html';
        }
      });
      layui.layer.close(index);
    });
  };

  // 添加全局 AJAX 拦截器处理未授权响应
  $.ajaxSetup({
    complete: function(xhr, status) {
      // 检查是否是未授权响应
      if (xhr.status === 401) {
        // 清除本地存储
        localStorage.removeItem('user_info');
        
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

  
  //对外暴露的接口
  exports('common', {});
});
